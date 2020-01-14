package quin

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	defaultPingCount = "5"

	pingOutput = regexp.MustCompile(`time=([.0-9]+) ms`)
)

// PingPeers ping all nodes we've found
func PingPeers(endpoints *Poller, port int) {
	peers := endpoints.getPeers()
	log.Printf("running ping for peers: %v", peers)
	var wg sync.WaitGroup
	for _, peer := range peers {
		peer := peer
		wg.Add(1)
		go func() {
			log.Printf("pinging %s", peer.Hostname)
			stats, err := ping(peer, &wg)
			if err != nil {
				log.Printf("ERROR: unable to ping peer %s: %s", peer.Hostname, err)
				return
			}
			log.Printf("ping %s %g", stats.Peer.Hostname, stats.RTTs)
			for _, rtt := range stats.RTTs {
				pings.With(prometheus.Labels{"hostname": stats.Peer.Hostname}).Observe(rtt)
			}
		}()
	}
}

type pingStats struct {
	Peer target
	RTTs []float64
}

func ping(peer target, wg *sync.WaitGroup) (*pingStats, error) {
	defer wg.Done()

	output, err := exec.Command("ping", "-c", defaultPingCount, peer.Hostname).CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("Failed to run ping: %v", err)
	}

	match := pingOutput.FindAllStringSubmatch(string(output), -1)
	pingData := pingStats{
		Peer: peer,
	}
	for _, ms := range match {
		log.Println(ms)
		ping, err := strconv.ParseFloat(ms[1], 64)
		if err != nil {
			return nil, fmt.Errorf("Failed to convert to int %s", match[2])
		}
		pingData.RTTs = append(pingData.RTTs, ping/1000)
	}

	return &pingData, nil
}
