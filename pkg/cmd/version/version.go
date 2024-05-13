package version

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/spf13/cobra"
	ver "mcgaunn.com/logwild/pkg/version"
)

var (
	versionCmdUse   string = "version"
	versionCmdShort string = "get version"
	versionCmdLong  string = "get version of logwild"
	versionString   string = ver.VERSION
)

func NewVersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   versionCmdUse,
		Short: versionCmdShort,
		Long:  versionCmdLong,
		RunE:  doRunVersionCmd,
	}

	return cmd
}

func doRunVersionCmd(cmd *cobra.Command, args []string) error {
	slog.Debug("got request to run version command", "args", args)
	fmt.Fprintf(os.Stdout, "%s\n", versionString)
	return nil
}
