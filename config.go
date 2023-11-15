package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Clusters []Cluster `yaml:"clusters"`
}

func CreateConfig() (Config, error) {
	var config Config
	var err error
	err = os.Mkdir(configDir(), 0666)
	if err != nil {
		return config, err
	}

	file, err := os.Create(configFilePath())
	if err != nil {
		return config, err
	}

	config = defaultConfig()
	stringConfig, err := yaml.Marshal(&config)
	if err != nil {
		return config, err
	}

	file.Write(stringConfig)

	return config, nil
}

func ReadConfig() Config {
	f, err := os.ReadFile(configFilePath())
	handleErr(err)

	var config Config
	err = yaml.Unmarshal(f, &config)
	handleErr(err)

	return config
}

func WriteConfig(c Config) (Config, error) {
	data, err := yaml.Marshal(&c)
	if err != nil {
		return c, err
	}

	err = os.WriteFile(configFilePath(), data, 0666)
	if err != nil {
		return c, err
	}

	return c, nil
}

func EditConfigInEditor() {
	editor := os.Getenv("EDITOR")

	if editor == "" {
		editor = "vi"
	}

	cmd := exec.Command(editor, configFilePath())
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	err := cmd.Run()

	if err != nil {
		log.Fatal(err)
	}
}

func configDir() string {
	return fmt.Sprintf("%s/.config/oc-cluster", os.Getenv("HOME"))
}

func configFilePath() string {
	return fmt.Sprintf("%s/clusters", configDir())
}

func handleErr(e error) {
	if e != nil {
		log.Fatalf("Error [%s]", e)
	}
}

func defaultConfig() Config {
	return Config{
		Clusters: []Cluster{},
	}
}
