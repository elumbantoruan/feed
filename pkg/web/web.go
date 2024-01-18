package web

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/elumbantoruan/feed/pkg/feed"
	"github.com/elumbantoruan/feed/pkg/web/storage"

	"github.com/chasefleming/elem-go"
	"github.com/chasefleming/elem-go/attrs"
	"github.com/chasefleming/elem-go/styles"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

var title = "News Feed"

type Handler struct {
	webStorage *storage.WebStorage
	logger     *slog.Logger
	firstTab   bool
}

func New(webStorage *storage.WebStorage, logger *slog.Logger) *Handler {
	return &Handler{
		webStorage: webStorage,
		logger:     logger,
	}
}

func (h *Handler) RenderFeedsRoute(c *fiber.Ctx) error {
	c.Type("html")
	ctx := context.Background()
	feeds, err := h.webStorage.GetArticles(ctx)
	if err != nil {
		return err
	}
	return c.SendString(h.renderFeeds(feeds))
}

func (h *Handler) UpdateFeedRoute(c *fiber.Ctx) error {
	newDate := utils.CopyString(c.FormValue("date"))
	if newDate != "" {
		title = fmt.Sprintf("%s: %s", title, newDate)
	}
	return c.Redirect("/")
}

func (h *Handler) createFeedNode(data feed.FeedSite[int64]) elem.Node {
	var dlist = elem.Dl(nil)

	italicStyle := styles.Props{
		styles.FontSize:  "10pt",
		styles.FontStyle: "italic",
	}

	termStyle := styles.Props{
		styles.FontSize:  "12pt",
		styles.FontStyle: "bold",
	}

	for _, article := range data.Articles {
		var dterm, dpublished, ddesc1, ddesc2 *elem.Element
		dterm = elem.Dt(attrs.Props{attrs.Style: termStyle.ToInline()}, elem.A(attrs.Props{attrs.Href: article.Link, attrs.Target: "_blank"}, elem.Text(article.Title)))

		authors := strings.Join(article.Authors, ", ")
		publishedDateAuthors := fmt.Sprintf("%s - %s", article.Published.String(), authors)
		dpublished = elem.Dd(attrs.Props{attrs.Style: italicStyle.ToInline()}, elem.Text(publishedDateAuthors))

		if article.Title != article.Description {
			ddesc1 = elem.Dd(nil, elem.Text(article.Description))
		}
		if article.Description != article.Content {
			ddesc2 = elem.Dd(nil, elem.Text(article.Content))
			_ = ddesc2
		}
		dlist.Children = append(dlist.Children, dterm)
		dlist.Children = append(dlist.Children, dpublished)
		if ddesc1 != nil {
			dlist.Children = append(dlist.Children, ddesc1)
		}
	}
	return elem.Li(nil,
		elem.H2(nil, elem.Text(data.Site.Site)),
		dlist,
	)
}

func (h *Handler) renderFeeds(feeds feed.FeedSites[int64]) string {

	inputDateStyle := styles.Props{
		styles.Width:           "200px",
		styles.Padding:         "10px",
		styles.MarginBottom:    "10px",
		styles.Border:          "1px solid #ccc",
		styles.BorderRadius:    "4px",
		styles.BackgroundColor: "#f9f9f9",
	}

	buttonDateStyle := styles.Props{
		styles.BackgroundColor: "#007BFF",
		styles.Padding:         "10px",
		styles.MarginBottom:    "10px",
		styles.Color:           "white",
		styles.BorderStyle:     "none",
		styles.BorderRadius:    "1px",
		styles.Cursor:          "pointer",
		styles.Width:           "100px",
		styles.FontSize:        "14px",
	}

	listContainerStyle := styles.Props{
		styles.ListStyleType: "none",
		styles.Padding:       "0",
		styles.Width:         "100%",
	}

	centerContainerStyle := styles.Props{
		styles.Margin:          "40px auto",
		styles.Padding:         "20px",
		styles.Border:          "1px solid #ccc",
		styles.BoxShadow:       "0px 0px 10px rgba(0,0,0,0.1)",
		styles.BackgroundColor: "#f9f9f9",
		styles.FontFamily:      "verdana",
	}

	headContent := elem.Head(nil,
		elem.Script(attrs.Props{attrs.Src: "https://unpkg.com/htmx.org"}),
	)

	bodyContent := elem.Div(
		attrs.Props{attrs.Style: centerContainerStyle.ToInline()},
		elem.H1(nil, elem.Text(title)),
		elem.Form(
			attrs.Props{attrs.Method: "post", attrs.Action: "/update"},
			elem.Table(nil,
				elem.Tr(nil,
					elem.Td(attrs.Props{attrs.Width: "100px"},
						elem.Input(
							attrs.Props{
								attrs.Type:        "text",
								attrs.Name:        "date",
								attrs.Placeholder: time.Now().String(),
								attrs.Style:       inputDateStyle.ToInline(),
							},
						),
					),
					elem.Td(attrs.Props{attrs.Width: "200px"},
						elem.Button(
							attrs.Props{
								attrs.Type:  "submit",
								attrs.Style: buttonDateStyle.ToInline(),
							},
							elem.Text("Update"),
						),
					),
				),
			),
		),
		elem.Ul(
			attrs.Props{attrs.Style: listContainerStyle.ToInline()},
			elem.TransformEach(feeds, h.createFeedNode)...,
		),
	)

	htmlContent := elem.Html(nil, headContent, bodyContent)

	return htmlContent.Render()
}
