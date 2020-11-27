# -- Traditional deployment -- #
cd /data/sakulali/go_workspace/src/gocool/server
killall gocool-server
killall gocool-server
go build -o gocool-server
nohup ./gocool-server &
ps aux | grep gocool-server | grep -v grep

sleep 1

cd /data/sakulali/go_workspace/src/gocool/client
killall gocool-client
go build -o gocool-client
ps aux | grep gocool-client | grep -v grep
./gocool-client

# -- Docker deployment -- #
# docker build . -t dock-gocool --network=bridge
# docker run -itd -p 7777:7777 dock-gocool
# docker ps -a | awk '{ print $1,$2 }' | grep gocool | awk '{print $1 }' | xargs docker stop
# docker rmi -f dock-gocool
# docker build . -t dock-gocool
# docker run -itd dock-gocool
# docker ps | grep dock-gocool

# -- Show logs -- #
# docker logs -f <container-id>