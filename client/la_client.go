package client

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/pkg/errors"
)

const (
	userAgent = "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.95 Safari/537.36"
	tokenName = "SUAT"
)

type LAClient struct {
	httpClient *http.Client
	laHost     string
	laToken    string
}

type MarketItem struct {
	Name              string
	YesterdayAvgPrice float64
	RecentPrice       float64
	LowestPrice       float64
}

func NewLAClient(httpClient *http.Client, laHost string, laToken string) *LAClient {
	return &LAClient{
		httpClient: httpClient,
		laHost:     laHost,
		laToken:    laToken,
	}
}

func (c *LAClient) postFormRequest(path string, values url.Values) (*http.Response, error) {
	req, err := http.NewRequest("POST", c.laHost+path, strings.NewReader(values.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Cookie", fmt.Sprintf("%s=%s", tokenName, c.laToken))

	return c.httpClient.Do(req)
}

func (c *LAClient) GetRefiningMaterials() ([]MarketItem, error) {
	params := make(url.Values)
	params.Add("firstCategory", "50000") // 강화재료
	params.Add("secondCategory", "0")    // 전체
	params.Add("tier", "3")              // 티어
	params.Add("isInit", "false")
	params.Add("sortType", "7")

	var totalItems []MarketItem
	pageNo := 1
	for {
		params.Set("pageNo", strconv.Itoa(pageNo))
		resp, err := c.postFormRequest("/Market/List_v2", params)
		if err != nil {
			return nil, err
		}
		items, err := parseMarketItems(resp.Body)
		closeBody(resp)
		if err != nil {
			return nil, err
		}
		if len(items) == 0 {
			break
		}
		totalItems = append(totalItems, items...)
		pageNo += 1
	}

	return totalItems, nil
}

func parseMarketItems(r io.ReadCloser) ([]MarketItem, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}

	var items []MarketItem
	for i := 0; i < 10; i++ {
		row := doc.Find(fmt.Sprintf("#tbodyItemList > tr:nth-child(%d)", i+1)).First()

		itemName := row.Find("td:nth-child(1) > div > span.name").First().Text()
		if len(itemName) == 0 {
			break
		}
		yesterdayAvg, err := toFloat64(row.Find("td:nth-child(2) > div > em").First().Text())
		if err != nil {
			return nil, errors.WithStack(err)
		}
		recentPrice, err := toFloat64(row.Find("td:nth-child(3) > div > em").First().Text())
		if err != nil {
			return nil, errors.WithStack(err)
		}
		lowestPrice, err := toFloat64(row.Find("td:nth-child(4) > div > em").First().Text())
		if err != nil {
			return nil, errors.WithStack(err)
		}

		items = append(items, MarketItem{
			Name:              itemName,
			YesterdayAvgPrice: yesterdayAvg,
			RecentPrice:       recentPrice,
			LowestPrice:       lowestPrice,
		})
	}

	return items, nil
}

func toFloat64(s string) (float64, error) {
	s = strings.Replace(s, ",", "", -1)
	return strconv.ParseFloat(s, 64)
}

func closeBody(resp *http.Response) {
	if _, err := io.Copy(io.Discard, resp.Body); err != nil {
		log.Fatal(err)
	}
	if err := resp.Body.Close(); err != nil {
		log.Fatal(err)
	}
}
