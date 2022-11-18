package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage:")
		fmt.Println(" startstopinstances (start|stop|status)")
		return
	}

	action := os.Args[1]
	if action != "start" && action != "stop" && action != "status" {
		fmt.Println("Usage:")
		fmt.Println(" startstopinstances (start|stop|status)")
		return
	}

	homedir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Unable to retriever user home directory.")
		return
	}

	f, err := os.Open(homedir + "/.aws/instances")
	if err != nil {
		fmt.Println("Error openning ~/.aws/instances.")
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	instances := []*string{}
	for scanner.Scan() {
		instances = append(instances, aws.String(scanner.Text()))
	}

	session := session.Must(session.NewSession(&aws.Config{Region: aws.String("us-east-1")}))
	svc := ec2.New(session)

	if action == "stop" {
		fmt.Println("Stopping instances: ")
		for _, i := range instances {
			fmt.Println(*i)
		}

		input := &ec2.StopInstancesInput{
			InstanceIds: instances,
		}

		_, err = svc.StopInstances(input)
		if err != nil {
			println(err.Error())
		}
	}

	if action == "start" {
		fmt.Println("Starting instances: ")
		for _, i := range instances {
			fmt.Println(*i)
		}
		input := &ec2.StartInstancesInput{
			InstanceIds: instances,
		}

		_, err = svc.StartInstances(input)
		if err != nil {
			println(err.Error())
		}
	}

	if action == "status" {
		result, err := svc.DescribeInstances(nil)
		if err != nil {
			panic(err.Error())
		}
		for _, reservation := range result.Reservations {
			for _, instance := range reservation.Instances {
				id := *instance.InstanceId
				name := ""
				for _, i := range instance.Tags {
					if *i.Key == "Name" {
						name = *i.Value
					}
				}
				state := *instance.State.Name
				ip := "none"
				if instance.PublicIpAddress != nil {
					ip = *instance.PublicIpAddress
				}

				fmt.Printf("%-20s %-20s %10s %15s\n", id, name, state, ip)
			}
		}
	}
}
