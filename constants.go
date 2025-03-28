package beacon

import (
	"time"
)

const (
	AppName      string = "beacon"
	AppDescShort string = "Allows user data to emit custom EventBridge events during processing"
	AppDescLong  string = "Allows user data to emit custom EventBridge events during processing. Returns 0 on success, 1 on failure."
)

const KB int = 1024

const DEFAULT_TIMEOUT time.Duration = 30 * time.Second

const DEFAULT_DETAIL_TYPE string = "User Data Beacon"

// AWS length/size constraints
// Ref: https://docs.aws.amazon.com/eventbridge/latest/APIReference/API_PutEvents.html
// Ref: https://docs.aws.amazon.com/eventbridge/latest/APIReference/API_PutEventsRequestEntry.html
const (
	EVENT_PAYLOAD_MAX_BYTES int = 256 * KB
	DETAIL_TYPE_MAX_LENGTH  int = 128
	RESOURCE_ARN_MAX_LENGTH int = 2048
)

const (
	STATUS_FAIL string = "fail"
	STATUS_INFO string = "info"
	STATUS_PASS string = "pass"
)
