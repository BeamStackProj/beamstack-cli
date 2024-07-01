package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Beamflow/beamflow-cli/src/types"
	"gopkg.in/yaml.v2"
)

func GetProfile(configFile string) (profile types.Profiles, err error) {
	// Check if the config file exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return profile, fmt.Errorf("config file not found at %s", configFile)
	} else if err != nil {
		return profile, fmt.Errorf("error checking config file: %v", err)
	}

	// Get file extension
	ext := filepath.Ext(configFile)

	// Check if the file type is supported
	switch ext {
	case ".yaml", ".yml":
		err = parseYAML(configFile, &profile)
	case ".json":
		err = parseJSON(configFile, &profile)
	default:
		return profile, fmt.Errorf("unsupported file type: %s", ext)
	}

	return profile, err
}

func parseYAML(filePath string, profile *types.Profiles) error {
	// Open the YAML file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// Read file contents
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	// Unmarshal YAML data into profile struct
	if err := yaml.Unmarshal(data, profile); err != nil {
		return fmt.Errorf("error parsing YAML: %v", err)
	}

	return nil
}

func parseJSON(filePath string, profile *types.Profiles) error {
	// Open the JSON file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// Read file contents
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	// Unmarshal JSON data into profile struct
	if err := json.Unmarshal(data, profile); err != nil {
		return fmt.Errorf("error parsing JSON: %v", err)
	}

	return nil
}

func SaveProfile(profile *types.Profiles) error {

	jsonData, err := json.MarshalIndent(profile, "", "    ")
	if err != nil {
		return fmt.Errorf("Error marshaling config to JSON, %s", err)
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("could not locate home directory %s", err)
	}

	fileName := filepath.Join(homeDir, ".beamstack", "profiles", fmt.Sprintf("%s.json", profile.Name))
	// Write the JSON data to a file
	err = os.WriteFile(fileName, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("Error writing profile file, %s", err)
	}
	return nil
}
