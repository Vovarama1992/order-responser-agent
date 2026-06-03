package fl

import (
	"fl-agent/internal/model"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var projectURLs = []string{
	"https://www.fl.ru/projects/category/programmirovanie/",
	"https://www.fl.ru/projects/category/saity/",
	"https://www.fl.ru/projects/category/ai-iskusstvenniy-intellekt/",
	"https://www.fl.ru/projects/category/avtomatizaciya-biznesa/",
}

func loadProjectsPage(url string) (*goquery.Document, error) {
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

func ParseProjects() ([]model.Order, error) {
	seen := make(map[string]bool)
	var orders []model.Order

	for _, url := range projectURLs {
		pageOrders, err := parseProjectsPage(url)
		if err != nil {
			return nil, err
		}

		limit := 5
		if len(pageOrders) < limit {
			limit = len(pageOrders)
		}

		for _, order := range pageOrders[:limit] {
			if seen[order.ID] {
				continue
			}

			seen[order.ID] = true
			orders = append(orders, order)
		}
	}

	return orders, nil
}

func parseProjectsPage(url string) ([]model.Order, error) {
	doc, err := loadProjectsPage(url)
	if err != nil {
		return nil, err
	}

	var orders []model.Order

	doc.Find("a[data-disposable-project-id]").Each(func(_ int, s *goquery.Selection) {
		orderID, ok := s.Attr("data-disposable-project-id")
		if !ok {
			return
		}

		href, ok := s.Attr("href")
		if !ok {
			return
		}

		title := strings.TrimSpace(s.Text())
		if title == "" {
			return
		}

		if !strings.HasPrefix(href, "http") {
			href = "https://www.fl.ru" + href
		}

		orders = append(orders, model.Order{
			Source: "fl",
			ID:     orderID,
			URL:    href,
			Title:  title,
		})
	})

	return orders, nil
}

func LoadForDebug() (*goquery.Document, error) {
	return loadProjectsPage(projectURLs[0])
}
