package main

import (
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecs"
	"github.com/fatih/color"
	"log"
	"strings"
)

var (
	runningStatus     = color.New(color.FgGreen)
	establishingStatus = color.New(color.FgYellow)
	stoppedStatus     = color.New(color.FgHiRed)
)

func main() {
	cluster := flag.String("cluster", "qa", "the cluster name")
	suffix := flag.String("suffix", "", "the service name suffix to look for")
	services := flag.String("services", "", "a comma separated list of services (max of 10)")

	flag.Parse()

	serviceNames := serviceNames(*suffix, *services)
	config := newConfig(cluster, serviceNames)
	result := listServicesData(config)

	for _, service := range result.Services {
		name := service.ServiceName
		status := service.Status
		desired := service.DesiredCount
		running := service.RunningCount

		prettyPrint(*name, *status, *desired, *running)
	}
}

// To be improved
func prettyPrint(serviceName, status string, desired, running int64) {
	if status == "ACTIVE" && desired == running {
		runningStatus.Printf("%s: status %s, desired: %d, running: %d\n", serviceName, status, desired, running)
	} else if status == "ACTIVE"  && desired != running || (status == "DRAINING") {
		establishingStatus.Printf("%s: status %s, desired: %d, running: %d\n", serviceName, status, desired, running)
	} else if status == "INACTIVE" && desired == running {
		stoppedStatus.Printf("%s: status %s, desired: %d, running: %d\n", serviceName, status, desired, running)
	}
}

func serviceNames(suffix, services string) []*string {
	names := strings.Split(services,",")
	var result []*string
	for _, name := range names {
		fullName := strings.Trim(fmt.Sprintf("%s%s", name, suffix), "")
		result = append(result,  aws.String(fullName))
	}

	return result
}

func newSession() *session.Session {
	session, err := session.NewSession(
		&aws.Config{Region: aws.String("us-east-1")})

	if err != nil {
		log.Fatal(err)
	}

	return session
}

func listServicesData(config *ecs.DescribeServicesInput) *ecs.DescribeServicesOutput {
	svc := ecs.New(newSession())
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
				fmt.Println(ecs.ErrCodeClusterNotFoundException, aerr.Error())
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

func newConfig(cluster *string, services []*string) *ecs.DescribeServicesInput{
	return &ecs.DescribeServicesInput{
		Cluster:  cluster,
		Services: services,
	}
}