package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchevents"
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
	sess := session.Must(session.NewSession(aws.NewConfig().
		WithMaxRetries(3),
	))

	var err error
	if instance_id == "" {
		instance_id, err = retrieveInstanceId(sess)
		if err != nil {
			return err
		}
	}

	event_svc := cloudwatchevents.New(sess)

	status := chooseStatusMessage()
	detail := fmt.Sprintf("{ \"Status\": \"%s\", \"Message\": \"%s\"}", status, message)
	detailType := "User Data"

	request_entry := cloudwatchevents.PutEventsRequestEntry{
		Detail:     &detail,
		DetailType: &detailType,
		Resources:  []*string{&instance_id},
		Source:     &project,
	}

	input := cloudwatchevents.PutEventsInput{
		Entries: []*cloudwatchevents.PutEventsRequestEntry{
			&request_entry,
		},
	}

	_, err = event_svc.PutEvents(&input)

	if err != nil {
		return err
	}

	return nil
}
