package main

var (
	userDataStatus string
	statusFail     bool
	statusInfo     bool
	statusPass     bool
	project        string
	permissions    bool
	configFile     string
	generateConfig bool
)

var requiredPermissions = []string{
	"events:PutEvents",
}
