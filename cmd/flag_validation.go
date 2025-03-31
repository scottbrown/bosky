package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

// handleFlagValidation ensures that at least one status flag is provided
// This is called before command execution to validate inputs
func handleFlagValidation(cmd *cobra.Command, args []string) error {
	if permissions || generateConfig {
		return nil
	}

	return validateRequiredFlags(cmd)
}

// validateRequiredFlags checks that at least one status flag is provided
func validateRequiredFlags(cmd *cobra.Command) error {
	if !statusFail && !statusInfo && !statusPass && userDataStatus == "" {
		return fmt.Errorf("at least one of --fail, --info, --pass, or --status is required")
	}

	return nil
}
