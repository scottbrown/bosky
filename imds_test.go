package beacon

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
)

func TestRetrieveInstanceID(t *testing.T) {
	tests := []struct {
		name             string
		setupMock        func(*MockIMDSClient)
		expectedID       string
		expectError      bool
		expectedErrorMsg string
	}{
		{
			name: "Successfully retrieve instance ID",
			setupMock: func(mockClient *MockIMDSClient) {
				mockClient.GetMetadataFunc = func(ctx context.Context, input *imds.GetMetadataInput, opts ...func(*imds.Options)) (*imds.GetMetadataOutput, error) {
					return &imds.GetMetadataOutput{}, nil
				}

				mockClient.GetInstanceIdentityDocumentFunc = func(ctx context.Context, input *imds.GetInstanceIdentityDocumentInput, opts ...func(*imds.Options)) (*imds.GetInstanceIdentityDocumentOutput, error) {
					return &imds.GetInstanceIdentityDocumentOutput{
						InstanceIdentityDocument: imds.InstanceIdentityDocument{
							InstanceID: "i-12345abcdef",
						},
					}, nil
				}
			},
			expectedID:  "i-12345abcdef",
			expectError: false,
		},
		{
			name: "IMDS not available",
			setupMock: func(mockClient *MockIMDSClient) {
				mockClient.GetMetadataFunc = func(ctx context.Context, input *imds.GetMetadataInput, opts ...func(*imds.Options)) (*imds.GetMetadataOutput, error) {
					return nil, errors.New("IMDS not available")
				}
			},
			expectedID:       "",
			expectError:      true,
			expectedErrorMsg: "Cannot lookup instance ID because EC2 metadata service is not available. Use --instance-id",
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
			expectedID:       "",
			expectError:      true,
			expectedErrorMsg: "Failed to get instance identity document",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockIMDSClient{}
			tt.setupMock(mockClient)

			instanceID, err := retrieveInstanceId(mockClient)

			if tt.expectError {
				if err == nil {
					t.Error("Expected an error but got nil")
				} else if tt.expectedErrorMsg != "" && err.Error() != tt.expectedErrorMsg {
					t.Errorf("Expected error message %q, got %q", tt.expectedErrorMsg, err.Error())
				}
			} else if err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}

			if instanceID != tt.expectedID {
				t.Errorf("Expected instance ID %q, got %q", tt.expectedID, instanceID)
			}
		})
	}
}
