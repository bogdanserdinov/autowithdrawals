package main

import (
	"fmt"
	"os"

	"github.com/caarlos0/env/v6"
	"github.com/joho/godotenv"
	"github.com/urfave/cli/v2"

	"autowithdrawals/services"
)

func Config() (*services.Config, error) {
	cfg := new(services.Config)

	err := godotenv.Overload(".env")
	if err != nil {
		return nil, err
	}

	err = env.Parse(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func main() {
	cfg, err := Config()
	if err != nil {
		fmt.Printf("could not get app config: %v", err)
		return
	}

	client, err := services.New(*cfg)
	if err != nil {
		fmt.Printf("could not initialize service: %v", err)
		return
	}

	app := &cli.App{
		Name:      "autowithdrawals",
		Usage:     "BTC auto withdrawals via Kaminari API",
		UsageText: "autowithdrawals [options] <command>",
		Commands: []*cli.Command{
			{
				Name:   "by-amount",
				Action: client.WithdrawByAmount,
			},
			{
				Name:   "by-date",
				Action: client.WithdrawByDate,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Printf("could not start app: %v", err)
	}
}
