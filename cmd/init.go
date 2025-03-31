package main

import (
	"os"
)

func init() {
	rootCmd.Flags().StringVar(&instanceID, FlagInstanceIDLong, FlagInstanceIDDefault, FlagInstanceIDDesc)
	rootCmd.Flags().StringVar(&project, FlagProjectLong, FlagProjectDefault, FlagProjectDesc)
	rootCmd.Flags().StringVarP(&userDataStatus, FlagStatusLong, FlagStatusShort, FlagStatusDefault, FlagStatusDesc)
	rootCmd.Flags().BoolVarP(&statusFail, FlagFailLong, FlagFailShort, FlagFailDefault, FlagFailDesc)
	rootCmd.Flags().BoolVarP(&statusInfo, FlagInfoLong, FlagInfoShort, FlagInfoDefault, FlagInfoDesc)
	rootCmd.Flags().BoolVarP(&statusPass, FlagPassLong, FlagPassShort, FlagPassDefault, FlagPassDesc)
	rootCmd.Flags().BoolVar(&permissions, FlagPermissionsLong, FlagPermissionsDefault, FlagPermissionsDesc)
	rootCmd.Flags().StringVarP(&configFile, FlagConfigLong, FlagConfigShort, FlagConfigDefault, FlagConfigDesc)

	// Support environment variables for instance-id and project
	if os.Getenv(EnvBeaconInstanceID) != "" && instanceID == FlagInstanceIDDefault {
		instanceID = os.Getenv(EnvBeaconInstanceID)
	}

	if os.Getenv(EnvBeaconProject) != "" && project == FlagProjectDefault {
		project = os.Getenv(EnvBeaconProject)
	}
}
