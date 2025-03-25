package bosky

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchevents"
)

// Mock implementations
type MockEventBridgeClient struct {
	PutEventsFunc func(context.Context, *cloudwatchevents.PutEventsInput, ...func(*cloudwatchevents.Options)) (*cloudwatchevents.PutEventsOutput, error)
}

func (m MockEventBridgeClient) PutEvents(ctx context.Context, input *cloudwatchevents.PutEventsInput, opts ...func(*cloudwatchevents.Options)) (*cloudwatchevents.PutEventsOutput, error) {
	return m.PutEventsFunc(ctx, input, opts...)
}

type MockIMDSClient struct {
	GetMetadataFunc                 func(context.Context, *imds.GetMetadataInput, ...func(*imds.Options)) (*imds.GetMetadataOutput, error)
	GetInstanceIdentityDocumentFunc func(context.Context, *imds.GetInstanceIdentityDocumentInput, ...func(*imds.Options)) (*imds.GetInstanceIdentityDocumentOutput, error)
}

func (m MockIMDSClient) GetMetadata(ctx context.Context, input *imds.GetMetadataInput, opts ...func(*imds.Options)) (*imds.GetMetadataOutput, error) {
	return m.GetMetadataFunc(ctx, input, opts...)
}

func (m MockIMDSClient) GetInstanceIdentityDocument(ctx context.Context, input *imds.GetInstanceIdentityDocumentInput, opts ...func(*imds.Options)) (*imds.GetInstanceIdentityDocumentOutput, error) {
	return m.GetInstanceIdentityDocumentFunc(ctx, input, opts...)
}

func TestEmitterChooseStatusMessage(t *testing.T) {
	tests := []struct {
		name           string
		emitter        Emitter
		expectedStatus string
	}{
		{
			name: "UserDataStatus is set",
			emitter: Emitter{
				UserDataStatus: "custom",
				StatusFail:     false,
				StatusInfo:     false,
				StatusPass:     false,
			},
			expectedStatus: "custom",
		},
		{
			name: "StatusFail is set",
			emitter: Emitter{
				UserDataStatus: "",
				StatusFail:     true,
				StatusInfo:     false,
				StatusPass:     false,
			},
			expectedStatus: STATUS_FAIL,
		},
		{
			name: "StatusPass is set",
			emitter: Emitter{
				UserDataStatus: "",
				StatusFail:     false,
				StatusInfo:     false,
				StatusPass:     true,
			},
			expectedStatus: STATUS_PASS,
		},
		{
			name: "StatusInfo is set",
			emitter: Emitter{
				UserDataStatus: "",
				StatusFail:     false,
				StatusInfo:     true,
				StatusPass:     false,
			},
			expectedStatus: STATUS_INFO,
		},
		{
			name: "No status is set",
			emitter: Emitter{
				UserDataStatus: "",
				StatusFail:     false,
				StatusInfo:     false,
				StatusPass:     false,
			},
			expectedStatus: STATUS_INFO, // Default
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status := tt.emitter.chooseStatusMessage()
			if status != tt.expectedStatus {
				t.Errorf("Expected status %s, got %s", tt.expectedStatus, status)
			}
		})
	}
}

func TestEmitterEmit(t *testing.T) {
	tests := []struct {
		name           string
		emitter        Emitter
		message        string
		setupMocks     func(*MockEventBridgeClient, *MockIMDSClient)
		expectedError  bool
		expectedErrMsg string
	}{
		{
			name: "Successfully emit event with instance ID provided",
			emitter: Emitter{
				StatusInfo: true,
				InstanceID: "i-12345",
				Project:    "test-project",
			},
			message: "Test message",
			setupMocks: func(ebClient *MockEventBridgeClient, imdsClient *MockIMDSClient) {
				ebClient.PutEventsFunc = func(ctx context.Context, input *cloudwatchevents.PutEventsInput, opts ...func(*cloudwatchevents.Options)) (*cloudwatchevents.PutEventsOutput, error) {
					return &cloudwatchevents.PutEventsOutput{}, nil
				}
			},
			expectedError: false,
		},
		{
			name: "Successfully emit event with instance ID fetched",
			emitter: Emitter{
				StatusInfo: true,
				InstanceID: "", // Empty, should use IMDS
				Project:    "test-project",
			},
			message: "Test message",
			setupMocks: func(ebClient *MockEventBridgeClient, imdsClient *MockIMDSClient) {
				ebClient.PutEventsFunc = func(ctx context.Context, input *cloudwatchevents.PutEventsInput, opts ...func(*cloudwatchevents.Options)) (*cloudwatchevents.PutEventsOutput, error) {
					return &cloudwatchevents.PutEventsOutput{}, nil
				}

				imdsClient.GetMetadataFunc = func(ctx context.Context, input *imds.GetMetadataInput, opts ...func(*imds.Options)) (*imds.GetMetadataOutput, error) {
					return &imds.GetMetadataOutput{}, nil
				}

				imdsClient.GetInstanceIdentityDocumentFunc = func(ctx context.Context, input *imds.GetInstanceIdentityDocumentInput, opts ...func(*imds.Options)) (*imds.GetInstanceIdentityDocumentOutput, error) {
					return &imds.GetInstanceIdentityDocumentOutput{
						InstanceIdentityDocument: imds.InstanceIdentityDocument{
							InstanceID: "i-fetched",
						},
					}, nil
				}
			},
			expectedError: false,
		},
		{
			name: "Fail to retrieve instance ID",
			emitter: Emitter{
				StatusInfo: true,
				InstanceID: "", // Empty, should use IMDS
				Project:    "test-project",
			},
			message: "Test message",
			setupMocks: func(ebClient *MockEventBridgeClient, imdsClient *MockIMDSClient) {
				imdsClient.GetMetadataFunc = func(ctx context.Context, input *imds.GetMetadataInput, opts ...func(*imds.Options)) (*imds.GetMetadataOutput, error) {
					return nil, errors.New("IMDS not available")
				}
			},
			expectedError:  true,
			expectedErrMsg: "Cannot lookup instance ID because EC2 metadata service is not available. Use --instance-id",
		},
		{
			name: "Fail to send event",
			emitter: Emitter{
				StatusInfo: true,
				InstanceID: "i-12345",
				Project:    "test-project",
			},
			message: "Test message",
			setupMocks: func(ebClient *MockEventBridgeClient, imdsClient *MockIMDSClient) {
				ebClient.PutEventsFunc = func(ctx context.Context, input *cloudwatchevents.PutEventsInput, opts ...func(*cloudwatchevents.Options)) (*cloudwatchevents.PutEventsOutput, error) {
					return nil, errors.New("Failed to send event")
				}
			},
			expectedError:  true,
			expectedErrMsg: "Failed to send event",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockEBClient := &MockEventBridgeClient{}
			mockIMDSClient := &MockIMDSClient{}

			tt.setupMocks(mockEBClient, mockIMDSClient)

			tt.emitter.EBClient = mockEBClient
			tt.emitter.IMDSClient = mockIMDSClient

			err := tt.emitter.Emit(tt.message)

			if tt.expectedError {
				if err == nil {
					t.Errorf("Expected error but got nil")
				} else if tt.expectedErrMsg != "" && err.Error() != tt.expectedErrMsg {
					t.Errorf("Expected error message %q, got %q", tt.expectedErrMsg, err.Error())
				}
			} else if err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}
