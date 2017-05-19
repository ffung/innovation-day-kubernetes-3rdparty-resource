#!/bin/bash
export GOPATH="/Users/fai/work/xebia/go"
export GOOS="linux"
export GOARCH="amd64"

go fmt
go get ./...
go build
cp innovation-day-kubernetes-3rdparty-resource environment-manager

eval $(minikube docker-env)
docker build -t environment-manager:bla .

kubectl delete pods --all
kubectl apply -f environment-manager.yaml
