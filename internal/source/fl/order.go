package fl

import (
	"fl-agent/internal/model"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/tidwall/gjson"
)

func loadOrderPage(url string) (*goquery.Document, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return goquery.NewDocumentFromReader(resp.Body)
}

func ParseOrder(order model.Order) (model.Order, error) {
	doc, err := loadOrderPage(order.URL)
	if err != nil {
		return order, err
	}

	doc.Find("script[type='application/ld+json']").Each(func(_ int, s *goquery.Selection) {
		json := s.Text()

		order.Title = gjson.Get(json, "name").String()
		order.Description = gjson.Get(json, "description").String()
		order.Budget = gjson.Get(json, "offers.price").String()
	})

	return order, nil
}
