# Bosky - User Data Event Emitter

Bosky is a way to send events from user data to AWS CloudWatch Events
which, in turn, can be consumed by any other AWS resources.

## Background

During an EC2 machine's startup sequence,
[cloud-init](https://cloud-init.io/) runs a special piece
of code called user-data.  Within user data, many companies often insert
Bash scripts to configure machines with runtime settings or deploy
applications.  If user data fails, it can often mean that the system did
not start correctly and, if it's behind an autoscaling group, that server
disappears without an administrator being warned that user data failed.

Using `bosky` while processing user data means that anyone interested
in the events within user data can be notified in _near_ real time.

## Important -- Costs

Using this application will incur costs.  CloudWatch Events is incredibly
cheap, but there is a cost associated with sending custom events through
this service, as well as the resources that consume the event.  Please
be aware of this when using this tool.

## Agent/System Installation

The binary is a self-contained, statically linked binary.  Once you copy
the binary to your system, that is all that it required.

### Install the Binary into the System

1. Go to the Releases page
1. Choose the latest stable version
1. Choose the binary for your architecture
1. Copy the binary to your system into `/usr/local/bin`.
1. Make sure `root` can find this binary in its `PATH`.
### Add an IAM Policy to Your EC2 Machine

Your EC2 machine's instance profile will require the `events:PutEvents`
IAM permission.

### Test It Out

And that's it.  Try testing it out.  You should receive an exit code of `0`
if the notification succeeds.

Since there is nothing on the other end of the event bus to consume the
event, it will be dropped.  This is the next installation step.

## Consumer Installation

A consumer is required to consume the user data events and do something
with them.  Consumers are parties that are interested in the user data
notifications from your servers.

The following installation will consume all user data events.  You can
filter as much or as little as you want with further rules.

1. Create a CloudWatch Event Rule.
1. Specify a rule pattern
  ```
  {
    "detailType": [
      "User Data"
    ]
  }
  ```
1. Create a target (for testing, use an SNS topic linked to your email address)

Now send a test event and you should receive an email with the details
of the user data event.

## Usage

We are going to announce to AWS that our system has started correctly
and finished processing user data successfully.  At the end of the user
data sequence, add the following command:

```bash
bosky --success "User data processed successfully"
```

When the system starts and hits this line, an event is sent to CloudWatch
Events and onto any matching CloudWatch Event rules.

## Testing

While this tool is intended for use on an EC2 machine, you can run it
from your home machine.  Make sure that you have your environment
configured with a working `AWS_ACCESS_KEY_ID` and `AWS_SECRET_ACCESS_KEY`
as well your IAM user has as the `events:PutEvents` permission.

## License

tl;dr MIT License

See [LICENSE](LICENSE) for more information.

