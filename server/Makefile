build:
	docker build -t echo/server .

push:
	if [ -z $ECR_IMAGE_PREFIX ]; then
		echo "ECR_IMAGE_PREFIX environment variable is not set."
		exit 1
	fi
	aws ecr get-login-password --region ${AWS_DEFAULT_REGION} | docker login --username AWS --password-stdin ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_DEFAULT_REGION}.amazonaws.com
	docker tag echo/server:latest ${ECR_IMAGE_PREFIX}/server:latest
	docker push ${ECR_IMAGE_PREFIX}/server:latest
run:
	docker run -it -p 50051:50051 echo/server

# The old way
# refresh: build push
# 	aws ecs update-service --cluster echo --service echo-app-EchoServerService-1MWS12Z1FZXEO --force-new-deployment

update: build push
	. ../.env
	./update-service.sh

start: build run
