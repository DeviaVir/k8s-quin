FROM golang:alpine AS builder

ENV GO111MODULE="on"
ENV CGO_ENABLED="0"

RUN apk add --update git

RUN mkdir -p /go/src/github.com/DeviaVir/k8s-quin

COPY . /go/src/github.com/DeviaVir/k8s-quin

RUN cd /go/src/github.com/DeviaVir/k8s-quin \
 && go mod vendor \
 && go build \
      -mod vendor \
      -o /go/bin/quin

FROM alpine
COPY --from=builder /go/bin/quin /usr/local/bin/quin
CMD ["/usr/local/bin/quin"]
