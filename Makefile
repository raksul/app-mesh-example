
generate:
	protoc ./echo.proto --go_out=plugins=grpc:./server
	protoc ./echo.proto --go_out=plugins=grpc:./client

build-infra:
	. ./.env
	./deploy.sh
