Below is the same metrics table sorted **P0 → P1 → P2**, with components grouped within each priority.

Priority scale: **P0** = cluster outage/data-loss emergency; **P1** = major production impact; **P2** = degradation or growing operational risk.

| Priority | Component | Datadog metric / service check | What it monitors | Failure condition and impact |
|---|---|---|---|---|
| P0 | API server | `kube_apiserver.slis.kubernetes_healthcheck`, `slis.kubernetes_healthcheck_total` | API-server internal health | Failed checks, no data, or all endpoints unavailable means the cluster API is partially or completely unavailable |
| P0 | API server | `kube_apiserver.apiserver_request_total.count` filtered to HTTP `5xx` | Server-side API request failures | Controllers, kubelets, operators, and users cannot reliably read or modify cluster state |
| P0 | API server → etcd | `kube_apiserver.etcd_request_duration_seconds.*`, `etcd_request_errors_total`, `etcd_requests_total` | API-server storage calls to etcd | Storage latency or errors affect almost every Kubernetes read/write operation |
| P0 | etcd | `etcd.prometheus.health` or `etcd.healthy` | Health and reachability of every etcd member | Failed members reduce redundancy; loss of quorum makes Kubernetes state unavailable |
| P0 | etcd | `etcd.server.has_leader` | Whether an etcd member sees a leader | `0` means the member is unavailable; no leader across quorum means complete control-plane outage |
| P0 | etcd | `etcd.server.proposals.pending`, `proposals.failed.total` | Raft proposals waiting or failing | Growing pending proposals or failures make Kubernetes writes slow or unavailable and may indicate quorum loss |
| P0 | etcd capacity | `etcd.mvcc.db.total.size.in_use.bytes`, `etcd.debugging.mvcc.db.total.size.in_bytes`, `etcd.server.quota.backend.bytes` | Logical database usage, allocated size, and backend quota | Reaching quota raises a `NOSPACE` alarm and causes Kubernetes state-changing requests to be rejected |
| P0 | etcd network | `etcd.network.peer.round_trip_time.seconds`, `peer.sent.failures.total`, `peer.received.failures.total`, `disconnected_peers.total` | Raft peer latency, errors, and disconnections | Followers lag, elections increase, and quorum may be lost |
| P1 | Kubelet | `kubernetes.kubelet.check` | Overall kubelet health | The node stops reliably reconciling pods; new and restarted containers may not start |
| P1 | Kubelet | `kubernetes.kubelet.check.ping` | Whether the kubelet endpoint responds | Control plane, Datadog, and other clients cannot communicate reliably with the node |
| P1 | Kubelet | `kubernetes.kubelet.check.syncloop` | Main kubelet pod-synchronization loop | Pod starts, stops, updates, probes, and volume changes stop being reconciled |
| P1 | Kubelet | `kubernetes.kubelet.pleg.last_seen`, `pleg.relist_duration.*`, `pleg.relist_interval.*`, `pleg.discard_events` | Pod Lifecycle Event Generator and runtime-state discovery | Kubelet loses track of container state; pods can remain stuck or incorrectly reported |
| P1 | Kubelet / CRI | `kubernetes.kubelet.runtime.errors`, `runtime.operations.duration.*`, `runtime.operations` | Container-runtime operations | Pods remain `Pending`, `ContainerCreating`, `Unknown`, or fail to terminate |
| P1 | Kubelet / CNI | `kubernetes.kubelet.network_plugin.latency.*`, `kubernetes.network.rx_errors`, `tx_errors`, `rx_dropped`, `tx_dropped` | CNI latency and container-network errors | Pod sandbox creation fails or pod traffic becomes unavailable/intermittent |
| P1 | Kubelet volumes | `kubernetes.kubelet.volume.stats.available_bytes`, `used_bytes`, `inodes_free`, `inodes_used` | Persistent-volume capacity and inode usage | Applications cannot write; databases may fail, become read-only, or risk corruption |
| P1 | Kubelet storage | `kubernetes.node.filesystem.usage_pct`, `node.image.filesystem.usage_pct`, `ephemeral_storage.usage` | Root, image, container-log, and ephemeral-storage pressure | Image pulls and log writes fail; pods are evicted and the node can enter `DiskPressure` |
| P1 | Kubelet | `kubernetes.kubelet.evictions` | Pods evicted by kubelet | Causes workload disruption and possible loss of non-persistent data |
| P1 | API server ↔ kubelet TLS | `kubernetes.apiserver.certificate.expiration.*` where available | Remaining lifetime of certificates used for authenticated communication | Expiry can break kubelet operations, including logs, exec, metrics, and node management |
| P1 | API server | `kube_apiserver.apiserver_request_total.count` filtered to HTTP `429`, `401`, and `403` | Throttled, unauthenticated, and forbidden API requests | Controllers may be throttled, legitimate clients may lose access, or a security issue may exist |
| P1 | API server | `kube_apiserver.request_duration_seconds.sum/count` | API latency by verb, resource, and request type | Causes slow deployments, controller lag, kubectl timeouts, and failed health checks |
| P1 | API server | `kube_apiserver.current_inflight_requests`, `longrunning_gauge` | Concurrent and long-running requests | Requests queue or time out when concurrency is exhausted |
| P1 | API server / APF | `kube_apiserver.flowcontrol_current_inqueue_requests`, `flowcontrol_rejected_requests_total.count`, `flowcontrol_request_wait_duration_seconds.*` | API Priority and Fairness queueing and rejection | Clients receive 429 responses; kubelets and controllers may be starved |
| P1 | API server | `kube_apiserver.apiserver_dropped_requests_total.count`, `apiserver_request_terminations_total.count` | Requests dropped or terminated for self-protection | Indicates overload severe enough that the API server is actively shedding work |
| P1 | API server admission | `kube_apiserver.admission_webhook_admission_latencies_seconds.*`, `apiserver_admission_webhook_fail_open_count.count`, `apiserver_admission_webhook_request_total.count` | Admission webhook latency, requests, and fail-open behavior | Object creation/update is blocked, delayed, or allowed without expected policy enforcement |
| P1 | API server authentication | `kube_apiserver.authentication_attempts.count`, `authentication_duration_seconds.*`, `authenticated_user_requests.count` | Authentication traffic, failures, and latency | Legitimate users/controllers lose access or the cluster may be experiencing credential attacks |
| P1 | API aggregation | `kube_apiserver.aggregator_unavailable_apiservice` | Availability of aggregated API services | Metrics APIs, custom metrics, extension APIs, and HPA integrations may fail |
| P1 | API server process | `kube_apiserver.process_cpu_total`, `process_resident_memory_bytes`, `go_goroutines`, `go_threads` | API-server CPU, memory, goroutines, and threads | Saturation or leaks increase latency and can eventually restart or crash the API server |
| P1 | API server audit | `kube_apiserver.audit_event.count` | Volume of audit events submitted to the backend | A drop to zero can mean lost security visibility; excessive volume can increase API latency |
| P1 | etcd | `etcd.server.leader.changes.seen.total` | Frequency of etcd leader elections | Rapid changes produce high API latency and indicate unstable disk, networking, or load |
| P1 | etcd | `etcd.server.proposals.committed.total` minus `proposals.applied.total` | Backlog between committed and applied entries | A continuously growing gap means a member is overloaded or unhealthy |
| P1 | etcd disk | `etcd.disk.wal.fsync.duration.seconds.*`, `disk.backend.commit.duration.seconds` | WAL persistence and backend commit latency | Slow disk causes slow writes, missed heartbeats, leader changes, and cluster instability |
| P1 | etcd | `etcd.server.apply.slow.total`, `read_indexes.slow.total`, `read_indexes.failed.total` | Slow apply operations and linearizable reads | API reads and writes stall, often because of disk or CPU saturation |
| P1 | etcd process | `etcd.process.open.fds`, `process.max.fds`, `process.resident.memory.bytes`, `process.cpu.seconds.total` | File descriptors, memory, and CPU | etcd may fail to open WAL/network files, panic, restart, or become unresponsive |
| P1 | Controller manager | `kube_controller_manager.up`, `prometheus.health` | Controller-manager availability | Deployments, ReplicaSets, Jobs, Nodes, endpoints, and other objects stop converging |
| P1 | Controller manager | `kube_controller_manager.leader_election.status`, `leader_election.transitions` | Leader availability and stability | Reconciliation pauses when no leader exists; frequent transitions create repeated interruption |
| P1 | Controller manager | `kube_controller_manager.queue.depth`, `queue.retries`, `queue.queue_duration.*` | Controller workqueue backlog, retries, and delay | Rollouts, node recovery, endpoint updates, Jobs, and other control loops become delayed |
| P1 | Controller manager | `queue.work_longest_duration`, `queue.work_unfinished_duration`, `queue.process_duration.*` | Stuck or long-running controller workers | A controller may be blocked, deadlocked, or waiting indefinitely on a dependency |
| P1 | Controller manager / nodes | `kube_controller_manager.nodes.unhealthy`, `nodes.evictions`, `nodes.count` | Unhealthy, evicted, and registered nodes | Pods may be mass-evicted and workload capacity can decline |
| P1 | Scheduler | `kube_scheduler.up`, `prometheus.health` | Scheduler availability | Existing pods continue running, but new and replacement pods cannot be placed |
| P1 | Scheduler | `kube_scheduler.leader_election.status` | Active scheduler leader | No leader means all new pod scheduling stops |
| P1 | Scheduler | `kube_scheduler.pending_pods` by queue | Active, backoff, and unschedulable queue sizes | Deployments, scaling, failover, and Jobs remain pending |
| P1 | Scheduler | `kube_scheduler.schedule_attempts` by result | Successful, unschedulable, and internal-error attempts | Internal errors prevent valid scheduling; widespread unschedulable results indicate capacity or constraint problems |
| P1 | Scheduler | `kube_scheduler.binding_duration.*`, `volume_scheduling_duration.*` | Pod binding and storage-aware scheduling latency | Pods are selected for nodes but remain pending or unbound |
| P1 | Scheduler | `kube_scheduler.client.http.requests`, `client.http.requests_duration.*` | Scheduler communication with the API server | API errors, throttling, or latency make the scheduler cache stale and cause binding failures |
| P2 | Kubelet | `kubernetes.kubelet.pod.start.duration`, `pod.worker.duration`, `pod.worker.start.duration` | Pod startup and synchronization latency | Produces slow deployments, slow autoscaling, and delayed workload recovery |
| P2 | Kubelet probes | `kubernetes.liveness_probe.failure.total`, `readiness_probe.failure.total`, `startup_probe.failure.total` | Health-probe failures | Liveness restarts containers, readiness removes endpoints, and startup failures block availability |
| P2 | Kubelet workloads | `kubernetes.containers.restarts`, `containers.state.waiting`, `containers.state.terminated` | Restart loops and containers that cannot run | Indicates `CrashLoopBackOff`, image, configuration, secret, or application failures |
| P2 | Kubelet process | `kubernetes.kubelet.cpu.usage`, `kubernetes.kubelet.memory.usage`, `kubernetes.kubelet.memory.rss` | Kubelet CPU and memory consumption | Saturation causes slow reconciliation and can lead to PLEG problems or kubelet OOM |
| P2 | Controller manager process | `goroutines`, `open_fds`, `max_fds`, plus component container CPU/memory | Resource saturation and leaks | Eventually causes slow reconciliation, leader loss, OOM, or restart |
| P2 | Scheduler | `kube_scheduler.scheduling.attempt_duration.*`, `scheduling.pod.scheduling_duration.*`, `scheduling.e2e_scheduling_duration.*` | Scheduling-cycle and end-to-end latency | Causes slow deployments, autoscaling, failover, and workload recovery |
| P2 | Scheduler | `kube_scheduler.pod_preemption.attempts`, `pod_preemption.victims.*` | Preemption activity | Indicates capacity/priority pressure and disruption of lower-priority workloads |
| P2 | Scheduler process | `goroutines`, `open_fds`, `max_fds`, plus component container CPU/memory | Scheduler process saturation and leaks | Eventually causes slow scheduling, leader loss, OOM, or restart |

Within each priority, I would implement alerts in this order:

**API server → etcd → kubelet → scheduler → controller manager → degradation/capacity metrics.**

Sources: [Datadog Kubernetes control-plane monitoring](https://docs.datadoghq.com/containers/kubernetes/control_plane/), [Datadog kubelet integration](https://docs.datadoghq.com/integrations/kubelet/), [official etcd metric semantics](https://etcd.io/docs/v3.6/metrics/).
