package quin

import (
	"fmt"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	clusterConfig = "kubernetes"
)

type target struct {
	Hostname string `json:"host"`
}

// Poller containing poll data
type Poller struct {
	config   string
	peers    []target
	requests chan chan []target
}

// NewEndpointGrabber this returns endpoints to use for operations
func NewEndpointGrabber(clusterConfig string, peerPollSec int, kubernetesInternal bool) *Poller {
	clientset, err := kubernetesClientset(kubernetesInternal)
	if err != nil {
		fmt.Errorf("Failing to grab endpoints: %s", err)
	}

	p := &Poller{
		config:   clusterConfig,
		requests: make(chan chan []target),
	}
	p.updatePeers(clientset)
	go p.run(clientset, peerPollSec)
	return p
}

func kubernetesClientset(internal bool) (*kubernetes.Clientset, error) {
	config, err := kubernetesConfig(internal)
	if err != nil {
		return nil, fmt.Errorf("ERROR: %s", err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("ERROR: %s", err.Error())
	}

	return clientset, nil
}

func kubernetesConfig(internal bool) (*rest.Config, error) {
	if internal {
		return rest.InClusterConfig()
	}

	usr, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("ERROR: %s", err.Error())
	}
	kubeconfig := string(filepath.Join(string(usr.HomeDir), ".kube", "config"))

	return clientcmd.BuildConfigFromFlags("", kubeconfig)
}

func (p *Poller) run(clientset *kubernetes.Clientset, peerPollSec int) {
	// Refresh peers every peerPollSec
	ticker := time.NewTicker(time.Duration(peerPollSec) * time.Second)
	for {
		select {
		case _ = <-ticker.C:
			p.updatePeers(clientset)
		case req := <-p.requests:
			peerCopy := make([]target, 0, len(p.peers))
			for _, target := range p.peers {
				peerCopy = append(peerCopy, target)
			}
			req <- peerCopy
		}
	}
}

func (p *Poller) getPeers() []target {
	req := make(chan []target)
	p.requests <- req
	return <-req
}

func (p *Poller) updatePeers(clientset *kubernetes.Clientset) {
	targets := []target{}

	nodes, err := clientset.CoreV1().Nodes().List(metav1.ListOptions{})
	for _, item := range nodes.Items {
		for _, address := range item.Status.Addresses {
			fmt.Printf("address.Type %s\n", string(address.Type))
			fmt.Printf("address.Address %s\n", string(address.Address))
			if strings.ToLower(string(address.Type)) == "internalip" {
				t := target{string(address.Address)}
				targets = append(targets, t)
			}
		}
	}
	if err != nil {
		k8sConnFailures.With(prometheus.Labels{"source": "kubernetes"}).Inc()
	}
	p.peers = targets
}
