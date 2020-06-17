build:
	docker build -t echo/client .

push:
	aws ecr get-login-password --region ap-northeast-1 | docker login --username AWS --password-stdin 373656256964.dkr.ecr.ap-northeast-1.amazonaws.com
	docker tag echo/client:latest 373656256964.dkr.ecr.ap-northeast-1.amazonaws.com/echo-client:latest
	docker push 373656256964.dkr.ecr.ap-northeast-1.amazonaws.com/echo/client:latest

run:
	docker run -it -p 8080:8080 echo/client

start: build run
