package main

import (
	"github.com/scottbrown/beacon"
)

// chooseStatusMessage determines which status to use based on command flags
// Priority order: custom status, fail, pass, info
// Default is info if no status flags are set
func chooseStatusMessage() beacon.Status {
	if userDataStatus != "" {
		return beacon.Status(userDataStatus)
	}

	if statusFail {
		return beacon.Status(beacon.STATUS_FAIL)
	}

	if statusPass {
		return beacon.Status(beacon.STATUS_PASS)
	}

	if statusInfo {
		return beacon.Status(beacon.STATUS_INFO)
	}

	return beacon.Status(beacon.STATUS_INFO)
}
