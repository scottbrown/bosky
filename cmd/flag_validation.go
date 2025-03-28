package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func handleFlagValidation(cmd *cobra.Command, args []string) error {
	if permissions {
		return nil
	}

	return validateRequiredFlags(cmd)
}

func validateRequiredFlags(cmd *cobra.Command) error {
	if !statusFail && !statusInfo && !statusPass && userDataStatus == "" {
		return fmt.Errorf("At least one of --fail, --info, --pass, or --status is required")
	}

	return nil
}
