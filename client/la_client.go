package client

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
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

func (c *LAClient) GetRefiningMaterials() {
	params := make(url.Values)
	params.Add("firstCategory", "50000") // 강화재료
	params.Add("secondCategory", "0")    // 전체
	params.Add("tier", "3")              // 티어
	params.Add("pageNo", "1")
	params.Add("isInit", "false")
	params.Add("sortType", "7")

	resp, err := c.postFormRequest("/Market/List_v2", params)
	if err != nil {
		log.Fatal(err)
	}
	defer closeBody(resp)

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	//Todo: Pagination
	println("아이템명 / 전일 평균 거래가 / 최근 거래가 / 최저가")
	for i := 0; i < 10; i++ {
		row := doc.Find(fmt.Sprintf("#tbodyItemList > tr:nth-child(%d)", i+1)).First()
		if row == nil {
			break
		}
		name := row.Find("td:nth-child(1) > div > span.name").First().Text()
		priceYesterday := row.Find("td:nth-child(2) > div > em").First().Text()
		priceRecent := row.Find("td:nth-child(3) > div > em").First().Text()
		priceLowest := row.Find("td:nth-child(4) > div > em").First().Text()

		println(name, priceYesterday, priceRecent, priceLowest)
	}
}

func closeBody(resp *http.Response) {
	if _, err := io.Copy(io.Discard, resp.Body); err != nil {
		log.Fatal(err)
	}
	if err := resp.Body.Close(); err != nil {
		log.Fatal(err)
	}
}
