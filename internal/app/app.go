package app

import (
	"context"
	"os"

	//"github.com/fredmayer/sentry/internal/models"
	"github.com/fredmayer/sentry/internal/models"
	"github.com/fredmayer/sentry/internal/services"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"
)

var configPath string

func Run() {
	// init logger
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

	cmd := &cli.Command{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "config",
				Aliases:     []string{"c"},
				Usage:       "Load configuration from `FILE`",
				Destination: &configPath,
			},
		},
		Name:   "sentry",
		Usage:  "servers service discover",
		Action: action,
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}

func action(ctx context.Context, cmd *cli.Command) error {
	// load configurations from yaml
	config, err := models.LoadConfig(configPath)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}

	// Load service
	s := services.NewServers(config)

	selected := cmd.Args().Get(0)

	// Run scanner
	err = s.Scan(selected)

	return err
}
