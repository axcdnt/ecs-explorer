package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/axcdnt/ecs-explorer/services"
)

func main() {
	clusterFlag := flag.String("cluster", "", "the cluster name")
	suffixFlag := flag.String("suffix", "", "the service name suffix to look for")
	servicesFlag := flag.String("services", "", "a comma separated list of services (max of 10)")

	flag.Parse()

	session := awsSession()
	serviceNames := parseServiceNames(*suffixFlag, *servicesFlag)

	ecsService := services.EcsService{
		Cluster:  clusterFlag,
		Services: serviceNames,
		Session:  session,
	}

	ecsService.List()
}

func parseServiceNames(suffix, services string) []*string {
	names := strings.Split(services, ",")
	var result []*string
	for _, name := range names {
		fullName := strings.Trim(fmt.Sprintf("%s%s", name, suffix), "")
		result = append(result, aws.String(fullName))
	}

	return result
}

func awsSession() *session.Session {
	session, err := session.NewSession(
		&aws.Config{Region: aws.String("us-east-1")})
	if err != nil {
		log.Fatal(err)
	}

	return session
}

