package beacon

import (
	"fmt"
)

type Emitter struct {
	UserDataStatus string
	StatusFail     bool
	StatusInfo     bool
	StatusPass     bool
	InstanceID     string
	Project        string

	EBClient   EventBridgeClient
	IMDSClient IMDSClient
}

func (e Emitter) chooseStatusMessage() string {
	if e.UserDataStatus != "" {
		return e.UserDataStatus
	}

	if e.StatusFail {
		return STATUS_FAIL
	}

	if e.StatusPass {
		return STATUS_PASS
	}

	if e.StatusInfo {
		return STATUS_INFO
	}

	return STATUS_INFO
}

func (e Emitter) Emit(message string) error {
	if e.InstanceID == "" {
		instanceID, err := retrieveInstanceId(e.IMDSClient)
		e.InstanceID = instanceID
		if err != nil {
			return err
		}
	}

	status := e.chooseStatusMessage()
	detail := fmt.Sprintf("{ \"Status\": \"%s\", \"Message\": \"%s\"}", status, message)
	detailType := "User Data"

	err := sendEvent(e.EBClient, detail, detailType, e.InstanceID, e.Project)
	if err != nil {
		return err
	}

	return nil
}
