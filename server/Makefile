build:
	docker build -t echo/server .

push:
	aws ecr get-login-password --region ap-northeast-1 | docker login --username AWS --password-stdin 373656256964.dkr.ecr.ap-northeast-1.amazonaws.com
	docker tag echo/server:latest 373656256964.dkr.ecr.ap-northeast-1.amazonaws.com/echo/server:latest
	docker push 373656256964.dkr.ecr.ap-northeast-1.amazonaws.com/echo/server:latest
run:
	docker run -it -p 50051:50051 echo/server

# The old way
# refresh: build push
# 	aws ecs update-service --cluster echo --service echo-app-EchoServerService-1MWS12Z1FZXEO --force-new-deployment

update: build push
	./update-service.sh

start: build run
