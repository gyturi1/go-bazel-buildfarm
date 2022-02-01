FROM golang:1.17-buster

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go install -ldflags "-s -w -extldflags '-static'" github.com/go-delve/delve/cmd/dlv@v1.7.2

CMD [ "/bin/bash" ]