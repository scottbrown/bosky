package beacon

import (
	"context"
	"errors"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
)

// RetrieveInstanceARN gets the current EC2 instance ARN from the metadata
// service.  If the metadata service is unavailable, it returns an error.
func RetrieveInstanceARN(ctx context.Context, client IMDSClient) (InstanceARN, error) {
	// Check if IMDS is available
	_, err := client.GetMetadata(ctx, &imds.GetMetadataInput{
		Path: "instance-id",
	})

	if err != nil {
		return "", errors.New("ec2 metadata service is not available.")
	}

	// Get instance identity document
	resp, err := client.GetInstanceIdentityDocument(ctx, &imds.GetInstanceIdentityDocumentInput{})
	if err != nil {
		return "", err
	}

	instanceID := resp.InstanceIdentityDocument.InstanceID
	region := resp.InstanceIdentityDocument.Region
	accountID := resp.InstanceIdentityDocument.AccountID

	arn := fmt.Sprintf("arn:aws:ec2:%s:%s:instance/%s", region, accountID, instanceID)

	return InstanceARN(arn), nil
}
