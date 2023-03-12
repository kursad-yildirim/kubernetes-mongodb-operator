#!/bin/bash
echo Run the operator
go mod tidy
clear
MATCH_NAMESPACE=tuff KUBECONFIG=/home/kyildiri/.kube/config make run
