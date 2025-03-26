package beacon

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/cloudwatchevents"
)

func TestSendEvent(t *testing.T) {
	tests := []struct {
		name           string
		detail         string
		detailType     string
		instanceID     string
		project        string
		clientBehavior func(ctx context.Context, input *cloudwatchevents.PutEventsInput, opts ...func(*cloudwatchevents.Options)) (*cloudwatchevents.PutEventsOutput, error)
		expectError    bool
	}{
		{
			name:       "Successfully send event",
			detail:     "{ \"Status\": \"pass\", \"Message\": \"Test message\" }",
			detailType: "User Data",
			instanceID: "i-12345",
			project:    "test-project",
			clientBehavior: func(ctx context.Context, input *cloudwatchevents.PutEventsInput, opts ...func(*cloudwatchevents.Options)) (*cloudwatchevents.PutEventsOutput, error) {
				// Validate input
				if len(input.Entries) != 1 {
					t.Errorf("Expected 1 entry, got %d", len(input.Entries))
				}
				entry := input.Entries[0]
				if *entry.Detail != "{ \"Status\": \"pass\", \"Message\": \"Test message\" }" {
					t.Errorf("Expected detail %q, got %q", "{ \"Status\": \"pass\", \"Message\": \"Test message\" }", *entry.Detail)
				}
				if *entry.DetailType != "User Data" {
					t.Errorf("Expected detailType %q, got %q", "User Data", *entry.DetailType)
				}
				if len(entry.Resources) != 1 || entry.Resources[0] != "i-12345" {
					t.Errorf("Expected resource i-12345, got %v", entry.Resources)
				}
				if *entry.Source != "test-project" {
					t.Errorf("Expected source %q, got %q", "test-project", *entry.Source)
				}

				return &cloudwatchevents.PutEventsOutput{}, nil
			},
			expectError: false,
		},
		{
			name:       "Fail to send event",
			detail:     "{ \"Status\": \"fail\", \"Message\": \"Error message\" }",
			detailType: "User Data",
			instanceID: "i-67890",
			project:    "test-project",
			clientBehavior: func(ctx context.Context, input *cloudwatchevents.PutEventsInput, opts ...func(*cloudwatchevents.Options)) (*cloudwatchevents.PutEventsOutput, error) {
				return nil, errors.New("Failed to send event")
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &MockEventBridgeClient{
				PutEventsFunc: tt.clientBehavior,
			}

			err := sendEvent(mockClient, tt.detail, tt.detailType, tt.instanceID, tt.project)

			if tt.expectError && err == nil {
				t.Error("Expected error but got nil")
			} else if !tt.expectError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}
