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

	validateFlags(*clusterFlag, *suffixFlag, *servicesFlag)

	session := awsSession()
	serviceNames := parseServiceNames(*suffixFlag, *servicesFlag)

	// aws ecs max limit for cluster queries
	sliceIn(serviceNames, 10)

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

func validateFlags(clusterFlag, suffixFlag, servicesFlag string) {
	if len(clusterFlag) == 0 {
		log.Fatalln("missing flag: -cluster")
	}

	if len(suffixFlag) == 0 {
		log.Fatalln("missing flag: -suffix")
	}

	if len(servicesFlag) == 0 {
		log.Fatalln("missing flag: -services")
	}
}

func sliceIn(services []*string, sliceSize int) [][]*string {
	numOfSlices := len(services) / sliceSize
	var allSlices [][]*string
	start := 0
	end := sliceSize

	for i := 0; i < numOfSlices; i++ {
		slice := services[start:end]
		allSlices = append(allSlices, slice)
		start = end
		end += 10
	}

	lastSlice := services[numOfSlices*sliceSize:]
	allSlices = append(allSlices, lastSlice)

	return allSlices
}