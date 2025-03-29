package beacon

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchevents"
)

type EventBridgeClient interface {
	PutEvents(context.Context, *cloudwatchevents.PutEventsInput, ...func(*cloudwatchevents.Options)) (*cloudwatchevents.PutEventsOutput, error)
}

type IMDSClient interface {
	GetMetadata(context.Context, *imds.GetMetadataInput, ...func(*imds.Options)) (*imds.GetMetadataOutput, error)
	GetInstanceIdentityDocument(context.Context, *imds.GetInstanceIdentityDocumentInput, ...func(*imds.Options)) (*imds.GetInstanceIdentityDocumentOutput, error)
}

type Status string

type Project string

type InstanceID string

type DetailType string

type Detail struct {
	Status  string `json:"Status"`
	Message string `json:"Message"`
}

func (d DetailType) Validate() error {
	size := len(d)
	if size == 0 {
		return fmt.Errorf("Detail type cannot be empty")
	}

	if size > DETAIL_TYPE_MAX_LENGTH {
		return fmt.Errorf("Detail type length of %d bytes exceeds %d bytes", size, DETAIL_TYPE_MAX_LENGTH)
	}

	return nil
}

func (i InstanceID) Validate() error {
	size := len(i)
	if size == 0 {
		return fmt.Errorf("Instance ID cannot be empty")
	}

	if size > RESOURCE_ARN_MAX_LENGTH {
		return fmt.Errorf("Instance ID length of %d bytes exceeds %d bytes", size, RESOURCE_ARN_MAX_LENGTH)
	}

	return nil
}
