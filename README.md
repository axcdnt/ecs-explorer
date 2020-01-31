# ecs-explorer

## What is ecs-explorer?
This is a small project to simplify read-only tasks from an ECS cluster perspective.
If you, like me, have multiple clusters and from time to time need to check services availability, this might help you.


## Motivation
The idea behind ecs-explorer is to give a flexible and colorized output perspective of your cluster. 
I created it because I enjoy CLI apps and I wanted to explore AWS ECS sdk.

It has served as an instructional code session to explain Go code for a few friends.

## How to

Usage:

```
▶ ./ecs-utils --help
Usage of ./ecs-utils:
  -cluster string
    	the cluster name (default "qa")
  -code string
    	the environment code (only needed for qa)
  -env string
    	the environment: qa/prod
  -services string
    	comma separated list of services (max of 10)
```

Output:

```
▶ ./ecs-utils --cluster=ecs-<cluster-name> --services=<service1,service2,service3> --env=<prod,qa>

service1: status ACTIVE, desired: 1, running: 1
service2: status ACTIVE, desired: 1, running: 1
service3: status INACTIVE, desired: 0, running: 0
```

Notice that, today, the `--env` param supports prod/qa as environments. It was my specific need, but can be easily adapted.

I hope you enjoy it!
