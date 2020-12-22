# Commands

```
# Config
docker version
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
docker stop a33d19d68ced
```

