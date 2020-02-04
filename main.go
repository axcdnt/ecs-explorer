package main

import (
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/axcdnt/ecs-explorer/services"
	"log"
	"strings"
)

func main() {
	clusterFlag := flag.String("cluster", "", "the cluster name")
	suffixFlag := flag.String("suffix", "", "the service name suffix to look for")
	servicesFlag := flag.String("services", "", "a comma separated list of services (max of 10)")

	flag.Parse()

	validateFlags(*clusterFlag, *suffixFlag, *servicesFlag)

	session := awsSession()
	serviceNames := parseServiceNames(*suffixFlag, *servicesFlag)
	// aws limits the query in 10 services/request
	queries := partitionIn(serviceNames, 10)

	// to be improved: run the queries in parallel
	for _, query := range queries {
		ecsService := services.EcsService{
			Cluster:  clusterFlag,
			Services: query,
			Session:  session,
		}

		ecsService.Query()
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
func partitionIn(services []*string, size int) [][]*string {
	numOfSlices := len(services) / size
	var allSlices [][]*string
	start := 0
	end := size

	for i := 0; i < numOfSlices; i++ {
		slice := services[start:end]
		allSlices = append(allSlices, slice)
		start = end
		end += 10
	}

	lastSlice := services[numOfSlices*size:]
	if len(lastSlice) > 0 {
		// corner case to avoid empty slices
		allSlices = append(allSlices, lastSlice)
	}

	return allSlices
}