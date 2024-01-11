# feed
News crawler

## Need

Feed is a news crawler.  The crawler downloads content from several RSS websites such as The Verge, Wired, Mashable, etc.  The crawler is triggered by CronJob, which makes gRPC call to store the content into MySQL database.  The CronJob and gRPC are containerized and its workload is managed in Kubernetes.  

I use this news crawler to accumulate news content into my desire format.  It allows me to archieve the data.  The infrastructure components such as Kubernetes and MySQL are hosted in my homelab, with the exception of Docker hub.

## Infrastructure

### Kubernetes homelab

I use [kubeadm](https://kubernetes.io/docs/reference/setup-tools/kubeadm/) to setup a Kubernetest cluster with master node and two worker nodes.

```
kubectl top nodes
NAME                    CPU(cores)   CPU%   MEMORY(bytes)   MEMORY%
k8smaster.edison.net    339m         2%     20483Mi         64%
k8sworker1.edison.net   153m         3%     3840Mi          54%
k8sworker2.edison.net   178m         2%     10143Mi         31%

kubectl get nodes
NAME                    STATUS   ROLES           AGE    VERSION
k8smaster.edison.net    Ready    control-plane   251d   v1.27.1
k8sworker1.edison.net   Ready    <none>          251d   v1.27.1
k8sworker2.edison.net   Ready    <none>          251d   v1.27.1
```

### Cronjob

The Cronjob is containerized and managed in Kubernetes.   
Link to [source code](https://github.com/elumbantoruan/feed/tree/main/cmd/cronjob).

```
kubectl get cronjobs
NAME               SCHEDULE      SUSPEND   ACTIVE   LAST SCHEDULE   AGE
newsfeed-cronjob   */5 * * * *   False     0        3m44s           6h41m
```

### gRPC

The gRPC is containerized and managed in Kubernetes.  
Link to [source code](https://github.com/elumbantoruan/feed/tree/main/cmd/grpc/server).
```
kubectl get services
NAME            TYPE       CLUSTER-IP     EXTERNAL-IP   PORT(S)          AGE
newsfeed-grpc   NodePort   10.97.150.39   <none>        9000:30008/TCP   6h50m
```

### pods
```
kubectl get pods
NAME                              READY   STATUS      RESTARTS   AGE
newsfeed-cronjob-28415755-p7m5j   0/1     Completed   0          12m
newsfeed-cronjob-28415760-q2c7n   0/1     Completed   0          7m6s
newsfeed-cronjob-28415765-cvnmr   0/1     Completed   0          2m6s
newsfeed-grpc-5fd97bfb69-k529n    1/1     Running     0          6h49m
newsfeed-grpc-5fd97bfb69-v9hgf    1/1     Running     0          6h49m
```

### MySQL

At this point, MySQL database is not managed in Kubernetes, rather as a local installation.  
Link to [database schema](https://github.com/elumbantoruan/feed/tree/main/pkg/storage/db-script).

### Docker hub
Link to [docker hub](https://hub.docker.com/repositories/edisonlt) for Cronjob and gRPC repositories.