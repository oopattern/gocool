# -- 外网如何访问容器 -- #
# ref: https://www.codenong.com/cs106463217/
# ref: https://www.jianshu.com/p/ee4ee15d3658
# docker run -d --name vm1 -p 80:80 nginx做端口映射（冒号 后的是容器内部的端口）
# docker port vm1查看容器端口映射情况 80/tcp->0.0.0.0:80
# iptables -t nat -S 查看防火墙策略
# ip addr show
# curl localhost 测试访问
# netstat -atnulp每当运行一个容器，就会开启一个docker-proxy进程
# docker run -d --name vm2 -p 808:80 nginx看到有2个docker-proxy进程

# -- Devcloud proxy settings -- #
# wget http://download.devcloud.oa.com/enable_internet_proxy.sh
# source  enable_internet_proxy.sh

# -- Traditional deployment -- #
# go build -o stringtest
# nohup ./stringtest &
# ps aux | grep stringtest | grep -v grep

# care about delete libc.so
# sln /lib64/libc-2.17.so /lib64/libc.so.6
# https://www.linuxidc.com/Linux/2017-02/140994.html

# -- Docker deployment -- #
# docker build . -t gokit-stringtest --network=bridge
# docker build . -t gokit-stringtest --network=host
# docker build . -t gokit-stringtest --network=default
# kubectl delete deployment stringtest
# kubectl get po --all-namespaces | awk '{ print $1,$2 }' | grep stringtest | awk '{print $2 }' | xargs kubectl delete po
# docker ps -a | awk '{ print $1,$2 }' | grep stringtest | awk '{print $1 }' | xargs docker stop
# docker rmi -f gokit-stringtest
# docker build . -t gokit-stringtest
# docker run -itd -p 8005:8005 gokit-stringtest
# docker ps | grep stringtest
# kubectl apply -f k8s-deploy.yaml
# kubectl get deployments
# kubectl get pods

# -- Show logs -- #
# docker logs -f <container-id>

# -- Start minikube -- #
# ref: https://github.com/kubernetes/minikube/issues/7903
# adduser k8s
# su - k8s
# groupadd docker
# usermod -aG docker k8s
# minikube start --driver=docker --cni=bridge
# docker ps

# -- Start minikube dashboard -- #
# debug: minikube dashboard --alsologtostderr -v=1
# ref: https://github.com/kubernetes/minikube/issues/5815

# -- Find local image with minikube -- #
# ref1: https://stackoverflow.com/questions/42564058/how-to-use-local-docker-images-with-minikube
# ref2: https://medium.com/swlh/how-to-run-locally-built-docker-images-in-kubernetes-b28fbc32cc1d
# ref3(k8s blog): https://kubernetes.io/blog/2019/03/28/running-kubernetes-locally-on-linux-with-minikube-now-with-kubernetes-1.14-support/

# -- Setup a local registry -- #
# https://github.com/kubernetes/website/issues/3596
# https://github.com/kubernetes/minikube/issues/1442
# https://minikube.sigs.k8s.io/docs/handbook/vpn_and_proxy/
# https://kubernetes.io/docs/tasks/administer-cluster/dns-debugging-resolution/
# https://kubernetes.io/zh/docs/tasks/administer-cluster/dns-debugging-resolution/
# https://github.com/kubernetes/minikube/issues/2350
# https://github.com/kubernetes/minikube/issues/1224
# https://segmentfault.com/a/1190000015639327
# https://stackoverflow.com/questions/42564058/how-to-use-local-docker-images-with-minikube
# config file: /home/k8s/.minikube/profiles/minikube/config.json
# kubectl logs --namespace=kube-system -l k8s-app=kube-dns
# kubectl describe pod coredns-f9fd979d6-lkrl6 -n kube-system
# kubectl -n kube-system delete pod -l k8s-app=kube-dns
# kubectl exec -it --namespace=kube-system etcd-minikube -- /bin/sh
# kubectl get po --all-namespaces -o wide
# kubectl get --all-namespaces svc
# kubectl get ep kube-dns --namespace=kube-system
# minikube start --alsologtostderr -v=1
# minikube start --driver=docker --cni=bridge
# eval $(minikube docker-env)
# docker inspect -f '{{.Name}} - {{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $(docker ps -aq)
# docker build -t foo:0.0.1 .
# kubectl run hello-foo --image=foo:0.0.1 --image-pull-policy=Never
# kubectl get pods

# -- Use a local registry -- #
# docker run -d -p 5500:5500 --restart=always --name registry registry:2
# docker tag gokit-stringtest localhost:5500/gokit-stringtest
# docker push localhost/gokit-stringtest
# docker pull localhost/gokit-stringtest


# -- Uninstall Minikube -- #
# minikube stop; minikube delete
# docker stop (docker ps -aq)
# rm -r ~/.kube ~/.minikube
# sudo rm /usr/local/bin/localkube /usr/local/bin/minikube
# systemctl stop '*kubelet*.mount'
# sudo rm -rf /etc/kubernetes/
# docker system prune -af --volumes
