# kube2ram

English | [简体中文](./README-zh_CN.md)

## Context

In current Alibaba ACK cluster runtime container, if you want to communication with other Alibaba Cloud resource as oss, slb..., you should call ECS metadata API to fetch the interim security credentials for RAM authorization.
 
But in a multi-tenanted containers based world, different users will deloy multiple containers in the same worker node, since only one worker RAM role could be fetch through ECS instance profile, the different containers could only sharing the security credentials mapping to the same RAM role, which is not acceptable from the security isolation perspective.

kube2ram is just the solution based on this issue, it could deploy on each worker node as a daemonst, and redirect the tracffic that is going to the ECS metadata server from bussiness container to kube2ram instance, make a call to the ECS API to fetch interim credentials and return these to the caller. Other calls will be still proxied to the ECS metadata server.

<br/> 
<img src="kube2ram.jpg" alt="kube2ram">

## Configuration

### RAM roles

Call the metadata server http://100.100.100.200/latest/meta-data/ram/security-credentials/ on worker node to get the default RAM role name (prefix with `KubernetesWorkerRole`),  then open RAM console-> RAM roles page, search and find the role with the given name, then go to the RAM policies page and append the policy below in Policy Document tab.（In this case you grant the target KubernetesWorkerRole to assume any other RAM role, and you can also set one or more RAM role in the `Resource` vaule for fine-grained control）：

```json
{
  "Version": "1",
  "Statement": [
    ...//omit the granted policies here
    {
      "Action": [
        "sts:AssumeRole"
      ],
      "Effect": "Allow",
      "Resource": "*"
    }
  ]
}
```
use your Alibaba Cloud account uid in `acs:ram::xxxxxxx:root`

also you need to find all of the target assumeroles, and make sure the account uid exists in trust `Principal` as below: 

```json
{
    "Statement": [
        {
            "Action": "sts:AssumeRole",
            "Effect": "Allow",
            "Principal": {
                "RAM": [
                    "acs:ram::xxxxxxx:root"
                ]
            }
        }
    ],
    "Version": "1"
}
```
* xxxxxxx need replace by the Alibaba Cloud account uid who own this worker node.


### kube2ram daemonset

Deploy the kube2ram daemonset with the template below, please notice that the kube2ram daemon and iptables rule need to run before all other pods that would require access to Alibaba Cloud resources.

```yaml
apiVersion: apps/v1
kind: DaemonSet
metadata:
  labels:
    app.kubernetes.io/name: kube2ram
  name: kube2ram
  namespace: kube-system
spec:
  selector:
   matchLabels:
    app.kubernetes.io/name: kube2ram
  template:
    metadata:
      labels:
        app.kubernetes.io/name: kube2ram
    spec:
      containers:
      - name: kube2ram
        image: registry.cn-hangzhou.aliyuncs.com/acs/kube2ram:1.0.0
        imagePullPolicy: Always
        args:
          - "--app-port=8181"
          - "--iptables=true"
          - "--host-ip=$(HOST_IP)"
          - "--host-interface=cni0"
          - "--verbose"
          - "--auto-discover-default-role"
        env:
        - name: HOST_IP
          valueFrom:
            fieldRef:
              apiVersion: v1
              fieldPath: status.podIP
        - name: NODE_NAME
          valueFrom:
            fieldRef:
              fieldPath: spec.nodeName
        ports:
        - containerPort: 8181
          protocol: TCP
        securityContext:
          privileged: true
      hostNetwork: true
      serviceAccountName: kube2ram
```

### iptables

To prevent containers from directly accessing the ECS metadata API and gaining unwanted access to Alibaba Cloud resources, the traffic to `100.100.100.200` must be proxied for all user containers.

```bash
iptables \
  --append PREROUTING \
  --protocol tcp \
  --destination 100.100.100.200 \
  --dport 80 \
  --in-interface docker0 \
  --jump DNAT \
  --table nat \
  --to-destination `http://100.100.100.200/latest/meta-data/private-ipv4`:8181
```

This rule can be added automatically by setting `--iptables=true`, setting the `HOST_IP` environment variable, and running the container in a privileged security context.

**Warning**: It is possible that other pods are started on an instance before kube2ram has started. Using `--iptables=true` (instead of applying the rule before starting the kubelet) **could give those pods the opportunity to access the real ECS metadata API, assume the role of the ECS instance and thereby have all permissions the instance role has** (including assuming potential other roles). Use with care if you don't trust the users of your kubernetes cluster or if you are running pods (that could be exploited) that have permissions to create other pods (e.g. controllers / operators).

Note that the interface --in-interface above or using the --host-interface cli flag may be different than docker0 depending on which virtual network you use e.g.

* flannel use `cni0`
* [Terway](https://github.com/AliyunContainerService/terway) or Calico,use `cali+` (如：cali1234567890)
* kops (on kubenet), use `cbr0`
* CNI, use `cni0`
* weave use `weave`
* [kube-router](https://github.com/cloudnativelabs/kube-router) use `kube-bridge`
* [OpenShift](https://www.openshift.org/) use `tun0`
* [Cilium](https://www.cilium.io) use `lxc+`

### kubernetes annotation

Add an `ram.aliyuncs.com/role` annotation to your pods with the role that you want to assume in this pod.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx-deployment
  labels:
    app: nginx
spec:
  replicas: 3
  selector:
    matchLabels:
      app: nginx
  template:
    metadata:
      annotations:
        ram.aliyuncs.com/role: kube2ram-arn
      labels:
        app: nginx
    spec:
      containers:
      - name: nginx
        image: nginx:1.9.1
        ports:
        - containerPort: 80
```

You can use `--default-role` to set a fallback role to use when annotation is not set.

annotation in other kubernetes deploy model，e.g.`CronJob`:

```yaml
apiVersion: batch/v1
kind: CronJob
metadata:
  name: my-cronjob
spec:
  schedule: "00 11 * * 2"
  concurrencyPolicy: Forbid
  startingDeadlineSeconds: 3600
  jobTemplate:
    spec:
      template:
        metadata:
          annotations:
            ram.aliyuncs.com/role: kube2ram-arn
        spec:
          restartPolicy: OnFailure
          containers:
          - name: job
            image: my-image
```

### Namespace Restrictions

By using the flag --namespace-restrictions you can enable a mode in which the roles that pods can assume is restricted by an annotation on the pod's namespace. This annotation should be in the form of a json array.

To allow the ACK pod specified above to run in the default namespace your namespace would look like the following.

```yaml
apiVersion: v1
kind: Namespace
metadata:
  annotations:
    ram.aliyuncs.com/allowed-roles: |
      ["role-arn"]
  name: default
```

_Notice:_ to allow all roles prefixed with my-custom-path/ to be assumed by pods in the default namespace, the default namespace would be annotated as follows:

```yaml
apiVersion: v1
kind: Namespace
metadata:
  annotations:
    ram.aliyuncs.com/allowed-roles: |
      ["my-custom-path/*"]
  name: default
```

If you prefer regexp to glob-based matching you can specify --namespace-restriction-format=regexp, then you can use a regexp in your annotation:

```yaml
apiVersion: v1
kind: Namespace
metadata:
  annotations:
    ram.aliyuncs.com/allowed-roles: |
      ["my-custom-path/.*"]
  name: default
```

### RBAC

This is the basic RBAC setup to get kube2ram working correctly when your cluster is using rbac. Below is the bare minimum to get kube2ram working.

```yaml
---
apiVersion: v1
kind: ServiceAccount
metadata:
  name: kube2ram
  namespace: kube-system
```


```yaml
---
apiVersion: v1
items:
  - apiVersion: rbac.authorization.k8s.io/v1
    kind: ClusterRole
    metadata:
      name: kube2ram
    rules:
      - apiGroups: [""]
        resources: ["namespaces","pods"]
        verbs: ["get","watch","list"]
  - apiVersion: rbac.authorization.k8s.io/v1
    kind: ClusterRoleBinding
    metadata:
      name: kube2ram
    subjects:
    - kind: ServiceAccount
      name: kube2ram
      namespace: kube-system
    roleRef:
      kind: ClusterRole
      name: kube2ram
      apiGroup: rbac.authorization.k8s.io
kind: List
```

### Debug

By using the --debug flag you can enable some extra features making debugging easier:

* `/debug/store` endpoint enabled to dump knowledge of namespaces and role association.

### Base ARN auto discovery

By using flag `--auto-discover-base-arn`，kube2ram will auto discovery the base ARN via ECS metadata server.

### Using KubernetesWoker instance role as default role

By using flag `--auto-discover-default-role`，kube2ram will auto discovery the ARN of the RAM role attached to ECS and use it as the fallback role to use when annotation is not set.


### Options

```bash
$ kube2ram --help
Usage of kube2ram:
      --api-server string                     Endpoint for the api server
      --api-token string                      Token to authenticate with the api server
      --app-port string                       Kube2iam server http port (default "8181")
      --auto-discover-base-arn                Queries ECS Metadata to determine the base ARN
      --auto-discover-default-role            Queries ECS Metadata to determine the default RAM Role and base ARN, cannot be used with --default-role, overwrites any previous setting for --base-role-arn
      --backoff-max-elapsed-time duration     Max elapsed time for backoff when querying for role. (default 2s)
      --backoff-max-interval duration         Max interval for backoff when querying for role. (default 1s)
      --base-role-arn string                  Base role ARN
      --ram-role-session-ttl                  Length of session when assuming the roles (default 15m)
      --debug                                 Enable debug features
      --default-role string                   Fallback role to use when annotation is not set
      --host-interface string                 Host interface for proxying ECS metadata (default "docker0")
      --host-ip string                        IP address of host
      --ram-role-key string                   Pod annotation key used to retrieve the RAM role (default "ram.aliyuncs.com/role")
      --iptables                              Add iptables rule (also requires --host-ip)
      --log-format string                     Log format (text/json) (default "text")
      --log-level string                      Log level (default "info")
      --metadata-addr string                  Address for the ECS metadata (default "100.100.100.200")
      --metrics-port string                   Metrics server http port (default: same as kube2ram server port) (default "8181")
      --namespace-key string                  Namespace annotation key used to retrieve the RAM roles allowed (value in annotation should be json array) (default "ram.aliyuncs.com/allowed-roles")
      --namespace-restriction-format string   Namespace Restriction Format (glob/regexp) (default "glob")
      --namespace-restrictions                Enable namespace restrictions
      --node string                           Name of the node where kube2ram is running
      --verbose                               Verbose
      --version                               Print the version and exits
``` 