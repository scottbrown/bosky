package main

import (
  "gopkg.in/urfave/cli.v1"
  "os"
)

var (
  userDataStatus string
  statusFail bool
  statusInfo bool
  statusPass bool
  instance_id string
  project string
)

func main() {
  app := cli.NewApp()
  app.Name = "bosky"
  app.Usage = "Allows user data to emit custom CloudWatch Events during processing.  Returns 0 on success, 1 on failure."
  app.UsageText = "bosky --fail \"Artifact download returned 404\""
  app.Author = "Scott Brown"
  app.Action = func (c *cli.Context) error {
    if c.NArg() < 1 {
      cli.ShowAppHelp(c)
      return nil
    }

    message := c.Args().Get(0)

    err := emitEvent(message)

    if err != nil {
      return cli.NewExitError(err, 1)
    }

    return nil
  }

  app.Flags = []cli.Flag {
    cli.StringFlag{
      Name: "instance-id",
      Usage: "Specifies the EC2 `INSTANCE_ID` instead of looking it up with the metadata service",
      Destination: &instance_id,
      EnvVar: "BOSKY_INSTANCE_ID",
    },
    cli.StringFlag{
      Name: "project",
      Usage: "Names the `PROJECT` as a source for the event.",
      Value: "unknown",
      Destination: &project,
      EnvVar: "BOSKY_PROJECT",
    },
    cli.StringFlag{
      Name: "status",
      Usage: "Emits an event with a custom `STATUS`",
      Destination: &userDataStatus,
    },
    cli.BoolFlag{
      Name: "fail, f",
      Usage: "Emits a failure event",
      Destination: &statusFail,
    },
    cli.BoolFlag{
      Name: "info, i",
      Usage: "Emits an informational event",
      Destination: &statusInfo,
    },
    cli.BoolFlag{
      Name: "pass, p",
      Usage: "Emits a successful event",
      Destination: &statusPass,
    },
  }

  app.Run(os.Args)
}

