package freelance

import (
	"fl-agent/internal/model"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const tasksURL = "https://freelance.ru/task"

func loadTasksPage() (*goquery.Document, error) {
	req, err := http.NewRequest(http.MethodGet, tasksURL, nil)
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

func loadTaskPage(url string) (*goquery.Document, error) {
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
	doc, err := loadTasksPage()
	if err != nil {
		return nil, err
	}

	var orders []model.Order

	doc.Find("article.task-card").Each(func(_ int, card *goquery.Selection) {
		link := card.Find("a.task-card__title-link").First()

		href, ok := link.Attr("href")
		if !ok {
			return
		}

		title := strings.TrimSpace(link.Text())
		if title == "" {
			title = strings.TrimSpace(link.AttrOr("title", ""))
		}

		if title == "" {
			return
		}

		id := extractTaskID(href)
		if id == "" {
			return
		}

		description := strings.TrimSpace(card.Find(".task-card__desc").First().Text())

		if !strings.HasPrefix(href, "http") {
			href = "https://freelance.ru" + href
		}

		orders = append(orders, model.Order{
			Source:      "freelance",
			ID:          id,
			URL:         href,
			Title:       title,
			Description: description,
		})
	})

	limit := 10
	if len(orders) < limit {
		limit = len(orders)
	}

	return orders[:limit], nil
}

func ParseOrder(order model.Order) (model.Order, error) {
	doc, err := loadTaskPage(order.URL)
	if err != nil {
		return order, err
	}

	description := strings.TrimSpace(doc.Find("meta[name='description']").AttrOr("content", ""))
	if description != "" {
		order.Description = description
	}

	budget := strings.TrimSpace(doc.Find(".tv-meta-item__val--budget").First().Text())
	order.Budget = cleanSpaces(budget)

	return order, nil
}

func extractTaskID(href string) string {
	parts := strings.Split(strings.Trim(href, "/"), "/")
	if len(parts) < 3 {
		return ""
	}

	if parts[0] != "task" || parts[1] != "view" {
		return ""
	}

	return parts[2]
}

func cleanSpaces(s string) string {
	s = strings.ReplaceAll(s, "\u00a0", " ")
	return strings.Join(strings.Fields(s), " ")
}
