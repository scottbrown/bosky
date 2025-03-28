package main

const (
	EnvBeaconInstanceID string = "BEACON_INSTANCE_ID"
	EnvBeaconProject    string = "BEACON_PROJECT"
)

const (
	FlagInstanceIDLong    string = "instance-id"
	FlagInstanceIDDesc    string = "Specifies the EC2 instance ID instead of looking it up with IMDS"
	FlagInstanceIDDefault string = ""

	FlagProjectLong    string = "project"
	FlagProjectDesc    string = "Names the PROJECT as a source for the event"
	FlagProjectDefault string = "unknown"

	FlagStatusLong    string = "status"
	FlagStatusShort   string = "s"
	FlagStatusDesc    string = "Emits an event with a custom status"
	FlagStatusDefault string = ""

	FlagFailLong    string = "fail"
	FlagFailShort   string = "f"
	FlagFailDesc    string = "Emits a failure event"
	FlagFailDefault bool   = false

	FlagInfoLong    string = "info"
	FlagInfoShort   string = "i"
	FlagInfoDesc    string = "Emits an informational event"
	FlagInfoDefault bool   = false

	FlagPassLong    string = "pass"
	FlagPassShort   string = "p"
	FlagPassDesc    string = "Emits a successful event"
	FlagPassDefault bool   = false

	FlagPermissionsLong    string = "permissions"
	FlagPermissionsDesc    string = "Displays IAM permissions required by the application"
	FlagPermissionsDefault bool   = false
)
