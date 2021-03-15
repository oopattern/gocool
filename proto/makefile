
.PHONY: build
build:
	protoc -I . --go_out=. --go-grpc_out=. --grpc-gateway_out=logtostderr=true:. ./*.proto

.PHONY: clean
clean:
	rm -f *pb.go
	rm -f *pb.gw.go
