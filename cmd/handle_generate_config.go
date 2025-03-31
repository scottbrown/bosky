package main

import (
	"context"
	"fmt"

	"github.com/scottbrown/beacon"

	"gopkg.in/yaml.v3"
)

func handleGenerateConfig(ctx context.Context) error {
	configFile := beacon.Config{
		InstanceID: "i-abcdef1234",
		Project:    "example",
	}

	template, err := yaml.Marshal(configFile)
	if err != nil {
		return err
	}

	fmt.Println(string(template))

	return nil
}
