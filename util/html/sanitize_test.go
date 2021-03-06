package html

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestSanitize(t *testing.T) {
	t.Helper()

	expect := `
<div>
1

<a href="https://google.com" title="unko" rel="nofollow">abc</a>
2
<iframe class="ql-video" frameborder="0" allowfullscreen="true" src="https://www.youtube.com/embed/oGsMGdhglb4?showinfo=0"></iframe>
3
<img src="https://example.com" title="unko"/>
4
</div>
`

	got := Sanitize(`
<div>
1
<script>console.log(1)</script>
<a href="https://google.com" title='unko' >abc</a>
2
<iframe class="ql-video" frameborder="0" allowfullscreen="true" src="https://www.youtube.com/embed/oGsMGdhglb4?showinfo=0"></iframe>
3
<img src="https://example.com" title='unko' />
4
</div>
`)

	if diff := cmp.Diff(expect, got); diff != "" {
		t.Fatalf("return value mismatch (-expect +got):\n%s", diff)
	}

	t.Logf("Return: %+#v", got)
}

func TestSanitizeExceptDataURI(t *testing.T) {
	got := Sanitize(testHTMLImageTag)

	if !strings.Contains(got, "img") {
		t.Fatalf("dataURI is gone: %+#v", got)
	}

	// t.Logf("Return: %+#v", got)
}
