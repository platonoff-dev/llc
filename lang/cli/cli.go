package cli

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"llc/lang/files"
	"llc/lang/object"
	"llc/lang/repl"
)

func init() {
	RootCmd.AddCommand(RunCmd)
}

var RootCmd = &cobra.Command{
	Use:   "llc",
	Short: "llc is a Learning Language Compiler",
	Long:  "llc is a Learning Language Compiler",
}

var RunCmd = &cobra.Command{
	Use:   "run [module]",
	Short: "run a Learning Language Compiler",
	Long:  "run a Learning Language Compiler",
	Args:  cobra.MaximumNArgs(1),
	Run:   runCommand,
}

func runCommand(command *cobra.Command, args []string) {
	if len(args) == 0 {
		repl.Start(os.Stdin, os.Stdout)
	} else {
		env := object.NewEnvironment()
		_, err := files.ReadFile(args[0], env)
		if err != nil {
			log.Fatal(err)
		}
	}
}
