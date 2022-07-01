package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/andreyAKor/8mar5d-hub/internal/app"
	clientsDevices "github.com/andreyAKor/8mar5d-hub/internal/clients/devices"
	"github.com/andreyAKor/8mar5d-hub/internal/configs"
	"github.com/andreyAKor/8mar5d-hub/internal/http/server"
	"github.com/andreyAKor/8mar5d-hub/internal/logging"
	metricsSensors "github.com/andreyAKor/8mar5d-hub/internal/metrics/sensors"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "8mar5d-hub",
	Short: "8mar5d hub service application",
	Long:  "The 8mar5d hub service is the most simplified service for reading periphery data and present thus on prometheus metrics and on service API.",
	RunE:  run,
}

func init() {
	pf := rootCmd.PersistentFlags()
	pf.StringVar(&cfgFile, "config", "", "config file")

	if err := cobra.MarkFlagRequired(pf, "config"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

//nolint:funlen
func run(cmd *cobra.Command, args []string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Init config
	cfg := &configs.Config{}
	if err := cfg.Init(cfgFile); err != nil {
		return errors.Wrap(err, "init config failed")
	}

	// Init logger
	l := logging.New(cfg.Logging.File, cfg.Logging.Level)
	if err := l.Init(); err != nil {
		return errors.Wrap(err, "init logging failed")
	}

	devicesClient, err := clientsDevices.New(cfg.Clients.Devices.Host)
	if err != nil {
		log.Fatal().Err(err).Msg("can't initialize devices client")
	}

	// Init http-server
	srv, err := server.New(cfg.HTTP.Host, cfg.HTTP.Port, cfg.HTTP.BodyLimit, devicesClient)
	if err != nil {
		log.Fatal().Err(err).Msg("can't initialize http-server")
	}

	// Init metrics
	sensorsMetrics, err := metricsSensors.New(cfg.Metrics.Devices.Interval, devicesClient)

	// Init and run app
	a, err := app.New(srv, sensorsMetrics)
	if err != nil {
		log.Fatal().Err(err).Msg("can't initialize app")
	}
	if err := a.Run(ctx); err != nil {
		log.Fatal().Err(err).Msg("app runnign fail")
	}

	log.Info().Msg("Started")

	// Graceful shutdown
	interruptCh := make(chan os.Signal, 1)
	signal.Notify(interruptCh, os.Interrupt, syscall.SIGTERM)
	<-interruptCh

	log.Info().Msg("Stopping...")

	if err := srv.Close(); err != nil {
		log.Fatal().Err(err).Msg("http-server closing fail")
	}
	if err := a.Close(); err != nil {
		log.Fatal().Err(err).Msg("app closing fail")
	}

	log.Info().Msg("Stopped")

	if err := l.Close(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return nil
}
