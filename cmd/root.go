package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"jenv-go/internal/config"
	"jenv-go/internal/env"
	"jenv-go/internal/java"
	"jenv-go/internal/sys"
	"os"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "jenv",
	Short: "A Java version manager",
	Long: `Jenv is a command-line tool for managing multiple Java versions.

It allows you to easily switch between different Java versions,
add new Java installations, and manage your Java environment.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Initialize configuration
		_, err := config.GetInstance()
		if err != nil {
			fmt.Printf("Error initializing configuration: %v\n", err)
			os.Exit(1)
		}

		// Backup environment variables if not already done
		if err := config.BackupEnvPath(); err != nil {
			fmt.Printf("Warning: Failed to backup environment variables: %v\n", err)
		}
	},
	// Add author and license information
	Version: "v2.0.0",
}

// Set author information and license
func init() {
	rootCmd.SetVersionTemplate(`{{with .Name}}{{printf "%s " .}}{{end}}{{printf "version %s" .Version}}
version: 0.1.0
Author: WhyWhatHow (https://github.com/WhyWhatHow)
Email: whywhathow.fun@gmail.com
License: Apache License 2.0
`)

	//Initialize configuration system
	cfg, err := config.GetInstance()
	if err != nil {
		fmt.Printf("Error initializing configuration: %v\n", err)
		os.Exit(1)
	}
	if !cfg.Initialized {
		//1.  backup environment variables
		java.Init()
		err := config.BackupEnvPath()

		if err != nil {
			fmt.Printf("Error backing up environment variables: %v\n", err)
			os.Exit(1)
		}
		//2. add JAVA_HOME,JAVA_HOME to PATH
		env.SetEnv("JAVA_HOME", cfg.SymlinkPath)

		cfg.Initialized = true
		err = cfg.Save()
	}
	// Add global flags
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

/*
*
1. Initialize configuration
2. Backup environment variables if not already done
3. Add global flags
4.
*/
func init() {
	if !sys.IsAdmin() {
		fmt.Errorf("please run this command with admin privileges")
		os.Exit(1)
	}
	//Initialize configuration system
	cfg, err := config.GetInstance()
	if err != nil {
		fmt.Printf("Error initializing configuration: %v\n", err)
		os.Exit(1)
	}
	if !cfg.Initialized {
		//1.  backup environment variables
		java.Init()
		err := config.BackupEnvPath()

		if err != nil {
			fmt.Printf("Error backing up environment variables: %v\n", err)
			os.Exit(1)
		}
		//2. add JAVA_HOME,JAVA_HOME to PATH
		env.SetEnv("JAVA_HOME", cfg.SymlinkPath)

		cfg.Initialized = true
		err = cfg.Save()
	}
	// Add global flags
}
