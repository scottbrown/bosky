package main

import (
	"github.com/scottbrown/beacon"
)

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
