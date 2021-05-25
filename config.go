package main

import (
	"encoding/json"
	"errors"
	"fmt"
	homedir "github.com/mitchellh/go-homedir"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
)

const (
	globalConfigVersion    = "1"
	globalConfigDir        = ".config/td"
	globalConfigWindowsDir = "td\\"
	globalConfigFile       = "td.json"
)

var (
	// ErrAliasNotFound is returned when the account alias does not exist
	ErrAliasNotFound = errors.New("alias does not exist")
	// ErrUnableLoadConfig is returned when the config cannot be loaded from config dir
	ErrUnableLoadConfig = errors.New("unable to load config from config directory")
	// ErrUnableSafeConfig is returned when the config cannot be safed to config dir
	ErrUnableSafeConfig = errors.New("unable to safe config from config directory")
	// ErrUnableLocateHomeDir is returned when the home dir cannot be located
	ErrUnableLocateHomeDir = errors.New("unable to locate home directory")
	// ErrInvalidArgument is returned when invalid arguments have been passed to function
	ErrInvalidArgument = errors.New("argument is invalid or was nil")
)

// customConfigDir contains the whole path to config dir. Only access via get/set functions.
var customConfigDir string

// SetConfigDir sets a custom go-t client config folder.
func SetConfigDir(configDir string) {
	customConfigDir = configDir
}

// GetConfigDir constructs go-t client config folder.
func GetConfigDir() (string, error) {
	if customConfigDir != "" {
		return customConfigDir, nil
	}
	homeDir, e := homedir.Dir()
	if e != nil {
		return "", ErrUnableLocateHomeDir
	}
	var configDir string
	// For windows the path is slightly different
	if runtime.GOOS == "windows" {
		configDir = filepath.Join(homeDir, globalConfigWindowsDir)
	} else {
		configDir = filepath.Join(homeDir, globalConfigDir)
	}
	return configDir, nil
}

// GetConfigPath onstructs go-t client configuration path.
func GetConfigPath() (string, error) {
	if customConfigDir != "" {
		return filepath.Join(customConfigDir, globalConfigFile), nil
	}
	dir, err := GetConfigDir()
	if err != nil {
		return "", ErrUnableLocateHomeDir
	}
	return filepath.Join(dir, globalConfigFile), nil
}

// SaveConfig saves the configuration file and returns error if any.
func ToFile(list *Todolist) error {
	if list == nil {
		return ErrInvalidArgument
	}

	err := createConfigDir()
	if err != nil {
		return err
	}

	// Save the config.
	file, err := GetConfigPath()
	if err != nil {
		return err
	}

	b, err := json.Marshal(list)
	if err != nil {
		return fmt.Errorf("%s, %s", ErrUnableSafeConfig, err)
	}

	err = ioutil.WriteFile(file, b, 0700)
	if err != nil {
		return fmt.Errorf("%s, %s", ErrUnableSafeConfig, err)
	}

	return nil
}

// LoadConfig loads the config file
func FromFile() (*Todolist, error) {
	file, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadFile(file)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("%s, %s", ErrUnableLoadConfig, err)
	}

	// Load list
	var list = Todolist{
		Todos: []*Todo{},
	}
	err = json.Unmarshal(b, &list)
	if err != nil {
		return nil, fmt.Errorf("%s, %s", ErrUnableLoadConfig, err)
	}

	return &list, nil
}

// createConfigDir creates config folder
func createConfigDir() error {
	p, err := GetConfigDir()
	if err != nil {
		return err
	}
	if e := os.MkdirAll(p, 0700); e != nil {
		return err
	}
	return nil
}

// IsConfigExists returns false if config doesn't exist.
func IsConfigExists() bool {
	configFile, err := GetConfigPath()
	if err != nil {
		return false
	}
	if _, e := os.Stat(configFile); e != nil {
		return false
	}
	return true
}
