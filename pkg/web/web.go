package web

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/elumbantoruan/feed/pkg/feed"
	"github.com/elumbantoruan/feed/pkg/grpc/client"

	"github.com/chasefleming/elem-go"
	"github.com/chasefleming/elem-go/attrs"
	"github.com/chasefleming/elem-go/styles"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

const baseTitle = "News Feed"

var title = baseTitle

type Handler struct {
	GRPCClient client.GRPCFeedClient
	Logger     *slog.Logger
}

func New(grpcClient client.GRPCFeedClient, logger *slog.Logger) *Handler {
	return &Handler{
		GRPCClient: grpcClient,
		Logger:     logger,
	}
}

func (h *Handler) RenderFeedsRoute(c *fiber.Ctx) error {
	c.Type("html")
	return c.SendString(h.renderFeeds(getFeeds()))
}

func (h *Handler) UpdateFeedRoute(c *fiber.Ctx) error {
	newDate := utils.CopyString(c.FormValue("date"))
	if newDate != "" {
		title = fmt.Sprintf("%s: %s", baseTitle, newDate)
	} else {
		title = baseTitle
	}
	return c.Redirect("/")
}

func (h *Handler) createFeedNode(data feed.Feed) elem.Node {
	var dlist = elem.Dl(nil)

	for _, article := range data.Articles {
		var dterm, ddesc1, ddesc2 *elem.Element
		dterm = elem.Dt(nil, elem.H3(nil, elem.A(attrs.Props{attrs.Href: article.Link}, elem.Text(article.Title))))
		if article.Title != article.Description {
			ddesc1 = elem.Dd(nil, elem.Text(article.Description))
		}
		if article.Description != article.Content {
			ddesc2 = elem.Dd(nil, elem.Text(article.Content))
		}
		dlist.Children = append(dlist.Children, dterm)
		if ddesc1 != nil {
			dlist.Children = append(dlist.Children, ddesc1)
		}
		if ddesc2 != nil {
			dlist.Children = append(dlist.Children, ddesc2)
		}
	}
	return elem.Li(nil,
		elem.H2(nil, elem.Text(data.Site)),
		dlist,
	)
}

func (h *Handler) renderFeeds(feeds feed.Feeds) string {

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
					elem.Td(attrs.Props{attrs.Width: "100px"}, elem.Input(attrs.Props{attrs.Type: "text", attrs.Name: "date", attrs.Placeholder: time.Now().String(), attrs.Style: inputDateStyle.ToInline()})),
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

func getFeeds() feed.Feeds {

	feeds := []feed.Feed{
		{
			Site:    "The Verge",
			SiteURL: "https://www.verge.com",
			Articles: []feed.Article{
				{
					Title:       "Samsung’s inexpensive Tab A9 Plus is now on sale in the US",
					Link:        "https://www.theverge.com/2024/1/12/24036471/samsung-galaxy-tab-a9-plus-wifi-5g-price-specs",
					Description: "Samsung’s inexpensive Tab A9 Plus is now on sale in the US",
					Published:   time.Now(),
					Content: `  

					<figure>
					  <img alt="Rendering of Galaxy Tab A9 Plus in silver, graphite, and navy." src="https://cdn.vox-cdn.com/thumbor/jPE313-2plFnWjea3eRexnowC4M=/0x57:1374x973/1310x873/cdn.vox-cdn.com/uploads/chorus_image/image/73050566/Galaxy_Tab_A9_A9_dl1_wider.5.jpg" />
						<figcaption><em>The A9 Plus comes in a 5G version that’s well under $300.</em> | Image: Samsung</figcaption>
					</figure>
				
				  <p id="C9MmQR">Samsung has quietly put <a href="https://go.redirectingat.com?id=66960X1514734&amp;xs=1&amp;url=https%3A%2F%2Fwww.samsung.com%2Fus%2Ftablets%2Fgalaxy-tab-a9-plus%2Fbuy%2Fgalaxy-tab-a9-plus-64gb-navy-wi-fi-sm-x210ndbaxar&amp;referrer=theverge.com&amp;sref=https%3A%2F%2Fwww.theverge.com%2F2024%2F1%2F12%2F24036471%2Fsamsung-galaxy-tab-a9-plus-wifi-5g-price-specs" rel="sponsored nofollow noopener" target="_blank">its budget Galaxy Tab A9 Plus</a> on sale following its <a href="https://news.samsung.com/global/samsung-galaxy-tab-a9-and-galaxy-tab-a9-entertainment-and-productivity-engineered-for-everyone">launch in October last year</a>. It starts at $219 for a Wi-Fi-only version, but unlike most other Android tablets around that price, you can pick up a version with 5G. The A9 Plus with a cellular connection costs $269, and you can take your pick from T-Mobile, Verizon, AT&amp;T, and US Cellular versions on <a href="https://go.redirectingat.com?id=66960X1514734&amp;xs=1&amp;url=http%3A%2F%2FSamsung.com&amp;referrer=theverge.com&amp;sref=https%3A%2F%2Fwww.theverge.com%2F2024%2F1%2F12%2F24036471%2Fsamsung-galaxy-tab-a9-plus-wifi-5g-price-specs" rel="sponsored nofollow noopener" target="_blank">Samsung.com</a>. <a href="https://9to5google.com/2024/01/12/galaxy-tab-a9-plus-us-release/"><em>9to5Google</em> first spotted</a> that the tablet had gone on sale.</p>
				<p id="Jnq9TM">The A9 Plus offers an 11-inch screen with a smooth 90Hz refresh rate — a rare feature at this price and definitely something you won’t find on an entry-level iPad. It has a 5-megapixel front-facing camera, an 8-megapixel rear camera, and comes with a 7,040mAh battery. The Wi-Fi version comes with either 4GB RAM / 64GB storage or 8GB RAM / 128GB of storage; the 5G version is just offered in the 4GB / 64GB configuration.</p>
				<p id="Q9AoQX">For well under $300, that’s an attractive deal on paper, but the A-series tablet misses out on a couple of notable features: S Pen stylus compatibility and an IP rating for water and dust resistance. You’ll have to step up to the Tab S9 series, which starts at $449, if you want either of those things.</p>
				<p id="XToOPT">There’s one more thing the A-series is missing, too: color. The A9 Plus is available in a very straight-laced navy, graphite, or silver. If you’re looking for mint or lavender, well, <a href="https://go.redirectingat.com?id=66960X1514734&amp;xs=1&amp;url=https%3A%2F%2Fwww.samsung.com%2Fus%2Ftablets%2Fgalaxy-tab-s9-fe%2Fbuy%2F&amp;referrer=theverge.com&amp;sref=https%3A%2F%2Fwww.theverge.com%2F2024%2F1%2F12%2F24036471%2Fsamsung-galaxy-tab-a9-plus-wifi-5g-price-specs" rel="sponsored nofollow noopener" target="_blank">you’ll have to pay up</a>.</p>
				
				`,
				},
				{
					Title:     "The tech industry’s layoffs and hiring freezes: all of the news",
					Link:      "https://www.theverge.com/2022/11/14/23458204/meta-twitter-amazon-apple-layoffs-hiring-freezes-latest-tech-industry",
					Published: time.Now(),
					Content: `  

					<figure>
					  <img alt="Image of three cardboard boxes stacked on top of each other, with frowning faces using an upside-down amazon logo for the mouth printed on them." src="https://cdn.vox-cdn.com/thumbor/JZcsP104bGcaNE4fc1gxumGDQkQ=/0x0:2040x1360/1310x873/cdn.vox-cdn.com/uploads/chorus_image/image/71628303/ngarun_181114_1777_amazon_0003.5.jpg" />
						<figcaption><em>Companies have been cutting costs.<span class="ql-cursor"></span></em> | Photo by Natt Garun / The Verge</figcaption>
					</figure>
				
				  <p>Companies across the tech industry have been laying off staff and reducing hiring after explosive growth during the pandemic.</p> <p id="HLjSV6">Over the last couple years, it feels like we’ve heard news of mass layoffs and hiring freezes from tech companies nearly every week, and since the beginning of 2024, there’s been a new wave of layoffs and firings.</p>
				<p id="GZO8uI">In the first few days of January 2024 alone:</p>
				<ul>
				<li id="r8pwZU">Google cut <a href="https://www.theverge.com/2024/1/11/24034124/google-layoffs-engineering-assistant-hardware">around a thousand employees</a>
				</li>
				<li id="Hgv02D">Discord cut <a href="https://www.theverge.com/2024/1/11/24034705/discord-layoffs-17-percent-employees">17 percent of its staff</a>
				</li>
				<li id="W7pGJE">Twitch cut <a href="https://www.theverge.com/2024/1/10/24032187/twitch-layoffs-video-game-industry">a third of its staff</a> (and Amazon fired hundreds <a href="https://www.theverge.com/2024/1/10/24032837/hundreds-of-amazon-prime-video-and-mgm-studios-workers-are-being-laid-off">from Amazon Prime Video and MGM Studios</a>)</li>
				<li id="Fo1kLD">Unity cut <a href="https://www.theverge.com/2024/1/8/24030695/unity-layoff-staff-25-percent">25 percent of its workforce</a>
				</li>
				<li id="3mqYNu">Humane cut <a href="https://www.theverge.com/2024/1/9/24032274/humane-layoffs-ai-pin">four percent of its employees</a>
				</li>
				</ul>
				<p id="ZkgyKS">And all that adds to the <a href="https://www.theverge.com/2023/1/20/23563706/google-layoffs-12000-jobs-cut-sundar-pichai">tens</a> <a href="https://www.theverge.com/2023/1/4/23539737/amazon-layoffs-thousands-17000">of</a> <a href="https://www.theverge.com/2023/3/14/23629726/meta-layoffs-mark-zuckerberg-facebook-tech-recruiting">thousands</a> of tech and gaming layoffs that hit in 2023.</p>
				<p id="i4saf5">Elizabeth Lopatto <a href="https://www.theverge.com/2023/1/26/23571659/tech-layoffs-facebook-google-amazon">spoke to experts</a> in an article published last year to try and answer the question of why so many layoffs are happening right now despite tech companies continuing to register sizable profits. One reason is that “investors have changed how they’re evaluating companies,” even if there’s a lack of evidence that the layoffs can help solve any of the problems they may have.</p>
				<p id="YN4N0t">Here’s all our coverage of the recent outbreak of layoffs from big tech, auto, crypto, gaming, and more.</p>
				
				`,
				},
			},
		},
		{
			Site:    "Wired",
			SiteURL: "https://www.wired.com",
			Articles: []feed.Article{
				{
					Title:       "Regulators Are Finally Catching Up With Big Tech",
					Link:        "https://www.wired.com/story/regulators-are-finally-catching-up-with-big-tech/",
					Description: "The lawless, Wild West era of AI and technology is almost at an end, as data protection authorities use new and existing legislation to get tough.",
					Published:   time.Now(),
					Content:     `The lawless, Wild West era of AI and technology is almost at an end, as data protection authorities use new and existing legislation to get tough.`,
				},
				{
					Title:       "Toyota's Robots Are Learning to Do Housework—By Copying Humans",
					Link:        "https://www.wired.com/story/fast-forward-toyota-robots-learning-housework/",
					Description: "Carmaker Toyota is developing robots capable of learning to do household chores by observing how humans take on the tasks. The project is an example of robotics getting a boost from generative AI.",
					Published:   time.Now(),
					Content:     `Carmaker Toyota is developing robots capable of learning to do household chores by observing how humans take on the tasks. The project is an example of robotics getting a boost from generative AI.`,
				},
			},
		},
	}
	return feeds
}
