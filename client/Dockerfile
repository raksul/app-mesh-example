FROM golang:1.14.2
MAINTAINER yizumi

RUN apt-get update && apt-get install -y unzip
WORKDIR /go/src/github.com/raksul/echo

# make protoc available for ... what?
RUN mkdir -p /tmp/protoc && \  
  curl -L https://github.com/protocolbuffers/protobuf/releases/download/v3.10.0/protoc-3.10.0-linux-x86_64.zip > /tmp/protoc/protoc.zip && \  
  cd /tmp/protoc && \  
  unzip protoc.zip && \  
  cp /tmp/protoc/bin/protoc /usr/local/bin && \  
  chmod go+rx /usr/local/bin/protoc && \  
  cd /tmp && \  
  rm -r /tmp/protoc

RUN go get -u google.golang.org/grpc
RUN go get -u github.com/golang/protobuf/protoc-gen-go

ENV PROJECT_HOME=$GOPATH/src/github.com/raksul/app-mesh-example/client

WORKDIR $PROJECT_HOME 
COPY . $PROJECT_HOME

EXPOSE 8080

CMD go run $PROJECT_HOME/client.go
