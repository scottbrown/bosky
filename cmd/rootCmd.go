package main

import (
	"context"
	"fmt"
	"os"

	"github.com/scottbrown/beacon"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/ec2/imds"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchevents"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     fmt.Sprintf("%s [message]", beacon.AppName),
	Short:   beacon.AppDescShort,
	Long:    beacon.AppDescLong,
	Example: fmt.Sprintf("%s --fail \"Artifact download returned 404\"", beacon.AppName),
	RunE:    handleRoot,
	PreRunE: handleFlagValidation,
	Version: beacon.VERSION,
}

// handleRoot processes the main command logic
func handleRoot(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithTimeout(context.Background(), beacon.DEFAULT_TIMEOUT)
	defer cancel()

	if permissions {
		return handleListPermissions(ctx)
	}

	if generateConfig {
		return handleGenerateConfig(ctx)
	}

	if len(args) < 1 {
		return fmt.Errorf("missing [message]")
	}

	// Load beacon configuration
	beaconConfig, err := beacon.LoadConfig(configFile)
	if err != nil && configFile != "" {
		return fmt.Errorf("failed to load config file: %w", err)
	}

	if beaconConfig != nil {
		if instanceIDNotProvided() && beaconConfig.InstanceID != "" {
			instanceID = beaconConfig.InstanceID
		}

		if project == FlagProjectDefault && beaconConfig.Project != "" {
			project = beaconConfig.Project
		}
	}

	message := args[0]

	cfg, err := config.LoadDefaultConfig(ctx, config.WithRetryMaxAttempts(3))
	if err != nil {
		return err
	}

	if instanceIDNotProvided() {
		instanceID, err = beacon.RetrieveInstanceID(ctx, imds.NewFromConfig(cfg))
		if err != nil {
			return err
		}
	}

	emitter := beacon.Emitter{
		InstanceID: beacon.InstanceID(instanceID),
		Project:    beacon.Project(project),
		EBClient:   cloudwatchevents.NewFromConfig(cfg),
	}

	status := chooseStatusMessage()

	err = emitter.Emit(ctx, status, message)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}

	return nil
}

func instanceIDNotProvided() bool {
	return instanceID == FlagInstanceIDDefault
}
