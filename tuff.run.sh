#!/bin/bash

go mog tidy
clear
MATCH_NAMESPACE=tuff KUBECONFIG=/home/kyildiri/.kube/config make run
