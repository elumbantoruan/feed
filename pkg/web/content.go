package web

import (
	"context"
	"fmt"
	"html"
	"log/slog"
	"strings"

	"github.com/elumbantoruan/feed/pkg/feed"
	"github.com/elumbantoruan/feed/pkg/web/storage"

	"github.com/chasefleming/elem-go"
	"github.com/chasefleming/elem-go/attrs"
	"github.com/chasefleming/elem-go/styles"
	"github.com/gofiber/fiber/v2"
)

func NewContent(webStorage *storage.WebStorage, logger *slog.Logger) *Handler {
	return &Handler{
		webStorage: webStorage,
		logger:     logger,
	}
}

func (h *Handler) RenderContentRoute(c *fiber.Ctx) error {
	c.Type("html")
	ctx := context.Background()
	feeds, err := h.webStorage.GetArticles(ctx)
	if err != nil {
		return err
	}
	return c.SendString(h.renderContent(feeds))
}

func (h *Handler) createContent(data feed.FeedSite[int64]) elem.Node {

	italicStyle := styles.Props{
		styles.FontSize:  "10pt",
		styles.FontStyle: "italic",
	}

	termStyle := styles.Props{
		styles.FontSize:  "12pt",
		styles.FontStyle: "bold",
	}

	clean := func(s string) string {
		c := html.UnescapeString(s)
		return strings.ReplaceAll(c, "â€™", "'")
	}

	var container = elem.Div(attrs.Props{attrs.ID: data.Site.Site, attrs.Class: "tabcontent"})

	for _, article := range data.Articles {
		var title, published, desc1, desc2 *elem.Element
		title = elem.P(attrs.Props{attrs.Style: termStyle.ToInline()}, elem.A(attrs.Props{attrs.Href: article.Link, attrs.Target: "_blank"}, elem.Text(clean(article.Title))))

		authors := strings.Join(article.Authors, ", ")
		publishedDateAuthors := fmt.Sprintf("%s - %s", article.Published.String(), authors)
		published = elem.P(attrs.Props{attrs.Style: italicStyle.ToInline()}, elem.Text(publishedDateAuthors))

		if article.Title != article.Description {
			desc1 = elem.P(nil, elem.Text(clean(article.Description)))
		}
		if article.Description != article.Content {
			desc2 = elem.P(nil, elem.Text(clean(article.Content)))
			_ = desc2
		}

		container.Children = append(container.Children, title)
		container.Children = append(container.Children, published)
		if desc1 != nil {
			container.Children = append(container.Children, desc1)
		}
	}
	return container
}

func (h *Handler) createTab(data feed.FeedSite[int64]) elem.Node {
	onclick := fmt.Sprintf("openContent(event, '%s')", data.Site.Site)
	props := map[string]string{
		attrs.Class: "tablinks",
		"onclick":   onclick,
	}
	if h.firstTab {
		h.firstTab = false
		props[attrs.ID] = "defaultopen"
	}
	button := elem.Button(props, elem.Text(data.Site.Site))

	return button
}

func (h *Handler) renderContent(feeds feed.FeedSites[int64]) string {

	headContent := elem.Head(nil,
		// elem.Script(attrs.Props{attrs.Src: "https://unpkg.com/htmx.org"}),
		elem.TextNode(script),
		elem.TextNode(style),
	)

	h.firstTab = true

	onload := "document.getElementById('defaultopen').click();"
	props := map[string]string{
		"onload": onload,
	}
	bodyContent := elem.Body(
		props,
		elem.H1(nil, elem.Text("News Feed")),
		elem.Div(attrs.Props{attrs.Class: "tab"}, elem.TransformEach(feeds, h.createTab)...),
		elem.Ul(
			nil,
			elem.TransformEach(feeds, h.createContent)...,
		),
	)

	htmlContent := elem.Html(nil, headContent, bodyContent)

	return htmlContent.Render()
}
