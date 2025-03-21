package main

import (
	"context"
	"fmt"
	"os"

	"github.com/RaihanurRahman2022/file-analytics/internal/processor"
	"github.com/RaihanurRahman2022/file-analytics/pkg/utils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "analyzer",
	Short: "A file analysis tool",
	Long:  `A tool for analyzing files, calculating hashes, and encoding/decoding content.`,
}

// analyzeCmd represents the analyze command
var analyzeCmd = &cobra.Command{
	Use:   "analyze [path]",
	Short: "Analyze files in the specified path",
	Long: `Analyze files in the specified path, processing them according to their type.
	Supports multiple file formats including text, JSON, and CSV.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("path argument is required")
		}

		path := args[0]
		if _, err := os.Stat(path); os.IsNotExist(err) {
			return fmt.Errorf("path does not exist: %s", path)
		}

		// Create processors
		processors := []processor.Processor{
			processor.NewTextProcessor(4096),
			processor.NewJSONProcessor(4096),
			processor.NewCSVProcessor(4096),
		}

		// Process files
		return processFiles(path, processors)
	},
}

// hashCmd represents the hash command
var hashCmd = &cobra.Command{
	Use:   "hash [file]",
	Short: "Calculate SHA256 hash of a file",
	Long:  `Calculate and display the SHA256 hash of the specified file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("file argument is required")
		}

		hash, err := utils.HashFile(args[0])
		if err != nil {
			return fmt.Errorf("failed to calculate hash: %w", err)
		}

		fmt.Printf("SHA256: %s\n", hash)
		return nil
	},
}

// encodeCmd represents the encode command
var encodeCmd = &cobra.Command{
	Use:   "encode [file]",
	Short: "Base64 encode a file",
	Long:  `Base64 encode the contents of the specified file and display the result.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return fmt.Errorf("file argument is required")
		}

		encoded, err := utils.Base64EncodeFile(args[0])
		if err != nil {
			return fmt.Errorf("failed to encode file: %w", err)
		}

		fmt.Printf("Base64: %s\n", encoded)
		return nil
	},
}

// decodeCmd represents the decode command
var decodeCmd = &cobra.Command{
	Use:   "decode [base64] [output]",
	Short: "Base64 decode to a file",
	Long:  `Decode base64 content and write it to the specified output file.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return fmt.Errorf("base64 content and output file arguments are required")
		}

		if err := utils.Base64DecodeFile(args[0], args[1]); err != nil {
			return fmt.Errorf("failed to decode file: %w", err)
		}

		fmt.Printf("Decoded content written to: %s\n", args[1])
		return nil
	},
}

// processFiles processes files in the given path using the provided processors
func processFiles(path string, processors []processor.Processor) error {
	// Create file filter
	filter := utils.CreateExtensionFilter(".txt", ".json", ".csv", ".tsv")

	// Walk through files
	return utils.WalkFiles(path, filter, func(filePath string) error {
		// Find appropriate processor
		var selectedProcessor processor.Processor
		for _, p := range processors {
			if p.CanHandle(filePath) {
				selectedProcessor = p
				break
			}
		}

		if selectedProcessor == nil {
			logrus.Warnf("No processor found for file: %s", filePath)
			return nil
		}

		// Process file
		result, err := selectedProcessor.Process(context.Background(), filePath)
		if err != nil {
			logrus.Errorf("Failed to process file %s: %v", filePath, err)
			return nil
		}

		// Log results
		logrus.Infof("Processed %s: %d lines, %d words, %d bytes in %v",
			filePath, result.Lines, result.Words, result.Bytes, result.Duration)

		return nil
	})
}

func init() {
	rootCmd.AddCommand(analyzeCmd)
	rootCmd.AddCommand(hashCmd)
	rootCmd.AddCommand(encodeCmd)
	rootCmd.AddCommand(decodeCmd)
}

func Execute() error {
	return rootCmd.Execute()
}
