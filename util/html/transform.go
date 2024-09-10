package html

import (
	"strings"

	"golang.org/x/xerrors"

	"github.com/PuerkitoBio/goquery"

	"github.com/eiicon-company/go-core/util/logger"
)

const (
	selectorDataURI = `img[src*='data:image/']`
)

var (
	// ErrTransformNoAttempt no attempt to trnasformation
	ErrTransformNoAttempt = xerrors.New("transformation was not attempt")
)

// TransformDataURI gives transform method which can be used for html transformation
// After that returns transformed HTML.
func TransformDataURI /*Transformer*/ (html string, transformer func( /*idx int, */ attr string) string) (string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "", xerrors.Errorf("must be returning non-nil error: %w", err)
	}

	attempt := 0

	doc.Find(selectorDataURI).Each(func(_ int, img *goquery.Selection) {
		attr, ok := img.Attr("src")
		if !ok {
			logger.Warnf("src attribute does not exists. It would be goquery matter")
			return
		}

		img.SetAttr("src", transformer(attr))
		attempt++
	})

	out, err := doc.Find("body").Html()
	if err != nil {
		return "", xerrors.Errorf("generate html failed: %w", err)
	}

	if attempt == 0 {
		return out, ErrTransformNoAttempt
	}

	return out, nil
}
