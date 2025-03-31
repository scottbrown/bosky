package main

var (
	userDataStatus string
	statusFail     bool
	statusInfo     bool
	statusPass     bool
	instanceID     string
	project        string
	permissions    bool
	configFile     string
)

var requiredPermissions = []string{
	"events:PutEvents",
}
