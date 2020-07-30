package html

import (
	"fmt"
	"strings"
	"testing"

	"golang.org/x/xerrors"

	"github.com/PuerkitoBio/goquery"
	"github.com/google/go-cmp/cmp"
	"github.com/yosssi/gohtml"
)

func TestTransformDataURIFunctionNoAttempt(t *testing.T) {
	t.Helper()

	body := `
	<div>A</div>
	<p>
		<img src="https://example.com" class="myclass1"/>
	</p>
	<div class="myclass4">B</div>
	<p>
		<img src="http://example.com" class="myclass2"/>
	</p>
	<div>C</div>
	<div class="myclass5">B</div>
	<p>
		<img src="//example.com" class="myclass3"/>
	</p>
	<div>D</div>
`

	_, err := TransformDataURI(body, func( /*idx int, */ attr string) string {
		return fmt.Sprintf("replaced-%s", attr)
	})
	if !xerrors.Is(err, ErrTransformNoAttempt) {
		t.Fatalf("must be returning ErrTransformNoAttempt error: %+#v", err)
	}
}

func TestTransformDataURIFunction(t *testing.T) {
	t.Helper()

	body := `
	<div>A</div>
	<p>
		<img src="data:image/jpeg;base64,/9j/aaaaaaaaaaaaaaaaaaaaaaa" class="myclass1"/>
	</p>
	<div class="myclass3">B</div>
	<p>
		<img src="data:image/png;base64,/9j/bbbbbbbbbbbbbbbbbbbbbbb" class="myclass2"/>
	</p>
	<div>C</div>
`

	got, err := TransformDataURI(body, func( /*idx int, */ attr string) string {
		if !strings.Contains(attr, "data:image") {
			t.Fatalf("must be contained data:image scheme")
		}

		return fmt.Sprintf("replaced-%s", attr)
	})
	if err != nil {
		t.Fatalf("must be returning non-nil error: %+#v", err)
	}

	expect := `
	<div>A</div>
	<p>
		<img src="replaced-data:image/jpeg;base64,/9j/aaaaaaaaaaaaaaaaaaaaaaa" class="myclass1"/>
	</p>
	<div class="myclass3">B</div>
	<p>
		<img src="replaced-data:image/png;base64,/9j/bbbbbbbbbbbbbbbbbbbbbbb" class="myclass2"/>
	</p>
	<div>C</div>
`

	if diff := cmp.Diff(gohtml.FormatWithLineNo(expect), gohtml.FormatWithLineNo(got)); diff != "" {
		t.Fatalf("return value mismatch (-expect +got):\n%s", diff)
	}

	t.Logf("Return: %+#v", got)
}

func TestTransformDataURIFromTwiceImages(t *testing.T) {
	t.Helper()

	body := `
	<div>A</div>
	<p>
		<img src="data:image/jpeg;base64,/9j/aaaaaaaaaaaaaaaaaaaaaaa" class="myclass1"/>
	</p>
	<div class="myclass3">B</div>
	<p>
		<img src="data:image/png;base64,/9j/bbbbbbbbbbbbbbbbbbbbbbb" class="myclass2"/>
	</p>
	<div>C</div>
`

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		t.Fatalf("must be returning non-nil error: %+#v", err)
	}

	doc.Find(selectorDataURI).Each(func(i int, img *goquery.Selection) {
		img.SetAttr("src", "replaced")
	})

	got, _ := doc.Find("body").Html()
	expect := `
	<div>A</div>
	<p>
		<img src="replaced" class="myclass1"/>
	</p>
	<div class="myclass3">B</div>
	<p>
		<img src="replaced" class="myclass2"/>
	</p>
	<div>C</div>
`

	if diff := cmp.Diff(gohtml.FormatWithLineNo(expect), gohtml.FormatWithLineNo(got)); diff != "" {
		t.Fatalf("return value mismatch (-expect +got):\n%s", diff)
	}

	t.Logf("Return: %+#v", got)
}

// nolint:funlen
func TestTransformDataURIFromEntireContent(t *testing.T) {
	t.Helper()

	body := testIncludeDataURIHTML

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		t.Fatalf("must be returning non-nil error: %+#v", err)
	}

	doc.Find(selectorDataURI).Each(func(i int, img *goquery.Selection) {
		img.SetAttr("src", "replaced")
	})

	got, _ := doc.Find("body").Html()
	expect := `
<p>
  <img src="replaced"/>
</p>
<h2>
  無数にある目標達成方法の中で…
</h2>
<p>
  　目標達成のためのメソッドは、当社のコンテンツ以外にもたくさん見かけます。
</p>
<p>
  <br/>
</p>
<p>
  　私が皆さんにいつも伝えていることとして、目標は考えられる限界まで高く大きく設定を、とお話ししています。
</p>
<p>
  <br/>
</p>
<p>
  　でもこの意見 はどうも、少数派のようですね。
</p>
<p>
  <br/>
</p>
<p>
  　私が見つけた「高すぎる大きすぎる目標は良くない」とする理由を見てみますと、概ね次のような理由です。
</p>
<ol>
  <li>
    高すぎる大きすぎる目標に、すぐに諦める気持ちになる
  </li>
  <li>
    高すぎる大きすぎる目標は、返って恐怖心を呼んでしまう
  </li>
  <li>
    変化を避ける生物の習性が、高すぎる大きすぎる目標を遠ざけてしまう
  </li>
</ol>
<p>
  <br/>
</p>
<p>
  　1. は心理的な視点でしょうか。
</p>
<p>
  　最初からできない言い訳を準備してしまっていますよ、という意味のようです。
</p>
<p>
  　2. は、実際によく聞くセミナー参加者の言葉です。
</p>
<p>
  　3. は脳科学由来でしょうか。　間違ってはいない分、惜しいですね。
</p>
<p>
  <br/>
</p>
<p>
  　これ らに共通しているところは、『意識』の範囲で何とかしよう、としている点です。
</p>
<p>
  <br/>
</p>
<h2>
  高く大きい方が良い理由 １
</h2>
<p>
  　改めて申しておきますが、当社のコンテンツは認知科学という『科学』で裏付けられており 、再現性があると証明されています。
</p>
<p>
  　その根拠を基に、想像を超えた大きな力を持つ『無意識』を味方に付けましょう、という考えです。
</p>
<p>
  <br/>
</p>
<p>
  　手段を講じなければ、無意識はほとんどの場合、意識の邪魔をします。足を引っ張ります。
</p>
<p>
  　変化を嫌う無意識の基本的な働きのため、変化したいという意識の邪魔をしてしまうわけですね。
</p>
<p>
  <br/>
</p>
<p>
  　今では、その無意識を味方に付ける方法が解明されました。
</p>
<p>
  　普通は 邪魔しかしない巨大な力を、味方にすることができるのです。
</p>
<p>
  <br/>
</p>
<p>
  　その条件のひとつとして、『高く大きな目標』であることが重要なのです。
</p>
<p>
  <br/>
</p>
<p>
  　セミナー参加者に
</p>
<p>
  「大きな目標には、恐怖心を感じる。」
</p>
<p>
  といわれました。
</p>
<p>
  <br/>
</p>
<p>
  　そう、その感情が必要なのです。
</p>
<p>
  　恐れて逃げるのではなく、そこから利用するのです。
</p>
<p>
  <br/>
</p>
<p>
  　先の３番目が惜しいのは、脳の特徴までは知っていな がら、『無意識』を利用するまでには至っていない点です。
</p>
<p>
  <br/>
</p>
<h2>
  高く大きい方が良い理由 ２
</h2>
<p>
  　もう一つ、セミナーではまだお話ししたことが少ない、夢や目標は高く大きい方が良い理由をお伝えしましょう。
</p>
<p>
  <br/>
</p>
<p>
  　私たちは、感情によって強く記憶に刷り込まれます。
</p>
<p>
  <br/>
</p>
<p>
  　想い出には、常に強い感情が伴っているのが普通です。
</p>
<p>
  　逆に、感情が伴わなかった出来事は、思い出すのも大変です。
</p>
<p>
  <br/>
</p>
<p>
  　できそうなことたくさんやってきても、心が動かないので記憶になかなか残りません。
</p>
<p>
  <br/>
</p>
<p>
  　あなたには、とてもできそうに無いと思っていたことを達成できて、大いに感動したり喜んだりした経験はありませんか？
</p>
<p>
  <br/>
</p>
<p>
  　そのような経験はしっかりと記憶に刻み込まれるため、その経緯や結果が出たときのことを、いくらでも楽しくお話しできるでしょう。
</p>
<p>
  <br/>
</p>
<p>
  　リストかメモを見返さないと思い出せない多 くの、できそうなことが予想通りできてきただけの人生と、
</p>
<p>
  　数は少なくても、いくらでも喜んで話せる、いや話したくてしょうがない、そんな高く大きな目標を達成できた人生と、
</p>
<p>
  どちらが楽しい人生と言えるでしょ うか。
</p>
<p>
  <br/>
</p>
<p>
  　これが、私が夢や目標を高く大きく持った方が良い理由です。
</p>
<p>
  <br/>
</p>
<p>
  　無意識の使いこなし方はここだけでは伝えきれませんが、様々な事例を通じてこれからも少しずつ、お伝えしていきます 。
</p>
<p>
  <br/>
</p>
<p>
  　ではまた。
</p>
<p>
  <br/>
</p>
<p>
  <br/>
</p>
<p>
  【新型コロナウィルスなんかぶっ飛ばせ！ キャンペーン】
</p>
<p>
  　受け取りましたか？
</p>
<p>
  <a href="https://arkbbw4u.com/info/free-present-textbook/" rel="nofollow">
    セミナーシナリオテキスト無料プレゼント
  </a>
</p>
<p>
  <br/>
</p>
<p>
  　半額以下！ のキャンペーン価格でご提供
</p>
<p>
  <a href="https://arkbbw4u.com/habit-design/" rel="nofollow">
    『脳と心を加速する技術 個別設 計実践』
  </a>
</p>
<p>
  <br/>
</p>
<p>
  <br/>
</p>
<p>
  　メールマガジンでは、読者向けキャンペーンも含めて全文お読みいただけます。
</p>
<p>
  　ご登録は、こちらの登録フォームからお願いします。
</p>
<p>
  <a href="https://www.mshonin.com/form/?id=318063475" rel="nofollow">
    『一生の知恵となる 認知 (脳) 科学コーチング メールマガジン』
  </a>
</p>
<p>
  <br/>
</p>
<p>
  　YouTube も少しずつアップしています。
</p>
<p>
  <a href="https://www.youtube.com/channel/UCbJI7XiXZ9c92wXEws84oyw" rel="nofollow">
    チャンネルはこちら
  </a>
</p>
<p>
  <br/>
</p>
<p>
  　当社セミナーは少人数制です。
</p>
<p>
  　参加者の具体的な悩みや案件にお応えすることで、理解を深めていただいております。
</p>
<p>
  <a href="https://arkbbw4u.com/seminar-info/" rel="nofollow">
    満足度平均 90％のセミナーはこちら
  </a>
</p>
<p>
  <br/>
</p>`

	if diff := cmp.Diff(gohtml.FormatWithLineNo(expect), gohtml.FormatWithLineNo(got)); diff != "" {
		t.Fatalf("return value mismatch (-expect +got):\n%s", diff)
	}

	t.Logf("Return: %+#v", got)
}
