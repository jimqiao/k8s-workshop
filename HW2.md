1. 构建本地镜像

root@cks-master:~# git clone https://github.com/jimqiao/k8s-workshop.git
Cloning into 'k8s-workshop'...
remote: Enumerating objects: 31, done.
remote: Counting objects: 100% (31/31), done.
remote: Compressing objects: 100% (27/27), done.
remote: Total 31 (delta 8), reused 2 (delta 0), pack-reused 0
Unpacking objects: 100% (31/31), done.
Checking connectivity... done.

root@cks-master:~# cd k8s-workshop

root@cks-master:~/k8s-workshop# docker build -t jimqiao/httpserver:v1.0 .
Sending build context to Docker daemon  88.58kB
Step 1/10 : FROM golang:alpine as build
 ---> 35cd8c8897b1
Step 2/10 : ENV GO111MODULE=on
 ---> Running in 0a363795ce04
Removing intermediate container 0a363795ce04
 ---> 21e77c60fd0d
Step 3/10 : ENV GOPROXY=https://goproxy.cn,direct
 ---> Running in 2376186fdaef
Removing intermediate container 2376186fdaef
 ---> 227fbd518e47
Step 4/10 : WORKDIR /usr/src
 ---> Running in e441da00ad84
Removing intermediate container e441da00ad84
 ---> c5e8d60e9516
Step 5/10 : COPY httpserver.go /usr/src
 ---> 58a94341b056
Step 6/10 : RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -ldflags="-s -w" -installsuffix cgo -o httpserver httpserver.go
 ---> Running in 12befc772e02
Removing intermediate container 12befc772e02
 ---> 1c7bdbd89126
Step 7/10 : FROM scratch as prod
 ---> 
Step 8/10 : COPY --from=build /usr/src/httpserver /
 ---> Using cache
 ---> 556df6916dbc
Step 9/10 : EXPOSE 80
 ---> Using cache
 ---> babfd5a9a6b9
Step 10/10 : ENTRYPOINT ["/httpserver"]
 ---> Using cache
 ---> 41111ff0d68d
Successfully built 41111ff0d68d
Successfully tagged jimqiao/httpserver:v1.0



2. 将镜像推送至 Docker 官方镜像仓库

root@cks-master:~/k8s-workshop# docker login
Authenticating with existing credentials...
WARNING! Your password will be stored unencrypted in /root/.docker/config.json.
Configure a credential helper to remove this warning. See
https://docs.docker.com/engine/reference/commandline/login/#credentials-store

Login Succeeded

root@cks-master:~/k8s-workshop# docker push jimqiao/httpserver:v1.0
The push refers to repository [docker.io/jimqiao/httpserver]
66ba0c8055f1: Layer already exists 
v1.1: digest: sha256:8e528130cfc69f6aeaef1b74fc89705f6d6f960f44717a853eff850670d71167 size: 528



3. 通过 Docker 命令本地启动 httpserver

root@cks-master:~/k8s-workshop# docker run --name=web -d -p 80:80 jimqiao/httpserver:v1.0
d548bba765dda49b3604c9ab5b0ba9220b73ccb0babe46778613a5e276cadd6a

root@cks-master:~/k8s-workshop# docker ps 
CONTAINER ID        IMAGE                     COMMAND             CREATED              STATUS              PORTS                NAMES
d548bba765dd        jimqiao/httpserver:v1.0   "/httpserver"       About a minute ago   Up About a minute   0.0.0.0:80->80/tcp   web

root@cks-master:~/k8s-workshop# curl http://127.0.0.1/healthz
200



4. 通过 nsenter 进入容器查看 IP 配置

root@cks-master:~/k8s-workshop# docker ps
CONTAINER ID        IMAGE                     COMMAND             CREATED             STATUS              PORTS                NAMES
d548bba765dd        jimqiao/httpserver:v1.0   "/httpserver"       20 minutes ago      Up 20 minutes       0.0.0.0:80->80/tcp   web

root@cks-master:~/k8s-workshop# docker inspect --format '{{ .State.Pid }}' d548bba765dd
16011

root@cks-master:~/k8s-workshop# nsenter -t 16011 -n ip addr
1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
    link/loopback 00:00:00:00:00:00 brd 00:00:00:00:00:00
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
24: eth0@if25: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc noqueue state UP group default 
    link/ether 02:42:ac:11:00:02 brd ff:ff:ff:ff:ff:ff link-netnsid 0
    inet 172.17.0.2/16 brd 172.17.255.255 scope global eth0
       valid_lft forever preferred_lft forever
