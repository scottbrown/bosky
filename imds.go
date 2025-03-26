package beacon

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
)

type IMDSClient interface {
	GetMetadata(context.Context, *imds.GetMetadataInput, ...func(*imds.Options)) (*imds.GetMetadataOutput, error)
	GetInstanceIdentityDocument(context.Context, *imds.GetInstanceIdentityDocumentInput, ...func(*imds.Options)) (*imds.GetInstanceIdentityDocumentOutput, error)
}

func retrieveInstanceId(client IMDSClient) (string, error) {
	ctx := context.TODO()

	// Check if IMDS is available
	_, err := client.GetMetadata(ctx, &imds.GetMetadataInput{
		Path: "instance-id",
	})

	if err != nil {
		return "", errors.New("Cannot lookup instance ID because EC2 metadata service is not available. Use --instance-id")
	}

	// Get instance identity document
	resp, err := client.GetInstanceIdentityDocument(ctx, &imds.GetInstanceIdentityDocumentInput{})
	if err != nil {
		return "", err
	}

	return resp.InstanceIdentityDocument.InstanceID, nil
}
