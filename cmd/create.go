package cmd

import (
	"github.com/dhillondeep/jailor/pkg/jailor"
	"github.com/dhillondeep/jailor/pkg/spec"
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
			parsed, err := spec.Parse(args[0])
			if err != nil {
				panic(err)
			}

			useCow, _ := cmd.Flags().GetBool("cow")
			force, _ := cmd.Flags().GetBool("force")
			clean, _ := cmd.Flags().GetBool("clean")

			if err := jailor.CreateJail(jailor.JailContext{
				UseCow: useCow,
				Force:  force,
				Clean:  clean,
			}, parsed, args[1]); err != nil {
				panic(err)
			}
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
