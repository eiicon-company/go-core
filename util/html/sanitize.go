package html

import "github.com/microcosm-cc/bluemonday"

// Sanitize returns a clean HTML to show on public content particularly blog
func Sanitize(htmlMayBeNotSafe string) string {
	p := bluemonday.UGCPolicy()
	p.AllowStyling()
	p.AllowElements("figcaption")
	p.AllowAttrs("width", "height", "max-width", "max-height").Matching(bluemonday.NumberOrPercent).Globally()

	// It allow embeds `Iflamely` code from `Medium Editor`. reference https://iframely.com/
	p.AllowAttrs("style").OnElements("div", "iframe")
	p.AllowAttrs("data-embed-code").OnElements("div")
	p.AllowAttrs("title", "src", "srcdoc", "allow", "allowfullscreen", "scrolling", "frameborder").OnElements("iframe")

	// p.AllowImages()
	p.AllowDataURIImages()

	return p.Sanitize(htmlMayBeNotSafe)
}

// Text extracts raw text contents which contains no html tags
func Text(html string) string {
	p := bluemonday.StrictPolicy()
	return p.Sanitize(html)
}
