SHELL := /bin/bash

TOOLS_DIR := $(CURDIR)/.tools/bin
export PATH := $(TOOLS_DIR):$(PATH)

.PHONY: help tools cluster-up cluster-down verify status

help: ## Show the available commands.
	@awk 'BEGIN {FS = ":.*## "; printf "Kubernetes local lab\n\n"} /^[a-zA-Z_-]+:.*## / {printf "  %-14s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

tools: ## Install the pinned Kind and kubectl binaries locally.
	@./scripts/install-tools.sh

cluster-up: tools ## Build the Kubernetes v1.36.2 node image and create the cluster.
	@./scripts/create-cluster.sh

cluster-down: tools ## Delete the local cluster.
	@./scripts/delete-cluster.sh

verify: tools ## Check the Kubernetes version and cluster health.
	@./scripts/verify-cluster.sh

status: tools ## Show nodes and all system pods.
	@kubectl --context "kind-$${CLUSTER_NAME:-interview-k8s}" get nodes -o wide
	@kubectl --context "kind-$${CLUSTER_NAME:-interview-k8s}" get pods --all-namespaces
