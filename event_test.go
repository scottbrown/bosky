package beacon

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

func TestEmitterEmit(t *testing.T) {
	tests := []struct {
		name           string
		emitter        Emitter
		status         Status
		message        string
		setupMocks     func(*MockEventBridgeClient, *MockIMDSClient)
		expectedError  bool
		expectedErrMsg string
	}{
		{
			name: "Successfully emit event with instance ID provided",
			emitter: Emitter{
				InstanceID: "arn:aws:ec2:region:0123456789012:instance/i-12345abcdef",
				Project:    "test-project",
			},
			status:  Status(STATUS_INFO),
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
				InstanceID: "", // Empty, should use IMDS
				Project:    "test-project",
			},
			status:  Status(STATUS_INFO),
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
              InstanceID: "arn:aws:ec2:region:0123456789012:instance/i-12345abcdef",
						},
					}, nil
				}
			},
			expectedError: false,
		},
		{
			name: "Fail to send event",
			emitter: Emitter{
        InstanceID: "arn:aws:ec2:region:0123456789012:instance/i-12345abcdef",
				Project:    "test-project",
			},
			status:  Status(STATUS_INFO),
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

			err := tt.emitter.Emit(context.Background(), tt.status, tt.message)

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
