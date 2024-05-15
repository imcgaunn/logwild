package cmd

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"mcgaunn.com/logwild/pkg/cmd/run"
	"mcgaunn.com/logwild/pkg/cmd/version"
	ver "mcgaunn.com/logwild/pkg/version"
)

var (
	debug           bool
	host            string
	port            string
	httpsPort       string
	portMetrics     string
	configPath      string
	certPath        string
	config          string
	otelServiceName string
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logwild",
		Short: "logwild command line",
		Long:  "logwild command line driver program",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// set up global logger with user-specified settings.
			// default to info instead of warning because existing info logs expect to always be printed
			logLevel := slog.LevelInfo
			if debug {
				logLevel = slog.LevelDebug
			}
			slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
				Level: logLevel,
			})))
			slog.Debug("if you see this, log level set to debug",
				"logLevel", logLevel,
				"args", args,
				"cmd", cmd)
			return nil
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			mp := viper.GetString("port-metrics")
			htp := viper.GetString("port")
			otelSvc := viper.GetString("otel-service-name")
			slog.Debug("the end",
				"args", args,
				"port-metrics", mp,
				"port", htp,
				"otel-service-name", otelSvc)
		},
	}

	p := cmd.PersistentFlags()
	p.BoolVar(&debug, "debug", false, "Enable debug logging")
	p.StringVar(&host, "host", "", "Host to which http server should bind")
	p.StringVar(&port, "port", "8888", "Port to which http server should bind")
	p.StringVar(&httpsPort, "https-port", "0", "Port to which secure http server should bind - 0 disables https")
	p.StringVar(&portMetrics, "port-metrics", "0", "Port to which prometheus metrics server should bind - 0 disables metrics")
	p.StringVar(&configPath, "config-path", "", "config dir path")
	p.StringVar(&config, "config", "config.yaml", "config file name within config dir")
	p.StringVar(&otelServiceName, "otel-service-name", "", "service name to report to otel address, disables tracing when not set")

	// bind flags and environment variables
	viper.BindPFlags(p)
	viper.RegisterAlias("configPath", "config-path")
	viper.RegisterAlias("backendUrl", "backend-url")
	hostname, _ := os.Hostname()
	viper.Set("hostname", hostname)
	viper.Set("version", ver.VERSION)
	viper.Set("revision", ver.REVISION)
	viper.SetEnvPrefix("logwild")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	// load config from file
	if _, fileErr := os.Stat(filepath.Join(viper.GetString("config-path"), viper.GetString("config"))); fileErr == nil {
		viper.SetConfigName(strings.Split(viper.GetString("config"), ".")[0])
		viper.AddConfigPath(viper.GetString("config-path"))
		if readErr := viper.ReadInConfig(); readErr != nil {
			fmt.Printf("Error reading config file, %v\n", readErr)
		}
	}

	// register subcommands
	cmd.AddCommand(version.NewVersionCmd())
	cmd.AddCommand(run.NewRunCmd())

	return cmd
}
