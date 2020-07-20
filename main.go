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

// aws limits the query in 10 services/request
const maxServicesPerRequest = 10

func main() {
	clusterFlag := flag.String("cluster", "", "the cluster name")
	suffixFlag := flag.String("suffix", "", "the service name suffix to look for")
	servicesFlag := flag.String("services", "", "a comma separated list of services")

	flag.Parse()

	validateFlags(*clusterFlag, *suffixFlag, *servicesFlag)

	session := awsSession()
	serviceNames := parseServiceNames(*suffixFlag, *servicesFlag)
	queries := partitionIn(serviceNames)

	resp := make(chan string)
	for _, query := range queries {
		ecsService := services.EcsService{
			Cluster:  clusterFlag,
			Services: query,
			Session:  session,
		}

		go ecsService.Query(resp)
	}
	for range serviceNames {
		fmt.Print(<-resp)
	}
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

// this function makes partitions of a slice based on the desired size
func partitionIn(services []*string) [][]*string {
	numOfSlices := len(services) / maxServicesPerRequest
	startRange := 0
	endRange := maxServicesPerRequest

	var allSlices [][]*string
	for i := 0; i < numOfSlices; i++ {
		slice := services[startRange:endRange]
		allSlices = append(allSlices, slice)
		startRange = endRange
		endRange += maxServicesPerRequest
	}

	lastSlice := services[numOfSlices*maxServicesPerRequest:]
	if len(lastSlice) > 0 {
		// corner case to avoid empty slices
		allSlices = append(allSlices, lastSlice)
	}

	return allSlices
}
