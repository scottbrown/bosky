package main

import (
	"context"
	"fmt"
	"os"

	"github.com/scottbrown/bosky"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchevents"
	"github.com/spf13/cobra"
)

var (
	userDataStatus string
	statusFail     bool
	statusInfo     bool
	statusPass     bool
	instance_id    string
	project        string
)

func main() {
	rootCmd := &cobra.Command{
		Use:     "bosky [message]",
		Short:   "Allows user data to emit custom CloudWatch Events during processing",
		Long:    "Allows user data to emit custom CloudWatch Events during processing. Returns 0 on success, 1 on failure.",
		Example: "bosky --fail \"Artifact download returned 404\"",
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			message := args[0]

			ctx := context.TODO()

			// Load the AWS SDK configuration
			cfg, err := config.LoadDefaultConfig(ctx, config.WithRetryMaxAttempts(3))
			if err != nil {
				return err
			}

			emitter := bosky.Emitter{
				UserDataStatus: userDataStatus,
				StatusFail:     statusFail,
				StatusInfo:     statusInfo,
				StatusPass:     statusPass,
				InstanceID:     instance_id,
				Project:        project,
				EBClient:       cloudwatchevents.NewFromConfig(cfg),
				IMDSClient:     imds.NewFromConfig(cfg),
			}

			err = emitter.Emit(message)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return err
			}

			return nil
		},
	}

	// Add flags matching the previous CLI options
	rootCmd.Flags().StringVar(&instance_id, "instance-id", "", "Specifies the EC2 INSTANCE_ID instead of looking it up with the metadata service")
	rootCmd.Flags().StringVar(&project, "project", "unknown", "Names the PROJECT as a source for the event")
	rootCmd.Flags().StringVar(&userDataStatus, "status", "", "Emits an event with a custom STATUS")
	rootCmd.Flags().BoolVarP(&statusFail, "fail", "f", false, "Emits a failure event")
	rootCmd.Flags().BoolVarP(&statusInfo, "info", "i", false, "Emits an informational event")
	rootCmd.Flags().BoolVarP(&statusPass, "pass", "p", false, "Emits a successful event")

	// Support environment variables for instance-id and project
	if os.Getenv("BOSKY_INSTANCE_ID") != "" && instance_id == "" {
		instance_id = os.Getenv("BOSKY_INSTANCE_ID")
	}

	if os.Getenv("BOSKY_PROJECT") != "" && project == "unknown" {
		project = os.Getenv("BOSKY_PROJECT")
	}

	// Add author info
	rootCmd.Version = bosky.VERSION
	rootCmd.SetVersionTemplate("bosky version {{.Version}}\nAuthor: Scott Brown\n")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
