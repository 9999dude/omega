# Local Kubernetes lab with Kind

This lab creates a local, two-node Kubernetes cluster for interview practice:

- Kubernetes API server: `v1.36.2`
- Kind CLI: `v0.32.0`
- kubectl CLI: `v1.36.2`
- Nodes: one control plane and two workers
- Per-node hard limit: 2 CPUs and 2 GiB RAM, with swap disabled
- Cluster name: `interview-k8s`
- kubectl context: `kind-interview-k8s`

## Prerequisite

Install and start [Docker](https://docs.docker.com/engine/install/). The setup
downloads its own pinned Kind and kubectl binaries into `.tools/`, so global
installations of those tools are not required.

For three nodes without oversubscription, configure Docker Desktop with at least
6 CPUs and 8 GB of memory. The node limits have a combined ceiling of 6 CPUs and
6 GiB; the extra memory leaves room for Docker's VM and container overhead.

The first setup is slower because Kind must build a local node image for the
exact Kubernetes version requested. The Kubernetes server bundle is cached under
`.tools/cache/`; an interrupted download resumes when `make cluster-up` is run
again.

## Create the cluster

From this directory, run:

```bash
make cluster-up
```

The command is safe to run again. If the cluster already exists, it verifies the
existing cluster instead of replacing it.

Kind does not expose per-node resource fields in `kind-config.yaml`. After Kind
creates the containers, the setup applies Docker hard limits to every node and
`make verify` checks those limits.

Kind nodes are containers rather than separate virtual machines. Kubernetes may
therefore display Docker Desktop's shared VM capacity when you inspect a Node,
instead of the container's smaller limit. Docker still enforces the 2 CPU / 2
GiB cap on each node; the verification command checks Docker's effective values.

Verify it at any time:

```bash
make verify
```

Inspect the cluster:

```bash
make status
```

Use kubectl directly with the course-pinned client:

```bash
./.tools/bin/kubectl --context kind-interview-k8s get nodes
./.tools/bin/kubectl --context kind-interview-k8s get pods --all-namespaces
```

## Delete the cluster

```bash
make cluster-down
```

The locally built `interview/kind-node:v1.36.2` Docker image is retained so the
next cluster creation is faster. Remove that image manually if disk space is
needed:

```bash
docker image rm interview/kind-node:v1.36.2
```

## Optional names

Override the default cluster name without changing any files:

```bash
CLUSTER_NAME=my-k8s-lab make cluster-up
CLUSTER_NAME=my-k8s-lab make cluster-down
```

## Why the node image is built locally

Kind selects the Kubernetes version through its node image. Kind did not publish
an official prebuilt node image for exactly `v1.36.2`; Kind `v0.32.0` ships with
Kubernetes `v1.36.1`. The setup therefore uses Kind's supported release-build
workflow to build a local node image from the official Kubernetes `v1.36.2`
server release before creating the cluster.
