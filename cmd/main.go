package main

import (
	"net/http"

	"github.com/kimtruth/LostarkMarket/client"
	"github.com/kimtruth/LostarkMarket/config"
)

func main() {
	cfg := config.NewConfig(config.NewSetting())

	laClient := client.NewLAClient(
		&http.Client{
			Timeout: cfg.Setting().LAHTTPTimeout,
		},
		cfg.Setting().LAHTTPHost,
		cfg.Setting().LAToken,
	)

	laClient.GetRefiningMaterials()
}
