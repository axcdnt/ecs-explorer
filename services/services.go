package services

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/fatih/color"
	"log"
)

var (
	activeStatus       = color.New(color.FgGreen)
	deactivatingStatus = color.New(color.FgYellow)
	inactiveStatus     = color.New(color.FgMagenta)
	stoppedStatus      = color.New(color.FgHiRed)
)

// EcsService represents an ECS service
type EcsService struct {
	Cluster  *string
	Services []*string
	Session  *session.Session
}

// Query list services according to params
func (e *EcsService) Query() {
	svcs := e.describe().Services
	for _, svc := range svcs {
		name := svc.ServiceName
		status := svc.Status
		desired := svc.DesiredCount
		running := svc.RunningCount

		prettyPrint(*name, *status, *desired, *running)
	}
}

// to be improved
func prettyPrint(serviceName, status string, desired, running int64){
	message := fmt.Sprintf("%s: status %s, desired: %d, running: %d", serviceName, status, desired, running)

	if status == "ACTIVE" {
		activeStatus.Println(message)
	}
	if status == "DEACTIVATING" {
		deactivatingStatus.Println(message)
	}
	if status == "INACTIVE" {
		inactiveStatus.Println(message)
	}
	if status == "STOPPED" {
		stoppedStatus.Println(message)
	}
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

