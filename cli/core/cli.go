package core

import (
	"github.com/PARTHIB-DEB/TermoTalks/cli/cmd"
	"github.com/spf13/cobra"
)

func Callcli() {
	var rootCmd = &cobra.Command{
		Use:   "TT [Commands] [Args] --flags",
		Short: "TT or TermoTalks is a platform to chat within the terminal",
	}
	// Preserve the original os.Args
	// originalArgs := os.Args
	// It assumes "TT" is the program name, adjust if "TT" comes from a different source
	// os.Args = []string{"TT"} // You might need to adjust "TT" if your base command name is different
	rootCmd.AddCommand(cmd.ProfCmd)
	rootCmd.AddCommand(cmd.LinkCmd)
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
	// os.Args = originalArgs
}
