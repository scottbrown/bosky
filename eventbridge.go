package bosky

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchevents"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchevents/types"
)

type EventBridgeClient interface {
	PutEvents(context.Context, *cloudwatchevents.PutEventsInput, ...func(*cloudwatchevents.Options)) (*cloudwatchevents.PutEventsOutput, error)
}

func sendEvent(client EventBridgeClient, detail, detailType, instanceID, project string) error {
	ctx := context.TODO()

	// Create the PutEvents input
	input := &cloudwatchevents.PutEventsInput{
		Entries: []types.PutEventsRequestEntry{
			{
				Detail:     aws.String(detail),
				DetailType: aws.String(detailType),
				Resources:  []string{instanceID},
				Source:     aws.String(project),
			},
		},
	}

	// Send the event
	_, err := client.PutEvents(ctx, input)
	if err != nil {
		return err
	}

	return nil
}
