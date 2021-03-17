
# 1. consul安装

# 2.1 consul启动, -dev为开发模式(非生产环境), 数据存储在内存中
nohup consul agent -dev -client=0.0.0.0 -data-dir=/data/tools/consul/data -node=consul-sakulali -log-file=/data/tools/consul/log/consul.log &

# 2.2 consul启动, 没有dev参数(生产环境), 数据存储在目录文件中
nohup consul agent -client=0.0.0.0 -data-dir=/data/tools/consul/data -node=consul-sakulali -log-file=/data/tools/consul/log/consul.log -advertise=9.134.118.145 -bootstrap-expect=1 -server=true -join=9.134.118.145 -ui=true &

ps aux | grep consul | grep -v grep | grep -v start_consul
