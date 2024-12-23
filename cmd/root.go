package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "fast-ansible",
	Short: "A fast tool for running ansible playbooks",
	Long:  `fast-ansible is a command line tool that leverages Go'sconcurrency and the Mitogen plugin to execute Ansible playbooks efficiently.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Fast-Ansible,More Efficient,More Faster!")
	},
}

func Execute() error {
	return rootCmd.Execute()
}
