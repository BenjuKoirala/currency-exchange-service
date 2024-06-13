package cmd

import (
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "exchange-cli",
	Short: "A CLI for currency exchange rates",
}

var GrpcServerEndpoint string

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Printf("Error while executing root command %v", err.Error())
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(loadConfig)
}

// loadConfig reads in config file
func loadConfig() {
	// Use config file from the flag.
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	configPath := filepath.Join(".", "exchange-cli")
	viper.AddConfigPath(configPath)
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file : %v", err.Error())
	}
	// Retrieve the gRPC server endpoint from config file
	GrpcServerEndpoint = viper.GetString("grpc.server_endpoint")
	if GrpcServerEndpoint == "" {
		log.Fatal("gRPC server endpoint not specified in config file")
	}
}
