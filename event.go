package bosky

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

type Emitter struct {
	UserDataStatus string
	StatusFail     bool
	StatusInfo     bool
	StatusPass     bool
	InstanceID     string
	Project        string
}

func (e Emitter) chooseStatusMessage() string {
	if e.UserDataStatus != "" {
		return e.UserDataStatus
	}

	if e.StatusFail {
		return STATUS_FAIL
	}

	if e.StatusPass {
		return STATUS_PASS
	}

	if e.StatusInfo {
		return STATUS_INFO
	}

	return STATUS_INFO
}

func (e Emitter) EmitEvent(message string) error {
	ctx := context.TODO()

	// Load the AWS SDK configuration
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRetryMaxAttempts(3))
	if err != nil {
		return err
	}

	if e.InstanceID == "" {
		e.InstanceID, err = retrieveInstanceId(cfg)
		if err != nil {
			return err
		}
	}

	// Create CloudWatch Events client
	client := cloudwatchevents.NewFromConfig(cfg)

	status := e.chooseStatusMessage()
	detail := fmt.Sprintf("{ \"Status\": \"%s\", \"Message\": \"%s\"}", status, message)
	detailType := "User Data"

	// Create the PutEvents input
	input := &cloudwatchevents.PutEventsInput{
		Entries: []types.PutEventsRequestEntry{
			{
				Detail:     aws.String(detail),
				DetailType: aws.String(detailType),
				Resources:  []string{e.InstanceID},
				Source:     aws.String(e.Project),
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
