# Commands

```
# Config
docker version
docker login
docker ps -a

# Run image
docker search docker/whalesay
docker run hello-world
docker run -it ubuntu-server /bin/bash
docker images

# Build docker image & stop
docker build -t docker-whale .
docker run -d -p 80:80 go-server
docker ps
docker kill a33d19d68ced

# Remove image or ALL images
docker rmi -f go-server
docker system prune // remove all unused containers
docker system prune -a // remove all images 

# Push on dockerhub / pull
docker tag go-server kevinjose/go-server:latest
docker run -d -p 80:80  kevinjose/go-server
// docker pull kevinjose/go-server
```

