package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Global variables demonstration
var (
	// ConfigFile holds the path to the configuration file
	configFile string

	// Debug mode flag
	debugMode bool

	// Version information
	version = "1.0.0"
)

// Constants demonstration
const (
	// Default configuration values
	defaultConfigPath = "configs/config.yaml"
	defaultBufferSize = 4096

	// File processing modes
	ModeSingle  = "single"
	ModeWatch   = "watch"
	ModeAnalyze = "analyze"
)

func main() {
	// Initialize the root command
	rootCmd := &cobra.Command{
		Use:   "analyzer",
		Short: "File Analytics System",
		Long: `A robust file processing and analytics system that demonstrates 
		various Go programming concepts while providing useful file analysis capabilities.`,
		Version: version,
		RunE:    run,
	}

	// Command-line flags demonstration
	rootCmd.PersistentFlags().StringVar(&configFile, "config", defaultConfigPath, "config file path")
	rootCmd.PersistentFlags().BoolVar(&debugMode, "debug", false, "enable debug mode")

	// Execute the root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// run implements the main logic of the application
// Demonstrates multiple return values
func run(cmd *cobra.Command, args []string) error {
	// Initialize configuration
	if err := initConfig(); err != nil {
		return fmt.Errorf("failed to initialize config: %w", err)
	}

	// Setup logging
	setupLogging()

	logrus.Info("File Analytics System started")
	return nil
}

// initConfig demonstrates error handling and file operations
func initConfig() error {
	if configFile != "" {
		// Get the absolute path
		absPath, err := filepath.Abs(configFile)
		if err != nil {
			return fmt.Errorf("failed to get absolute path: %w", err)
		}

		// Set the config file path
		viper.SetConfigFile(absPath)
	} else {
		// Search for config in default locations
		viper.AddConfigPath(".")
		viper.AddConfigPath("./configs")
		viper.SetConfigName("config")
	}

	// Read the config file
	if err := viper.ReadInConfig(); err != nil {
		// Demonstrates type assertion in error handling
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			logrus.Warn("No config file found, using defaults")
		} else {
			return fmt.Errorf("failed to read config file: %w", err)
		}
	}

	return nil
}

// setupLogging demonstrates conditional logic and configuration
func setupLogging() {
	// If/else demonstration
	if debugMode {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("Debug mode enabled")
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	// Customize logging format
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}
