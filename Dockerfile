FROM golang:alpine as build
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct
WORKDIR /usr/src
COPY httpserver.go /usr/src
RUN GOOS=linux CGO_ENABLED=0 GOARCH=amd64 go build -ldflags="-s -w" -installsuffix cgo -o httpserver httpserver.go

FROM scratch as prod
COPY --from=build /usr/src/httpserver /
EXPOSE 80
ENTRYPOINT ["/httpserver"]
