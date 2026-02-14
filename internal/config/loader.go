package config

import (
	"flag"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"os"
)

const (
	commonConfigFlagName    = "config"
	commonConfigFlagUsage   = "path to config file"
	defaultCommonConfigPath = "./config/config.yaml"

	environmentConfigFlagName  = "env-config"
	environmentConfigFlagUsage = "path to environment config file"

	envVarCommonConfigPath      = "PXR_SSO_CONFIG_PATH"
	envVarEnvironmentConfigPath = "PXR_SSO_ENV_CONFIG_PATH"
)

type Loader struct {
	cfg Config

	commonConfigPath      string
	environmentConfigPath string
}

func NewLoader() *Loader {
	return &Loader{}
}

// MustLoad loads config and panics if any error occurs.
func (l *Loader) MustLoad() *Config {
	l.fetchPath()

	if err := l.loadCommonConfig(); err != nil {
		panic(err)
	}

	_ = l.loadEnvironmentConfig()

	return &l.cfg
}

// fetchPath fetches path from command line flag or environment variable.
// Priority: flag > env > default.
func (l *Loader) fetchPath() {
	l.flagParse()
	l.envVarParse()
}

func (l *Loader) flagParse() {
	var commonConfigPath string
	var environmentConfigPath string
	flag.StringVar(&commonConfigPath, commonConfigFlagName, defaultCommonConfigPath, commonConfigFlagUsage)
	flag.StringVar(&environmentConfigPath, environmentConfigFlagName, "", environmentConfigFlagUsage)
	flag.Parse()

	l.setCommonConfigPathIfEmpty(commonConfigPath)
	l.setEnvironmentConfigPathIfEmpty(environmentConfigPath)
}

func (l *Loader) envVarParse() {
	_ = godotenv.Load()

	commonConfigPath := os.Getenv(envVarCommonConfigPath)
	l.setCommonConfigPathIfEmpty(commonConfigPath)

	environmentConfigPath := os.Getenv(envVarEnvironmentConfigPath)
	l.setEnvironmentConfigPathIfEmpty(environmentConfigPath)
}

func (l *Loader) setCommonConfigPathIfEmpty(path string) {
	if l.commonConfigPath == "" {
		l.commonConfigPath = path
	}
}

func (l *Loader) setEnvironmentConfigPathIfEmpty(path string) {
	if l.environmentConfigPath == "" {
		l.environmentConfigPath = path
	}
}

func (l *Loader) loadCommonConfig() error {
	if l.commonConfigPath == "" {
		return fmt.Errorf("common config path is empty")
	}

	if err := l.loadByPath(l.commonConfigPath); err != nil {
		return err
	}

	return nil
}

func (l *Loader) loadEnvironmentConfig() error {
	if l.environmentConfigPath == "" {
		return fmt.Errorf("environment config path is empty")
	}

	if err := l.loadByPath(l.environmentConfigPath); err != nil {
		return err
	}

	return nil
}

func (l *Loader) loadByPath(path string) error {
	file, err := os.Stat(path)
	if os.IsNotExist(err) {
		return fmt.Errorf("%s: %s", "config file does not exists", path)
	}

	if err := cleanenv.ReadConfig(path, &l.cfg); err != nil {
		return fmt.Errorf("%s %s: %w", "cannot read", file.Name(), err)
	}

	return nil
}
