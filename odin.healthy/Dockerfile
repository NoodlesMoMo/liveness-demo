
FROM golang:1.11.5-alpine3.8 AS build

RUN mkdir -p $GOPATH/src/odin.healthy && mkdir -p /root/myapp

COPY . $GOPATH/src/odin.healthy

WORKDIR $GOPATH/src/odin.healthy

RUN go build -o /root/myapp/healthy

RUN echo -n "123456789" > /tmp/abc.txt

EXPOSE 8080

CMD /root/myapp/healthy

