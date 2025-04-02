package beacon

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
)

func TestRetrieveInstanceARN(t *testing.T) {
	tests := []struct {
		name             string
		setupMock        func(*MockIMDSClient)
		expectedARN      InstanceARN
		expectError      bool
		expectedErrorMsg string
	}{
		{
			name: "Successfully retrieve instance ARN",
			setupMock: func(mockClient *MockIMDSClient) {
				mockClient.GetMetadataFunc = func(ctx context.Context, input *imds.GetMetadataInput, opts ...func(*imds.Options)) (*imds.GetMetadataOutput, error) {
					return &imds.GetMetadataOutput{}, nil
				}

				mockClient.GetInstanceIdentityDocumentFunc = func(ctx context.Context, input *imds.GetInstanceIdentityDocumentInput, opts ...func(*imds.Options)) (*imds.GetInstanceIdentityDocumentOutput, error) {
					return &imds.GetInstanceIdentityDocumentOutput{
						InstanceIdentityDocument: imds.InstanceIdentityDocument{
							AccountID:  "0123456789012",
							Region:     "us-east-1",
							InstanceID: "i-12345abcdef",
						},
					}, nil
				}
			},
			expectedARN: "arn:aws:ec2:us-east-1:0123456789012:instance/i-12345abcdef",
			expectError: false,
		},
		{
			name: "IMDS not available",
			setupMock: func(mockClient *MockIMDSClient) {
				mockClient.GetMetadataFunc = func(ctx context.Context, input *imds.GetMetadataInput, opts ...func(*imds.Options)) (*imds.GetMetadataOutput, error) {
					return nil, errors.New("IMDS not available")
				}
			},
			expectedARN:      "",
			expectError:      true,
			expectedErrorMsg: "ec2 metadata service is not available.",
		},
		{
			name: "Error getting instance identity document",
			setupMock: func(mockClient *MockIMDSClient) {
				mockClient.GetMetadataFunc = func(ctx context.Context, input *imds.GetMetadataInput, opts ...func(*imds.Options)) (*imds.GetMetadataOutput, error) {
					return &imds.GetMetadataOutput{}, nil
				}

				mockClient.GetInstanceIdentityDocumentFunc = func(ctx context.Context, input *imds.GetInstanceIdentityDocumentInput, opts ...func(*imds.Options)) (*imds.GetInstanceIdentityDocumentOutput, error) {
					return nil, errors.New("Failed to get instance identity document")
				}
			},
			expectedARN:      "",
			expectError:      true,
			expectedErrorMsg: "Failed to get instance identity document",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockIMDSClient{}
			tt.setupMock(mockClient)

			instanceARN, err := RetrieveInstanceARN(context.Background(), mockClient)

			if tt.expectError {
				if err == nil {
					t.Error("Expected an error but got nil")
				} else if tt.expectedErrorMsg != "" && err.Error() != tt.expectedErrorMsg {
					t.Errorf("Expected error message %q, got %q", tt.expectedErrorMsg, err.Error())
				}
			} else if err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			if instanceARN != tt.expectedARN {
				t.Errorf("Expected instance ARN %q, got %q", tt.expectedARN, instanceARN)
			}
		})
	}
}
