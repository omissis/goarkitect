#!/bin/sh

cd /root
curl https://dl.google.com/go/go1.19.1.linux-arm64.tar.gz --output go.tgz
tar xvfz go.tgz
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.49.0
apt-get install -y git make
export PATH=/root/go/bin:$PATH
