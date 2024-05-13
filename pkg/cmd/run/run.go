package run

import (
	"log/slog"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"mcgaunn.com/iankubetrace/pkg/api/http"
	"mcgaunn.com/iankubetrace/pkg/signals"
)

var (
	runCmdUse   string = "run"
	runCmdShort string = "run server"
	runCmdLong  string = "run server with the configured options"
)

func NewRunCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   runCmdUse,
		Short: runCmdShort,
		Long:  runCmdLong,
		RunE:  doRunCmd,
	}
	return cmd
}

func doRunCmd(cmd *cobra.Command, args []string) error {
	slog.Debug("got request to start server", "args", args)

	// unmarshal server config with viper
	var srvCfg http.Config
	if err := viper.Unmarshal(&srvCfg); err != nil {
		slog.Error("config unmarshal failed", "err", err)
		return err
	}
	// log version and server port
	slog.Info("Starting iankubetrace",
		slog.String("version", viper.GetString("version")),
		slog.String("revision", viper.GetString("revision")),
		slog.String("port", srvCfg.Port))

	// start http server implementing health checks, etc.
	srv, _ := http.NewServer(&srvCfg, slog.Default())
	httpServer, healthy, ready := srv.ListenAndServe()

	// set up signal handlers to manage graceful shutdown
	stopCh := signals.SetupSignalHandler()
	sd, _ := signals.NewShutdown(srvCfg.ServerShutdownTimeout, slog.Default())
	sd.Graceful(stopCh, httpServer, healthy, ready)
	return nil
}
