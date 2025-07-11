package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/samuelea/gator/internal/database"
)

type State struct {
	Config *Config
	DbQueries *database.Queries
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	Handlers map[string]func(*State, Command) error
}

type Config struct {
	DBUrl string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

const configfileName = ".gatorconfig.json"

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", homeDir, configfileName), nil
}

func Read() (*Config, error) {
	configPath, err := getConfigFilePath()

	if err != nil {
		return nil, err
	}

	file, err := os.ReadFile(string(configPath))

	if err != nil {
		return nil, err
	}

	var gatorConfig Config

	err = json.Unmarshal(file, &gatorConfig)

	if err != nil {
		return nil, err
	}

	return &gatorConfig, err
}

func (c *Config) SetUser(user string) error {
	c.CurrentUserName = user

	err := write(c)

	if err != nil {
		return err
	}

	return nil
}

func write(c *Config) error {
	configPath, err := getConfigFilePath()

	if err != nil {
		return err
	}

	newFile, err := json.Marshal(&c)

	if err != nil {
		return err
	}

	err = os.WriteFile(configPath, newFile, 0644)

	if err != nil {
		return err
	}

	return nil
}

func (c *Commands) Run(state *State, command Command) error {
	handler, ok := c.Handlers[command.Name]

	if !ok {
		return fmt.Errorf("error: command %s not found", command.Name)
	}

	return handler(state, command)
}
