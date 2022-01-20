#!/bin/bash
#!/bin/sh
CURRENTPATH=$(pwd)
cd ../../../../../..
export GOPATH=$(pwd)
cd $CURRENTPATH
rm -fr *.go
cd protoc || exit
protoc --experimental_allow_proto3_optional\
  -I . \
  -I ${GOPATH}/src/github.com/protocolbuffers/protobuf/src/google/protobuf \
  -I ${GOPATH}/src/github.com/JiangNan7Guai/protoc-gen-validate \
  -I ${GOPATH}/src/github.com/protocolbuffers/protobuf/src \
  -I ${GOPATH}/src/github.com/googleapis/googleapis \
   \
  --go_out=.. \
  --go-grpc_out=.. \
  --grpc-gateway_out=.. \
  *.proto