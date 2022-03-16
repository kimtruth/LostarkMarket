package main

import (
	"fmt"
	"log"
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

	items, err := laClient.GetRefiningMaterials()
	if err != nil {
		log.Fatalf("%+v", err)
	}

	println("아이템명 / 전일 평균 거래가 / 최근 거래가 / 최저가")
	for i, item := range items {
		fmt.Printf("%d %s %.1f %.1f %.1f\n", i+1, item.Name, item.YesterdayAvgPrice, item.RecentPrice, item.LowestPrice)
	}
}
