/*
Package beacon provides functionality for sending EC2 user data events to
EventBridge, enabling real-time monitoring of EC2 instance startup
processes.

# Overview

During EC2 instance startup, cloud-init executes user data scripts to configure the machine
and deploy applications. If these scripts fail, instances may terminate without alerting
administrators, especially in autoscaling groups.

Beacon solves this by emitting custom events to EventBridge at critical points in your user
data execution, allowing:

  - Real-time monitoring of user data script execution
  - Alerting on failures
  - Tracking of deployment steps
  - Auditing of instance initialization

# Usage

Beacon can send three types of events:
  - Success (pass): Indicates successful completion of a step
  - Failure (fail): Indicates failure of a step
  - Informational (info): Provides status information
  - Custom: User-defined status

Example:

	emitter := beacon.Emitter{
		InstanceID: beacon.InstanceID("i-1234567890abcdef0"),
		Project:    beacon.Project("my-project"),
		EBClient:   cloudwatchevents.NewFromConfig(cfg),
	}

	// Send a success event
	err := emitter.Emit(ctx, beacon.Status(beacon.STATUS_PASS), "Configuration completed successfully")

	// Send a failure event
	err := emitter.Emit(ctx, beacon.Status(beacon.STATUS_FAIL), "Failed to download artifact")

	// Send an informational event
	err := emitter.Emit(ctx, beacon.Status(beacon.STATUS_INFO), "Starting deployment")

	// Send a custom status event
	err := emitter.Emit(ctx, beacon.Status("warning"), "Low disk space detected")

# Required IAM Permissions

The EC2 instance or execution environment requires the following IAM permission:

	{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Action": "events:PutEvents",
				"Resource": "*"
			}
		]
	}

# Environment Variables

Beacon recognizes the following environment variables:

  - BEACON_INSTANCE_ID: Override EC2 instance ID detection
  - BEACON_PROJECT: Set project name (default: "unknown")

# EventBridge Integration

Events emitted by Beacon can be captured in EventBridge using a pattern like:

	{
		"detailType": [
			"User Data"
		]
	}

Additional filtering options:
  - By project: "source": ["your-project-name"]
  - By status: "detail": {"Status": ["fail"]}
*/
package beacon
