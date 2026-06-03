package fl

import (
	"fl-agent/internal/model"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const projectsURL = "https://www.fl.ru/projects/category/programmirovanie/"

func loadProjectsPage() (*goquery.Document, error) {
	req, err := http.NewRequest(http.MethodGet, projectsURL, nil)
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
	doc, err := loadProjectsPage()
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
	return loadProjectsPage()
}
