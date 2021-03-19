
# 1. k8s安装
https://github.com/ubuntu/microk8s
https://microk8s.io/docs

# 2. k8s启动, 生产环境
ps aux | grep kube | grep -v grep

# 3. 常用命令
# microk8s.kubectl get all --all-namespaces // 查看所有资源
# microk8s.kubectl -n kube-system get secret // 列出所有授权token
# microk8s.kubectl -n kube-system describe secret kubernetes-dashboard-token-9hl6q // 查看dashboard的token
# microk8s.kubectl cluster-info // 查看集群信息
# microk8s dashboard-proxy // 运行代理dashboard, 在浏览器输入部署的机器外网ip, 然后在浏览器输入验证的token

# 4. 编译构建k8s
#