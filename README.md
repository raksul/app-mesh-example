

# Setup

## Infrastructure setup

First, make sure you copy the following file and edit its contents:
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

Once you have the .env setup, run the following script to initialize VPC, ECR/ECS and App Mesh.
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

Access the above load balancer and make sure that you have output like this:
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

# Behind the scene

Here is the diagram that explains the infrastructure:

![diagram](https://github.com/raksul/app-mesh-example/raw/master/doc/infrastructure.jpg)

## The data flow

* From the Web Browser, App Load Balancer accepts a HTTP request on port 80.
* App Load Balacer forwards the HTTP request to one of the `echo_client` services.
* `echo_client` service then calls `echo_server` through gRPC, proxying a virtual service and virtual router.
* V.Router's routing forwards the traffic to one of the healthy virtual node.

# Updating the server code

Now, knowing that the `echo_client` is sucessfully communicating with `echo_service` via gRPC, let's update `echo_service` to a newer version.

## Update the server-side code:

Let's just update the version indicated in `server.go`:
```go:server/server.go
- const version = "0.9"
+ const version = "1.0"
```

Test the code by running the server on local environment.
```bash:Screen 1
$ go run server.go
2020/06/18 23:37:32 Echo Server Version 1.0
2020/06/18 23:37:32 Starting on port: 50051, ssl: false
```

Run the client in a separate console:
```bash:Screen 2
$ cd client
$ go run client.go
2020/06/18 23:38:30 Echo Client 0.1
2020/06/18 23:38:30 Echo Client connecting to localhost:50051
2020/06/18 23:38:30 Listening to port: 8080
```

Finally, call the client via HTTP:
```bash:Screen 3
$ curl localhost:8080/?name=Yizumi
Response from the server: Hello, Yizumi-san! (Said 192.168.11.8, Version 1.0)
```
Now you see that the server is responding with Version 1.0.

## Refresh the ECS Service

Let's now see if we can update the `echo_server` seamlessly without breaking the client.

First, let's run the following `curl` command to keep the request running during the update:
```bash:Terminal 1
$ while true; do echo; curl "http://echo-publi-upxd2gg2nc16-261471328.ap-northeast-1.elb.amazonaws.com/?name=Yusuke"; done
Response from the server: Hello, Yusuke-san! (Said 169.254.172.42, Version 0.9)
Response from the server: Hello, Yusuke-san! (Said 169.254.172.42, Version 0.9)
Response from the server: Hello, Yusuke-san! (Said 169.254.172.42, Version 0.9)
...
```

Then update the by running the update command.
```bash:Terminal 2
$ cd server
$ make update
```

This will build a new docker image, pushes it to ECS, then creates a new ECS service, wait for the new service to be stable, then switches traffic from old instances to new instances.

At one point you'll see something like this:
```bash:Terminal 2
Waiting for ECS Service to be in RUNNING state...
Tasks are starting (0/2)...
Tasks are starting (0/2)...
Tasks are starting (1/2)...
Tasks started
Updating traffic route
Routing 50% of traffic to the new service
Updating traffic route
```

Note that in Terminal 1, you'll begin 50% of traffic being forward to the newer service:
```bash:Terminal 1
Response from the server: Hello, Yusuke-san! (Said 169.254.172.42, Version 0.9)
Response from the server: Hello, Yusuke-san! (Said 169.254.172.42, Version 1.0)
Response from the server: Hello, Yusuke-san! (Said 169.254.172.42, Version 1.0)
Response from the server: Hello, Yusuke-san! (Said 169.254.172.42, Version 0.9)
Response from the server: Hello, Yusuke-san! (Said 169.254.172.42, Version 1.0)
Response from the server: Hello, Yusuke-san! (Said 169.254.172.42, Version 1.0)
Response from the server: Hello, Yusuke-san! (Said 169.254.172.42, Version 0.9)
Response from the server: Hello, Yusuke-san! (Said 169.254.172.42, Version 0.9)
Response from the server: Hello, Yusuke-san! (Said 169.254.172.42, Version 1.0)
```

Finally, the script will update the routing to send all traffic to the new version.
```bash:Terminal 2
Routing 100% of traffic to the new service
```

Note the change in Terminal 1:
```bash:Terminal 1
Response from the server: Hello, Yusuke-san! (Said 169.254.172.42, Version 1.0)
Response from the server: Hello, Yusuke-san! (Said 169.254.172.42, Version 1.0)
Response from the server: Hello, Yusuke-san! (Said 169.254.172.42, Version 1.0)
Response from the server: Hello, Yusuke-san! (Said 169.254.172.42, Version 1.0)
Response from the server: Hello, Yusuke-san! (Said 169.254.172.42, Version 1.0)
Response from the server: Hello, Yusuke-san! (Said 169.254.172.42, Version 1.0)
Response from the server: Hello, Yusuke-san! (Said 169.254.172.42, Version 1.0)
```

The Server 0.9 can now be safely destroyed. (WIP)

## Summary

This deploy strategy is very similar to Blue-Green Deploy available in ECS, however, it clearly works with gRPC, and does not break the connection between client and server.

Let's look at the diagram to see what just happened behind the scene.

![diagram](https://github.com/raksul/app-mesh-example/raw/master/doc/service-update.jpg)

In hindsight, this is what happened:
* The requests coming from `echo_client` was served by the Virtual Service, which has a routing definition saying "I'm routing 100% traffic to version 0.9"
* Then you uploaded the Version 1.0 of `calc_server`. We'll let couple seconds pass until the instaces become stable. (Figure 1)
* Once confirmed that Version 1.0 instances are running in "HEALTHY" state, we'll begin routing 50% of the traffic to Version 1.0 instances. (Figure 2)
* It then waits couple more seconds and begin routing ALL requests to Version 1.0 instances.
* Once confirmed that all traffic now routes to the new service, we can safely destroy the old node.
