package main

import (
	"context"
	"fmt"
)

// handleListPermissions displays the IAM permissions required by the
// application.
func handleListPermissions(ctx context.Context) error {
	fmt.Println("Required IAM permissions:")

	for _, p := range requiredPermissions {
		fmt.Printf("- %s\n", p)
	}

	fmt.Println("\nJSON Policy Document:")
	fmt.Println(`{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": "events:PutEvents",
      "Resource": "*"
    }
  ]
}`)

	return nil
}
