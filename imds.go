package beacon

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
)

// RetrieveInstanceID gets the current EC2 instance ID from the metadata
// service.  If the metadata service is unavailable, it returns an error
// suggesting the use of `--instance-id`.
func RetrieveInstanceID(ctx context.Context, client IMDSClient) (string, error) {
	// Check if IMDS is available
	_, err := client.GetMetadata(ctx, &imds.GetMetadataInput{
		Path: "instance-id",
	})

	if err != nil {
		return "", errors.New("cannot lookup instance ID because EC2 metadata service is not available. Use --instance-id")
	}

	// Get instance identity document
	resp, err := client.GetInstanceIdentityDocument(ctx, &imds.GetInstanceIdentityDocumentInput{})
	if err != nil {
		return "", err
	}

	return resp.InstanceIdentityDocument.InstanceID, nil //nolint:all
}
