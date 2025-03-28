package beacon

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchevents"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchevents/types"
)

type Emitter struct {
	InstanceID InstanceID
	Project    Project

	EBClient   EventBridgeClient
	IMDSClient IMDSClient
}

func (e Emitter) Emit(ctx context.Context, status Status, message string) error {
	d := Detail{
		Status:  string(status),
		Message: message,
	}
	detail, err := json.Marshal(d)
	if err != nil {
		return err
	}

	detailType := DetailType(DEFAULT_DETAIL_TYPE)

	if err := detailType.Validate(); err != nil {
		return err
	}

	if err := e.InstanceID.Validate(); err != nil {
		return err
	}

	input := &cloudwatchevents.PutEventsInput{
		Entries: []types.PutEventsRequestEntry{
			{
				Detail:     aws.String(string(detail)),
				DetailType: aws.String(string(detailType)),
				Resources:  []string{string(e.InstanceID)},
				Source:     aws.String(string(e.Project)),
			},
		},
	}

	if err := validateEventPayloadSize(input); err != nil {
		return err
	}

	_, err = e.EBClient.PutEvents(ctx, input)
	if err != nil {
		return err
	}

	return nil
}

func validateEventPayloadSize(input *cloudwatchevents.PutEventsInput) error {
	totalBytes := 0
	for _, j := range input.Entries {
		totalBytes += len(*j.Detail)
		totalBytes += len(*j.DetailType)
		for _, i := range j.Resources {
			totalBytes += len(i)
		}
		totalBytes += len(*j.Source)

		if totalBytes > EVENT_PAYLOAD_MAX_BYTES {
			return fmt.Errorf("Payload size of %d bytes exceeds %d bytes", totalBytes, EVENT_PAYLOAD_MAX_BYTES)
		}
	}

	return nil
}
