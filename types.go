package beacon

import (
	"context"
	"fmt"
	"regexp"

	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchevents"
)

// EventBridgeClient defines the interface for EventBridge operations
type EventBridgeClient interface {
	PutEvents(context.Context, *cloudwatchevents.PutEventsInput, ...func(*cloudwatchevents.Options)) (*cloudwatchevents.PutEventsOutput, error)
}

// IMDSClient defines the interface for AWS EC2 Instance Metadata Service
// operations
type IMDSClient interface {
	GetMetadata(context.Context, *imds.GetMetadataInput, ...func(*imds.Options)) (*imds.GetMetadataOutput, error)
	GetInstanceIdentityDocument(context.Context, *imds.GetInstanceIdentityDocumentInput, ...func(*imds.Options)) (*imds.GetInstanceIdentityDocumentOutput, error)
}

// Status represents the status of an event (pass, fail, info, or custom)
type Status string

// Project represents the project name for event source identification
type Project string

// InstanceID represents an EC2 instance identifier
type InstanceID string

// DetailType represents the EventBridge detail type
type DetailType string

// Detail contains the main payload for an EventBridge event
type Detail struct {
	Status  string `json:"Status"`  // Status of the event
	Message string `json:"Message"` // Message associated with the event
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
		return fmt.Errorf("instance ID cannot be empty")
	}

	if size > RESOURCE_ARN_MAX_LENGTH {
		return fmt.Errorf("instance ID length of %d bytes exceeds %d bytes", size, RESOURCE_ARN_MAX_LENGTH)
	}

	arnPattern := regexp.MustCompile(`^arn:aws:ec2:[a-z0-9-]+:[0-9]+:instance/i-[a-z0-9]+$`)

	if !arnPattern.MatchString(string(i)) {
		return fmt.Errorf("invalid format. Must be a valid EC2 instance ARN")
	}

	return nil
}
