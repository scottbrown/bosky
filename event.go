package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchevents"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchevents/types"
)

const (
	STATUS_FAIL string = "fail"
	STATUS_INFO string = "info"
	STATUS_PASS string = "pass"
)

func chooseStatusMessage() string {
	if userDataStatus != "" {
		return userDataStatus
	}

	if statusFail {
		return STATUS_FAIL
	}

	if statusPass {
		return STATUS_PASS
	}

	if statusInfo {
		return STATUS_INFO
	}

	return STATUS_INFO
}

func emitEvent(message string) error {
	ctx := context.TODO()

	// Load the AWS SDK configuration
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRetryMaxAttempts(3))
	if err != nil {
		return err
	}

	if instance_id == "" {
		instance_id, err = retrieveInstanceId(cfg)
		if err != nil {
			return err
		}
	}

	// Create CloudWatch Events client
	client := cloudwatchevents.NewFromConfig(cfg)

	status := chooseStatusMessage()
	detail := fmt.Sprintf("{ \"Status\": \"%s\", \"Message\": \"%s\"}", status, message)
	detailType := "User Data"

	// Create the PutEvents input
	input := &cloudwatchevents.PutEventsInput{
		Entries: []types.PutEventsRequestEntry{
			{
				Detail:     aws.String(detail),
				DetailType: aws.String(detailType),
				Resources:  []string{instance_id},
				Source:     aws.String(project),
			},
		},
	}

	// Send the event
	_, err = client.PutEvents(ctx, input)
	if err != nil {
		return err
	}

	return nil
}
