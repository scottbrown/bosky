package main

import (
	"errors"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
	"github.com/aws/aws-sdk-go/aws/session"
)

func retrieveInstanceId(sess *session.Session) (string, error) {
	metadata_svc := ec2metadata.New(sess)

	if !metadata_svc.Available() {
		// are we not on an EC2 machine?
		// ...or maybe the service is down?
		return "", errors.New("Cannot lookup instance ID because EC2 metadata service is not available.  Use --instance-id")
	}

	identity_document, err := metadata_svc.GetInstanceIdentityDocument()

	if err != nil {
		return "", err
	}

	return identity_document.InstanceID, nil
}
