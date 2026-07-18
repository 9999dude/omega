These fields control three different areas: service discovery, scheduler priority/preemption, and node resource protection.

| Field | Set by | Mainly used by |
|---|---|---|
| `enableServiceLinks` | Pod author/defaulting | kubelet |
| `preemptionPolicy` | Pod/PriorityClass/defaulting | kube-scheduler |
| `priority` | Priority admission from a `PriorityClass` | kube-scheduler and kubelet |
| `qosClass` | Computed by Kubernetes | kubelet and Linux resource management |

### `enableServiceLinks: true`

This tells the kubelet to inject information about Kubernetes Services into the container as environment variables.

For a Service named `redis-primary`, the container might receive:

```text
REDIS_PRIMARY_SERVICE_HOST=10.96.20.15
REDIS_PRIMARY_SERVICE_PORT=6379
REDIS_PRIMARY_PORT=tcp://10.96.20.15:6379
```

The field defaults to `true`. It exists mainly for compatibility with Docker’s old “container links” behavior. The kubelet constructs these variables before starting the container. [Kubernetes Pod API reference](https://kubernetes.io/docs/reference/kubernetes-api/core/pod-v1/)

Significance:

- Only Services that already exist when the Pod starts can be represented.
- The variables do not update when Services are added or changed.
- It does not control Service networking, routing, CoreDNS, or whether the Pod can contact a Service.
- Modern applications normally discover Services through DNS, such as `redis-primary.default.svc.cluster.local`.

Setting it to `false` disables ordinary Service-link environment variables. Kubernetes API-related variables such as `KUBERNETES_SERVICE_HOST` may still be supplied for compatibility.

You can inspect what was injected with:

```bash
kubectl exec kubernetes-bootcamp-66d5c96cb-ndvsp -- env
```

### `preemptionPolicy: PreemptLowerPriority`

This controls what the scheduler is permitted to do when this Pod is waiting to be scheduled but no node has enough available resources.

`PreemptLowerPriority` means:

> The scheduler may terminate Pods whose numeric priority is lower than this Pod’s priority, if doing so would make room for this Pod.

The other value is:

```yaml
preemptionPolicy: Never
```

A `Never` Pod can still be placed ahead of lower-priority Pods in the scheduler queue, but it cannot remove existing Pods to make room.

For this particular Pod:

```yaml
priority: 0
preemptionPolicy: PreemptLowerPriority
```

It can preempt only Pods with a priority below `0`. It cannot preempt other ordinary priority-`0` Pods because victim priority must be strictly lower. Higher-priority Pods can still preempt this Pod. [Pod priority and preemption](https://kubernetes.io/docs/concepts/scheduling-eviction/pod-priority-preemption/)

This field matters only when the Pod is pending and unschedulable. It has no ongoing effect while this Pod is happily running.

### `priority: 0`

This is the Pod’s resolved numeric importance.

Normally you configure a symbolic `PriorityClass`:

```yaml
spec:
  priorityClassName: production-critical
```

For example:

```yaml
apiVersion: scheduling.k8s.io/v1
kind: PriorityClass
metadata:
  name: production-critical
value: 100000
preemptionPolicy: PreemptLowerPriority
```

The priority admission controller resolves that class and stores its numeric value in `spec.priority`.

Your Pod has no `priorityClassName`. If the cluster has no global default `PriorityClass`, Kubernetes resolves that to `0`.

The value is used for:

- Scheduler queue order: higher-priority pending Pods are generally considered first.
- Scheduler preemption: a pending Pod can remove only lower-priority Pods.
- Kubelet node-pressure eviction: priority participates in deciding which running Pods to evict.

A priority of `0` simply means normal/default importance. It does not mean “disabled.”

### `qosClass: BestEffort`

Unlike the other three fields, this is under `status`, not `spec`:

```yaml
status:
  qosClass: BestEffort
```

You do not set it directly. Kubernetes computes it from CPU and memory requests and limits.

Your container has:

```yaml
resources: {}
```

Therefore it has:

- No CPU request
- No memory request
- No CPU limit
- No memory limit

That produces the `BestEffort` QoS class. Kubernetes has three QoS classes:

| QoS class | Resource configuration | Protection |
|---|---|---|
| `Guaranteed` | Every container has equal CPU/memory requests and limits | Highest |
| `Burstable` | At least one CPU/memory request or limit, but not Guaranteed | Medium |
| `BestEffort` | No CPU/memory requests or limits | Lowest |

Consequences for this Pod:

- It has no reserved CPU or memory.
- The scheduler accounts for effectively zero requested CPU and memory when placing it.
- It can consume spare resources while the node has capacity.
- It has no explicit CPU or memory ceiling in this manifest.
- It is among the most vulnerable Pods when the node experiences resource pressure.
- On Linux, BestEffort containers receive a high `oom_score_adj`, making them likely targets if the kernel must kill something during an out-of-memory event.

More precisely, kubelet eviction ranks Pods using whether usage exceeds requests, Pod priority, and usage relative to requests. Because a BestEffort Pod has zero requests, any positive usage exceeds its requests. [QoS classes](https://kubernetes.io/docs/concepts/workloads/pods/pod-qos/) and [node-pressure eviction](https://kubernetes.io/docs/concepts/scheduling-eviction/node-pressure-eviction/)

### How the four fields interact here

```text
Pod created
   │
   ├─ kubelet injects Service environment variables
   │     enableServiceLinks: true
   │
   ├─ scheduler orders and places the Pod
   │     priority: 0
   │     preemptionPolicy: PreemptLowerPriority
   │
   └─ kubelet manages the running Pod under resource pressure
         qosClass: BestEffort
         priority: 0
```

The important distinction is:

- Scheduler preemption primarily uses `priority`, not QoS.
- Node-pressure eviction uses resource usage versus requests and `priority`; `BestEffort` strongly predicts vulnerability because its requests are zero.

Because this Pod belongs to a ReplicaSet, any changes should be made to the Deployment’s Pod template. For a real workload, defining resource requests and limits is usually the most important improvement:

```yaml
resources:
  requests:
    cpu: 100m
    memory: 128Mi
  limits:
    cpu: 500m
    memory: 256Mi
```

That would change the Pod from `BestEffort` to `Burstable`.
