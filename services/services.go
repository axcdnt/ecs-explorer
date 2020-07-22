package services

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/fatih/color"
)

var (
	activeStatus           = color.New(color.FgGreen)
	deactivatingStatus     = color.New(color.FgYellow)
	inactiveStatus         = color.New(color.FgMagenta)
	stoppedStatus          = color.New(color.FgHiRed)
	containsErrorStatus    = color.New(color.FgHiRed)
	activeNotRunningStatus = color.New(color.FgHiYellow)
)

// EcsService represents an ECS service
type EcsService struct {
	Cluster  *string
	Services []*string
	Session  *session.Session
}

// Query list services according to params
func (e *EcsService) Query(resp chan<- string) {
	svcDescr := e.describe()
	for _, failure := range svcDescr.Failures {
		reportFailure(failure, resp)
	}
	for _, svc := range svcDescr.Services {
		reportStatus(svc, resp)
	}
}

func reportFailure(failure *ecs.Failure, resp chan<- string) {
	resp <- containsErrorStatus.Sprintf("Failure on fetch service info. Arn: %s. Reason: %s\n", *failure.Arn, *failure.Reason)
}

func reportStatus(service *ecs.Service, resp chan<- string) {
	name := *service.ServiceName
	status := *service.Status
	desired := *service.DesiredCount
	running := *service.RunningCount

	// task definition stuff
	desiredRevision := taskRevision(*service.TaskDefinition)
	latestRevision := taskRevision(*service.Deployments[0].TaskDefinition)

	message := fmt.Sprintf(
		"%s: status %s, desired: %d, running: %d, desired revision: %s, latest running revision: %s",
		name, status, desired, running, desiredRevision, latestRevision,
	)

	if len(service.Deployments) == 1 && isLatestRevisionRunning(desiredRevision, latestRevision) && desired == running {
		resp <- activeStatus.Sprintln(message)
	} else if status == "DEACTIVATING" {
		resp <- deactivatingStatus.Sprintln(message)
	} else if status == "INACTIVE" {
		resp <- inactiveStatus.Sprintln(message)
	} else if status == "STOPPED" {
		resp <- stoppedStatus.Sprintln(message)
	} else {
		resp <- activeNotRunningStatus.Sprintln(message)
	}
}

func taskRevision(taskDefinition string) string {
	separatorIndex := strings.LastIndex(taskDefinition, ":")
	return taskDefinition[separatorIndex+1:]
}

func isLatestRevisionRunning(desiredRevision, latestRevision string) bool {
	return desiredRevision == latestRevision
}

func (e *EcsService) describe() *ecs.DescribeServicesOutput {
	svc := ecs.New(e.Session)
	config := newConfig(e.Cluster, e.Services)
	result, err := svc.DescribeServices(config)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case ecs.ErrCodeServerException:
				fmt.Println(ecs.ErrCodeServerException, aerr.Error())
			case ecs.ErrCodeClientException:
				fmt.Println(ecs.ErrCodeClientException, aerr.Error())
			case ecs.ErrCodeInvalidParameterException:
				fmt.Println(ecs.ErrCodeInvalidParameterException, aerr.Error())
			case ecs.ErrCodeClusterNotFoundException:
				log.Fatalln("an error occurred while trying to find cluster: cluster not found")
			case ecs.ErrCodeAccessDeniedException:
				// this is possibly a frequent error
				log.Fatalln("an error occurred while trying to access aws: invalid credentials or related")
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return nil
	}

	return result
}

func newConfig(cluster *string, services []*string) *ecs.DescribeServicesInput {
	return &ecs.DescribeServicesInput{
		Cluster:  cluster,
		Services: services,
	}
}
