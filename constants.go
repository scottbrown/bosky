package beacon

import (
	"time"
)

const (
	// AppName is the name of the application
	AppName string = "beacon"
	// AppDescShort is a short description of the application
	AppDescShort string = "Allows user data to emit custom EventBridge events during processing"
	// AppDescLong is a detailed description of the application
	AppDescLong string = "Allows user data to emit custom EventBridge events during processing. Returns 0 on success, 1 on failure."
)

const KB int = 1024

// DEFAULT_TIMEOUT defines the maximum time an operation can take
const DEFAULT_TIMEOUT time.Duration = 30 * time.Second

// DEFAULT_DETAIL_TYPE defines the default EventBridge detail type
const DEFAULT_DETAIL_TYPE string = "User Data Beacon"

// AWS length/size constraints
// Ref: https://docs.aws.amazon.com/eventbridge/latest/APIReference/API_PutEvents.html
// Ref: https://docs.aws.amazon.com/eventbridge/latest/APIReference/API_PutEventsRequestEntry.html
const (
	EVENT_PAYLOAD_MAX_BYTES int = 256 * KB
	DETAIL_TYPE_MAX_LENGTH  int = 128
	RESOURCE_ARN_MAX_LENGTH int = 2048
)

// Status constants for events
const (
	// STATUS_FAIL indicates a failure event
	STATUS_FAIL string = "fail"

	// STATUS_INFO indicates an informational event
	STATUS_INFO string = "info"

	// STATUS_PASS indicates a successful event
	STATUS_PASS string = "pass"
)
