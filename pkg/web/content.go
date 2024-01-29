package web

import (
	"context"
	"fmt"
	"html"
	"log/slog"
	"strings"

	"github.com/elumbantoruan/feed/pkg/feed"
	"github.com/elumbantoruan/feed/pkg/web/storage"
	"go.opentelemetry.io/otel"

	"github.com/chasefleming/elem-go"
	"github.com/chasefleming/elem-go/attrs"
	"github.com/chasefleming/elem-go/styles"
	"github.com/gofiber/fiber/v2"
)

func NewContent(webStorage *storage.WebStorage, logger *slog.Logger) *Handler {
	tracer := otel.Tracer("newsfeed-web")

	return &Handler{
		webStorage: webStorage,
		logger:     logger,
		tracer:     tracer,
	}
}

func (h *Handler) RenderContentRoute(c *fiber.Ctx) error {
	c.Type("html")
	ctx := context.Background()
	ctx, span := h.tracer.Start(ctx, "FeedService.AddSiteFeed")
	defer span.End()

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
		replacer := strings.NewReplacer("position:", "", "absolute;", "", `allowfullscreen=""`, "", "allowfullscreen", "")
		return replacer.Replace(c)
	}

	var container = elem.Div(attrs.Props{attrs.ID: data.Site.Site, attrs.Class: "tabcontent"})

	expandCollapse := func(action string, containerID string) map[string]string {
		script := fmt.Sprintf("showHideRow('%s')", containerID)

		props := map[string]string{
			action: script,
		}
		if action == "onclick" {
			props["class"] = "pointer"
		} else if action == "ondblclick" {
			props["class"] = "hide"
			props["id"] = containerID
		}
		return props
	}

	for i, article := range data.Articles {
		var title, published, desc1, desc2 *elem.Element
		title = elem.P(attrs.Props{attrs.Style: termStyle.ToInline()}, elem.A(attrs.Props{attrs.Href: article.Link, attrs.Target: "_blank"}, elem.Text(clean(article.Title))))

		authors := strings.Join(article.Authors, ", ")
		publishedDateAuthors := fmt.Sprintf("%s - %s", article.Published.String(), authors)
		published = elem.P(attrs.Props{attrs.Style: italicStyle.ToInline()}, elem.Text(publishedDateAuthors))

		if article.Title != article.Description {
			desc1 = elem.P(nil, elem.Text(clean(article.Description)))
		}
		contID := fmt.Sprintf("contS%dC%d", data.Site.ID, i)

		if article.Description != article.Content {
			desc2 = elem.Div(expandCollapse("ondblclick", contID), elem.Text(clean(article.Content)))
		}

		container.Children = append(container.Children, title)
		container.Children = append(container.Children, published)
		if desc1 != nil {
			container.Children = append(container.Children, desc1)
		}
		if desc2 != nil {
			expander := elem.Label(expandCollapse("onclick", contID), elem.Text("   &#8595"))
			published.Children = append(published.Children, expander)
			// container.Children = append(container.Children, expander)
			container.Children = append(container.Children, desc2)
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
		elem.Meta(attrs.Props{attrs.HTTPequiv: "Content-Type", attrs.Content: "text/html", attrs.Charset: "utf-8"}),
		elem.Script(attrs.Props{attrs.Src: "https://cdnjs.cloudflare.com/ajax/libs/jquery/3.3.1/jquery.min.js"}),
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
