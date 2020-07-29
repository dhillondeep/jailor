package cmd

import (
	"github.com/spf13/cobra"
)

var (
	createCmd = &cobra.Command{
		Use:   "create",
		Short: "Creates chroot/jail environment",
		Long: `This will create jail environment based on configuration provided.
It creates filesystem isolation and restricts daemons into filesystem sub-tree to enchance security.

Usage: jailor create [flags] <spec path> <dir path> `,
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			// TODO: parse spec

			// TODO: call create api
		},
	}
)

func init() {
	createCmd.PersistentFlags().BoolP("link", "l", false, "hard link files instead of copying")
	createCmd.PersistentFlags().BoolP("force", "f", false,
		"if an existing destination file connot be opened, remote it and try again")
	createCmd.PersistentFlags().BoolP("clean", "c", false,
		"remove existing destination file before attempting to open it")
	createCmd.PersistentFlags().BoolP("cow", "", false,
		"perform lightweight copies using CoW (copy on write)")
	rootCmd.AddCommand(createCmd)
}
