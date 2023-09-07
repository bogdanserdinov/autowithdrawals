package services

import (
	"log"

	kaminarigosdk "github.com/BoostyLabs/kaminari-go-sdk"
	kaminariclient "github.com/BoostyLabs/kaminari-go-sdk/client"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"github.com/urfave/cli/v2"
)

type Config struct {
	MaxAmount       int64  `json:"maxAmount"`
	WithdrawAddress string `json:"withdrawAddress"`
	CronSpecDate    string `json:"cronSpecDate"`
	ApiKey          string `json:"apiKey"`
	ApiUrl          string `json:"apiUrl"`
}

type Client struct {
	cfg            Config
	kaminariClient kaminarigosdk.Interface
}

func New(cfg Config) (*Client, error) {
	kaminariCfg := kaminariclient.Config{
		ApiKey: cfg.ApiKey,
		ApiUrl: cfg.ApiUrl,
	}

	kaminariClient, err := kaminariclient.DefaultClient(&kaminariCfg)
	if err != nil {
		return nil, errors.Wrap(err, "could create kaminari client")
	}

	return &Client{
		cfg:            cfg,
		kaminariClient: kaminariClient,
	}, nil
}

func (client *Client) WithdrawByAmount(ctx *cli.Context) error {
	c := cron.New()
	_, err := c.AddFunc("@hourly", func() {
		if err := client.withdrawByAmount(); err != nil {
			log.Printf("could not initiate auto withdrawing by amount: %v", err)
		}
	})
	if err != nil {
		log.Printf("could not register auto withdrawing callback by amount")
	}

	c.Run()
	return nil
}

// WithdrawByAmount initiates withdraw from Kaminari to user wallet, when balance > amount.
func (client *Client) withdrawByAmount() error {
	balance, err := client.kaminariClient.GetBalance()
	if err != nil {
		return errors.Wrap(err, "could not retrieve user balance")
	}

	if balance.TotalBalance-balance.FrozenAmount >= client.cfg.MaxAmount {
		estimate, err := client.kaminariClient.EstimateIOChainTx(&kaminarigosdk.EstimateOnChainTxRequest{
			BitcoinAddress: client.cfg.WithdrawAddress,
			Amount:         balance.TotalBalance - balance.FrozenAmount,
		})

		err = client.kaminariClient.SendOnChainPayment(&kaminarigosdk.SendOnChainPaymentRequest{
			BitcoinAddress: client.cfg.WithdrawAddress,
			Amount:         balance.TotalBalance - balance.FrozenAmount - estimate.Fee,
			MerchantID:     "",
		})
		if err != nil {
			return errors.Wrap(err, "could not withdraw full balance")
		}
	}

	return nil
}

func (client *Client) WithdrawByDate(ctx *cli.Context) error {
	c := cron.New()
	_, err := c.AddFunc(client.cfg.CronSpecDate, func() {
		if err := client.withdrawByDate(); err != nil {
			log.Printf("could not initiate auto withdrawing by date: %v", err)
		}
	})
	if err != nil {
		log.Printf("could not register auto withdrawing callback by date")
	}

	c.Run()
	return nil
}

// WithdrawByAmount initiates withdraw from Kaminari user wallet to specified address by date.
func (client *Client) withdrawByDate() error {
	balance, err := client.kaminariClient.GetBalance()
	if err != nil {
		return errors.Wrap(err, "could not retrieve user balance")
	}

	estimate, err := client.kaminariClient.EstimateIOChainTx(&kaminarigosdk.EstimateOnChainTxRequest{
		BitcoinAddress: client.cfg.WithdrawAddress,
		Amount:         balance.TotalBalance - balance.FrozenAmount,
	})

	// defines min transfer amount in satoshi.
	const minOnChainTxAmount = 3000

	availableAmountToWithdraw := balance.TotalBalance - balance.FrozenAmount - estimate.Fee
	if availableAmountToWithdraw >= minOnChainTxAmount {
		err = client.kaminariClient.SendOnChainPayment(&kaminarigosdk.SendOnChainPaymentRequest{
			BitcoinAddress: client.cfg.WithdrawAddress,
			Amount:         availableAmountToWithdraw,
			MerchantID:     "",
		})
		if err != nil {
			return errors.Wrap(err, "could not withdraw full balance")
		}
	}

	return nil
}
