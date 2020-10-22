# ecs-explorer
[![CircleCI Build Status](https://circleci.com/gh/axcdnt/ecs-explorer/tree/master.svg?style=shield)](https://circleci.com/gh/axcdnt/ecs-explorer/tree/master)

## What is ecs-explorer?
This is a small project to simplify read-only tasks for an AWS ECS cluster.
If you, like me, have multiple clusters and from time to time need to check services _availability_, this might help you.


## Motivation
The idea behind ecs-explorer is to give a flexible and colorized output perspective of your cluster. 
I created it because I enjoy CLI apps and I wanted to explore [AWS ECS SDK](https://docs.aws.amazon.com/sdk-for-go/api/index.html).

I also used the codebase as an instructional code session to explain Go for a group of friends.

## How to

Build: 

```
go build -o ecs-explorer main.go 
```

Usage:

```
▶ ./ecs-explorer --help
Usage of ./ecs-explorer:
  -cluster string
    	the cluster name
  -services string
    	a comma separated list of services
  -suffix string
    	the service name suffix to look for
```

Output:

```
▶ ./ecs-explorer --cluster=ecs-<cluster-name> --services=<service1,service2,service3> --suffix=-foo

service1-foo: status ACTIVE, desired: 1, running: 1
service2-foo: status ACTIVE, desired: 1, running: 1
service3-foo: status INACTIVE, desired: 0, running: 0
```

I hope you enjoy it!
