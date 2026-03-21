package cmd

import (
	"errors"
	"fmt"
	"go-admin-api/common/global"
	"os"

	"github.com/go-admin-team/go-admin-core/sdk/pkg"

	"github.com/spf13/cobra"

	"go-admin-api/cmd/config"
	"go-admin-api/cmd/version"
)

var rootCmd = &cobra.Command{
	Use:          "go-admin-api",
	Short:        "go-admin-api",
	SilenceUsage: true,
	Long:         `go-admin-api`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			tip()
			return errors.New(pkg.Red("requires at least one arg"))
		}
		return nil
	},
	PersistentPreRunE: func(*cobra.Command, []string) error { return nil },
	Run: func(cmd *cobra.Command, args []string) {
		tip()
	},
}

func tip() {
	usageStr := `Welcome to ` + pkg.Green(`go-admin-api `+global.Version) + `. Use ` + pkg.Red(`-h`) + ` to view available commands.`
	usageStr1 := `You can also refer to https://doc.go-admin.dev/guide/ksks for more information.`
	fmt.Printf("%s\n", usageStr)
	fmt.Printf("%s\n", usageStr1)
}

func init() {
	rootCmd.AddCommand(version.StartCmd)
	rootCmd.AddCommand(config.StartCmd)
}

//Execute : apply commands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}