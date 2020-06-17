

# Setup

## Infrastructure setup

Before running the following script, make sure you copy the following file and edit its contents:
```
$ cp .env-sample .env
$ vim .env
```
Variables:

* `PROJECT_NAME` Any arbitrary project name.  Use 'echo' if you don't have any preference.
* `AWS_DEFAULT_REGION` Your preferred AWS region
* `AWS_ACCOUNT_ID` Your account ID as you see [here](https://console.aws.amazon.com/billing/home?#/account)
* `ENVOY_IMAGE` Your preferred Envoy image. No need to edit.
* `KEY_PAIR` Name of Key Pair you'd like to use to setup the infrastructure. Find it [here](https://ap-northeast-1.console.aws.amazon.com/ec2/v2/home?region=ap-northeast-1#KeyPairs)

Once you have the .env setup, simply run the following script to initialize VPC, ECR/ECS and App Mesh.
```
$ make build-infra
```

When done, you should have listed:
```
Bastion endpoint:
54.65.206.60
Public endpoint:
http://echo-Publi-UPXD2GG2NC16-261471328.ap-northeast-1.elb.amazonaws.com
```

Go access the above load balancer and make sure that you have output like this:
```
$ curl -i "http://echo-publi-upxd2gg2nc16-261471328.ap-northeast-1.elb.amazonaws.com/?name=Yusuke"
HTTP/1.1 200 OK
Date: Wed, 17 Jun 2020 14:57:37 GMT
Content-Type: text/plain; charset=utf-8
Content-Length: 79
Connection: keep-alive
x-envoy-upstream-service-time: 2
server: envoy

Response from the server: Hello, Yusuke-san! (Said 169.254.172.42, Version 0.9)
```

Congrats! You're accessing the web service supported by gRPC using App Mesh!

# What's happening behind the scene?

Here is the diagram that explains the infrastructure:

