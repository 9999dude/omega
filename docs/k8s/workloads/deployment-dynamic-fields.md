These fields belong to four different layers; they are not all part of the Deployment’s `spec`:

| Location | Written/used primarily by |
|---|---|
| `metadata.*` | API server and Kubernetes clients/controllers |
| `spec.*` | Deployment controller |
| `spec.template.spec.*` | Scheduler and kubelet after Pods are created |
| `status.conditions[]` | Deployment controller, monitoring tools, `kubectl` |

## Deployment metadata

### `deployment.kubernetes.io/revision: "1"`

A Deployment rollout revision number.

- Revision `1` means this is the first version of the Pod template.
- It increments when `.spec.template` changes—for example, changing the container image.
- Scaling `replicas` does not create a new rollout revision.
- The Deployment controller records the revision on the corresponding ReplicaSet.
- `kubectl rollout history` and `kubectl rollout undo` use this history.
- It is controller-managed and should not be edited manually.

It is not the same thing as `generation` or `resourceVersion`. [Kubernetes revision annotation documentation](https://kubernetes.io/docs/reference/labels-annotations-taints/#deployment-kubernetes-io-revision)

### `creationTimestamp: "2026-07-15T15:27:26Z"`

The API-server time when this particular Deployment object was created.

It is mainly useful for:

- Displaying object age.
- Auditing and troubleshooting.
- Determining the lifetime of this particular object.

It does not determine when the containers started, and Kubernetes does not automatically delete an object based on this timestamp.

### `generation: 1`

The version of the Deployment’s desired state.

For a Deployment, the API server increments this when the `spec` changes:

```text
metadata.generation = latest desired-state version
status.observedGeneration = version processed by controller
```

In your object:

```yaml
metadata:
  generation: 1
status:
  observedGeneration: 1
```

This means the Deployment controller has processed the latest `spec`.

An important distinction:

- Changing `replicas` normally increments `generation`.
- Changing `replicas` does not increment the rollout revision because the Pod template did not change.
- Status-only updates do not increment `generation`.

### `uid: 92588eaa-e58c-44a5-bc17-153040c68e33`

A cluster-wide unique identity for this particular occurrence of the Deployment.

Suppose you run:

```bash
kubectl delete deployment kubernetes-bootcamp
kubectl create deployment kubernetes-bootcamp ...
```

The new Deployment has the same name but a different UID. Kubernetes therefore knows it is a completely new object.

UIDs are important for:

- Owner references.
- Garbage collection.
- Preventing a newly created object with the same name from being mistaken for an old object.
- Safely identifying objects in controllers and API requests.

For example, the generated ReplicaSet normally has an owner reference containing this Deployment’s name and UID. The garbage collector uses that relationship when the Deployment is deleted. [Kubernetes object names and UIDs](https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#uids)

### `resourceVersion: "2012"`

An opaque storage version for the current representation of this object.

It changes whenever the stored Deployment changes, including many:

- `spec` changes.
- `status` updates.
- Annotation or label changes.

Clients use it for:

1. Optimistic concurrency control—preventing one client from silently overwriting another client’s newer change.
2. Watching resources from a known point.
3. Detecting whether an object has changed.

`"2012"` does not mean the Deployment was updated 2,012 times. Clients should treat it as an opaque API-server value and pass it back unchanged when required. [Kubernetes API resource-version documentation](https://kubernetes.io/docs/reference/using-api/api-concepts/#resource-versions)

### Revision versus generation versus resourceVersion

| Field | Changes when | Purpose |
|---|---|---|
| Rollout revision | Pod template changes | Rollout history and rollback |
| `generation` | Desired `spec` changes | Controller reconciliation tracking |
| `resourceVersion` | Stored object changes, including status | Concurrency and watches |
| `uid` | Never; replacement gets a new UID | Object identity |

## Deployment settings

### `revisionHistoryLimit: 10`

The Deployment controller may retain up to 10 old ReplicaSets after rollouts.

Those old ReplicaSets contain previous Pod templates and make rollback possible:

```bash
kubectl rollout history deployment/kubernetes-bootcamp
kubectl rollout undo deployment/kubernetes-bootcamp
```

It means 10 old ReplicaSets, not 10 old Pods. The current ReplicaSet is separate from that limit.

The Deployment controller performs the cleanup. Setting it to `0` allows old ReplicaSets to be cleaned up after rollout completion, which removes normal rollback history. [Deployment API reference](https://kubernetes.io/docs/reference/kubernetes-api/apps/deployment-v1/)

### `progressDeadlineSeconds: 600`

The Deployment has 600 seconds to demonstrate rollout progress.

Progress includes events such as:

- A new ReplicaSet being created.
- New Pods becoming ready.
- More updated replicas becoming available.
- Old replicas being scaled down.

If progress stalls beyond the deadline, the Deployment controller reports something similar to:

```yaml
type: Progressing
status: "False"
reason: ProgressDeadlineExceeded
```

This field:

- Detects and reports a stalled rollout.
- Does not stop the controller from continuing to reconcile.
- Does not automatically roll back the Deployment.
- Is not counted while the Deployment is paused.

Monitoring systems and `kubectl rollout status` can surface the resulting failure condition. [Deployment progress-deadline documentation](https://kubernetes.io/docs/reference/kubernetes-api/apps/deployment-v1/)

## Container termination information

These fields belong under:

```yaml
spec:
  template:
    spec:
      containers:
        - terminationMessagePath: /dev/termination-log
          terminationMessagePolicy: File
```

### `terminationMessagePath: /dev/termination-log`

A file inside the container where the application can write a short final message before exiting:

```sh
echo "Configuration validation failed" > /dev/termination-log
exit 1
```

After the container exits, the kubelet reads the file and places its contents into the Pod status:

```yaml
status:
  containerStatuses:
    - state:
        terminated:
          message: Configuration validation failed
```

This is separate from normal `stdout`/`stderr` container logs. Kubernetes does not automatically copy all application logs into this file.

### `terminationMessagePolicy: File`

Tells the kubelet to use only the contents of `terminationMessagePath` for the termination message, whether the container succeeded or failed.

The alternative is:

```yaml
terminationMessagePolicy: FallbackToLogsOnError
```

That policy uses the end of the container logs when:

- The termination-message file is empty, and
- The container exited with an error.

Termination messages are intended to be brief; the kubelet truncates oversized messages. [Kubernetes termination-message documentation](https://kubernetes.io/docs/tasks/debug/debug-application/determine-reason-pod-failure/)

## Pod infrastructure settings

### `dnsPolicy: ClusterFirst`

The kubelet configures the Pod so DNS queries go through the cluster DNS service, normally CoreDNS.

This enables names such as:

```text
my-service
my-service.default
my-service.default.svc.cluster.local
```

Queries outside the cluster domain are forwarded by cluster DNS to an upstream DNS server.

Despite its name, `Default` is not Kubernetes’ default Pod DNS policy—`ClusterFirst` is. [Kubernetes Pod DNS policies](https://kubernetes.io/docs/concepts/services-networking/dns-pod-service/#pod-s-dns-policy)

### `schedulerName: default-scheduler`

Selects which scheduler is responsible for assigning the Pod to a node.

The sequence is:

```text
Deployment controller
    → creates ReplicaSet
        → creates Pod
            → default-scheduler chooses Node
                → kubelet runs Pod
```

If omitted, the API server normally defaults it to `default-scheduler`.

Clusters can run custom schedulers:

```yaml
schedulerName: gpu-specialized-scheduler
```

If no running scheduler recognizes the specified name, the Pod normally remains `Pending` without being assigned to a node. [Kubernetes scheduler configuration](https://kubernetes.io/docs/reference/scheduling/config/#multiple-profiles)

## Deployment condition timestamps

These timestamps belong to one particular entry under:

```yaml
status:
  conditions:
    - type: Progressing
```

### `lastTransitionTime`

The last time that condition’s `status` changed between:

```text
True ↔ False ↔ Unknown
```

For your `Progressing` condition:

```yaml
lastTransitionTime: "2026-07-15T15:27:26Z"
status: "True"
```

The condition became `True` when the rollout began and remained `True`.

### `lastUpdateTime`

The last time the controller updated that condition, even when its Boolean status did not change.

In your example:

```yaml
lastTransitionTime: "2026-07-15T15:27:26Z"
lastUpdateTime: "2026-07-15T15:27:39Z"
reason: NewReplicaSetAvailable
```

The timeline is:

```text
15:27:26  Progressing became True
15:27:39  It was still True, but the controller updated it because
          the new ReplicaSet became available
```

Therefore:

- Different timestamps: the condition stayed in the same status but received newer information.
- Same timestamps: the update also changed its status, as happened for your `Available=True` condition.

These timestamps are condition-specific; every item in `status.conditions` has its own timestamps. The Deployment controller uses condition information for rollout/progress tracking, while `kubectl`, dashboards, alerts, and monitoring systems use it to explain the current state. [Deployment condition API reference](https://kubernetes.io/docs/reference/kubernetes-api/apps/deployment-v1/#DeploymentCondition)
