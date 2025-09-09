package config

import (
	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
	"github.com/spf13/viper"
)

func NewMidtransClient(config *viper.Viper) *coreapi.Client {
	var client coreapi.Client

	env := midtrans.Sandbox
	if config.GetBool("MIDTRANS_IS_PRODUCTION") {
		env = midtrans.Production
	}

	client.ClientKey = config.GetString("MIDTRANS_CLIENT_KEY")

	client.New(
		config.GetString("MIDTRANS_SERVER_KEY"),
		env,
	)

	return &client
}
