package main

import (
	"testing"

	"github.com/scottbrown/beacon"
)

func TestChooseStatusMessage(t *testing.T) {
	tests := []struct {
		name           string
		userDataStatus string
		statusFail     bool
		statusInfo     bool
		statusPass     bool
		expectedStatus beacon.Status
	}{
		{
			name:           "UserDataStatus is set",
			userDataStatus: "custom",
			statusFail:     false,
			statusInfo:     false,
			statusPass:     false,
			expectedStatus: beacon.Status("custom"),
		},
		{
			name:           "StatusFail is set",
			userDataStatus: "",
			statusFail:     true,
			statusInfo:     false,
			statusPass:     false,
			expectedStatus: beacon.Status(beacon.STATUS_FAIL),
		},
		{
			name:           "StatusPass is set",
			userDataStatus: "",
			statusFail:     false,
			statusInfo:     false,
			statusPass:     true,
			expectedStatus: beacon.Status(beacon.STATUS_PASS),
		},
		{
			name:           "StatusInfo is set",
			userDataStatus: "",
			statusFail:     false,
			statusInfo:     true,
			statusPass:     false,
			expectedStatus: beacon.Status(beacon.STATUS_INFO),
		},
		{
			name:           "No status is set",
			userDataStatus: "",
			statusFail:     false,
			statusInfo:     false,
			statusPass:     false,
			expectedStatus: beacon.Status(beacon.STATUS_INFO), // Default
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userDataStatus = tt.userDataStatus
			statusFail = tt.statusFail
			statusInfo = tt.statusInfo
			statusPass = tt.statusPass

			status := chooseStatusMessage()
			if status != tt.expectedStatus {
				t.Errorf("Expected status %s, got %s", tt.expectedStatus, status)
			}
		})
	}
}
