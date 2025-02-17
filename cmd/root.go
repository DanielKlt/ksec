package cmd

import (
	"fmt"
	"os"

	"github.com/DanielKlt/ksec/cmd/secrets"
	"github.com/DanielKlt/ksec/pkg/environment"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "ksec",
	Aliases: []string{"ks"},
	Short:   "cli to work with secrets in k8s",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		var err error
		if err = environment.GetEnv().GetViper().BindPFlag("namespace", cmd.PersistentFlags().Lookup("namespace")); err != nil {
			return err
		}
		if err = environment.GetEnv().GetViper().BindPFlag("kubeconfig", cmd.PersistentFlags().Lookup("kubeconfig")); err != nil {
			return err
		}
		return err
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(VersionCmd)
	rootCmd.AddCommand(secrets.GetCmd)
}
