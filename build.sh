# -- Traditional deployment -- #
# go build -o gocool
# nohup ./gocool &
# ps aux | grep gocool | grep -v grep

# -- Docker deployment -- #
# docker build . -t dock-gocool --network=bridge
# docker run -itd -p 7777:7777 dock-gocool
docker ps -a | awk '{ print $1,$2 }' | grep gocool | awk '{print $1 }' | xargs docker stop
docker rmi -f dock-gocool
docker build . -t dock-gocool
docker run -itd dock-gocool
docker ps | grep dock-gocool

# -- Show logs -- #
# docker logs -f <container-id>