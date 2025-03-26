![Beacon](beacon.small.png)

# Beacon - EC2 User Data Event Emitter

Beacon sends events from EC2 user data scripts to AWS CloudWatch Events, enabling real-time monitoring of startup processes.

## Overview

During EC2 instance startup, [cloud-init](https://cloud-init.io/) executes user data scripts to configure the machine and deploy applications. If these scripts fail, instances may terminate without alerting administrators, especially in autoscaling groups.

Beacon solves this by emitting custom events to CloudWatch at critical points in your user data execution, allowing:
- Real-time monitoring of user data script execution
- Alerting on failures
- Tracking of deployment steps
- Auditing of instance initialization

## Installation

### Binary Installation

1. Download the appropriate binary for your architecture from the [Releases page](https://github.com/scottbrown/beacon/releases)
2. Copy to your instance: `sudo cp beacon /usr/local/bin/`
3. Make executable: `sudo chmod +x /usr/local/bin/beacon`

### Required IAM Permission

EC2 instances using Beacon require the `events:PutEvents` permission:

```json
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
```

Add this to your instance profile or role.

## Usage

Beacon can send three types of events: success (pass), failure (fail), or informational (info).

### Basic Usage

```bash
# Send success event
beacon --pass "User data processed successfully"

# Send failure event
beacon --fail "Failed to download application artifact"

# Send informational event
beacon --info "Starting application deployment"

# Custom status
beacon --status "warning" "Disk space below threshold"
```

### Environment Variables

- `BEACON_INSTANCE_ID`: Override instance ID detection
- `BEACON_PROJECT`: Set project name (default: "unknown")

### Command Line Options

```
Usage:
  beacon [message] [flags]

Flags:
  -f, --fail            Emits a failure event
  -h, --help            Help for beacon
      --info            Emits an informational event
      --instance-id     Specifies the EC2 INSTANCE_ID instead of looking it up
  -p, --pass            Emits a successful event
      --project string  Names the PROJECT as a source for the event (default "unknown")
      --status string   Emits an event with a custom STATUS
  -v, --version         Version for beacon
```

## CloudWatch Events Setup

### Creating the Rule

1. Open the CloudWatch console
2. Navigate to Events → Rules
3. Create a new rule
4. For the event pattern:

```json
{
  "detailType": [
    "User Data"
  ]
}
```

5. Additional filtering options:
   - By project: `"source": ["your-project-name"]`
   - By status: `"detail": {"Status": ["fail"]}`

### Target Configuration

Connect your rule to targets like:
- SNS topics for email/SMS notifications
- Lambda functions for custom processing
- SQS queues for event processing
- CloudWatch Alarm actions

## Integration Examples

### Basic User Data Success Tracking

```bash
#!/bin/bash

# Start user data execution
beacon --info "Starting user data execution"

# Install dependencies
apt-get update
if [ $? -eq 0 ]; then
  beacon --info "System packages updated"
else
  beacon --fail "Failed to update system packages"
  exit 1
fi

# Deploy application
./deploy_app.sh
if [ $? -eq 0 ]; then
  beacon --pass "Application deployed successfully"
else
  beacon --fail "Application deployment failed"
  exit 1
fi
```

### Project-Based Tracking

```bash
#!/bin/bash
export BEACON_PROJECT="webapp-fleet"

beacon --info "Starting webapp deployment"
# Deployment steps...
beacon --pass "Webapp successfully deployed"
```

## Building From Source

Prerequisites:
- Go 1.24+
- [Task](https://taskfile.dev) (v3.x)

```bash
# Clone the repository
git clone https://github.com/scottbrown/beacon.git
cd beacon

# Build
task build

# Run tests
task test

# Build for all platforms
task dist

# Create release artifacts (requires VERSION env var)
task release VERSION=1.1.0
```

You can list all available tasks with:

```bash
task
```

## Cost Considerations

CloudWatch Events has associated costs:
- $1.00 per million custom events
- Additional costs for targets (SNS, Lambda, etc.)

While costs are minimal for most use cases, be mindful when implementing at scale.

## License

MIT License - See [LICENSE](LICENSE) for details.

Copyright © Scott Brown
