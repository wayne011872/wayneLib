VERSION=`git describe --tags`
BUILD_TIME=`date +%FT%T%z`
LDFLAGS=-ldflags "-X main.Version=${V} -X main.BuildTime=${BUILD_TIME}"

gen-code:
	protoc --go_out=. --go-grpc_out=. ${SER}/proto/*.proto
