package search

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/olivere/elastic"

	"github.com/eiicon-company/go-core/util"
)

type (
	// Command defines command behavior
	Command interface {
		Search(ctx context.Context, search *elastic.SearchService) (*Result, error)
		Bulk(ctx context.Context, bulk *elastic.BulkService) (*Result, error)
		PostDocument(ctx context.Context, name string, id int, doc string) (*Result, error)
		DeleteDocument(ctx context.Context, name string, id int) (*Result, error)
		UpdateByScript(ctx context.Context, name string, id int, script string, params map[string]interface{}) (*Result, error)
		UpsertByScript(ctx context.Context, name string, id int, script string, params, upsert map[string]interface{}) (*Result, error)
		CreateIndex(ctx context.Context, name string, index string) (*Result, error)
		DeleteIndex(ctx context.Context, name string) (*Result, error)
		Aliases(ctx context.Context, name string) (*Result, error)
		PutAlias(ctx context.Context, name, alias string) (*Result, error)
		UpdateAliases(ctx context.Context, name, old, new string) (*Result, error)
	}

	// command defines interfaces as elasticsearch api.
	command struct {
		Env      util.Environment
		ESClient *elastic.Client
	}

	// Result has common to return a value
	Result struct {
		Res interface{} // ES Result Buffer
		Err error
	}
)

// Indices returns values which matches alias name
func (cr *Result) Indices(alias string) []string {
	switch value := cr.Res.(type) {
	case *elastic.AliasesResult:
		return value.IndicesByAlias(alias)
	default:
		return []string{}
	}
}

// JSON returns value as JSON
func (cr *Result) JSON() []byte {
	bytes, _ := json.Marshal(cr.Res)
	return bytes
}

// Values returns significant values which was choosen along with any es result
func (cr *Result) Values() interface{} {
	switch value := cr.Res.(type) {
	default:
		return value
	case *elastic.AliasesResult:
		return value.Indices
	case *elastic.IndicesCreateResult:
		return cr.JSON()
	case *Result:
		return value.Values()
	}
}

// MakeIndexName returns name with timestamp suffix
func MakeIndexName(name string) string {
	return fmt.Sprintf("%s_%d", name, time.Now().UnixNano())
}

// RestoreIndexName returns remove timestamp suffix
func RestoreIndexName(name string) string {
	return strings.Split(name, "_")[0]
}

func (c *command) do(ctx context.Context, fn func(chan *Result)) (*Result, error) {
	rch := make(chan *Result, 1)

	go fn(rch)

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case cr := <-rch:
		return cr, cr.Err
	}
}

func (c *command) Search(ctx context.Context, search *elastic.SearchService) (*Result, error) {
	fn := func(rch chan *Result) {
		res, err := search.Pretty(c.Env.IsDebug()).Do(ctx)
		rch <- &Result{Res: res, Err: err}
	}

	return c.do(ctx, fn)
}

func (c *command) Bulk(ctx context.Context, bulk *elastic.BulkService) (*Result, error) {
	fn := func(rch chan *Result) {
		res, err := bulk.Do(ctx)
		rch <- &Result{Res: res, Err: err}
	}

	return c.do(ctx, fn)
}

func (c *command) PostDocument(ctx context.Context, name string, id int, doc string) (*Result, error) {
	fn := func(rch chan *Result) {
		res, err := c.ESClient.Index().
			Pretty(c.Env.IsDebug()).
			Index(name).Type(name).Id(strconv.Itoa(id)).BodyString(doc).Do(ctx)
		rch <- &Result{Res: res, Err: err}
	}

	return c.do(ctx, fn)
}

func (c *command) UpdateByScript(ctx context.Context, name string, id int, script string, params map[string]interface{}) (*Result, error) {
	fn := func(rch chan *Result) {
		script := elastic.NewScript(script).Params(params).Lang("painless")

		res, err := c.ESClient.Update().
			Pretty(c.Env.IsDebug()).Index(name).Type(name).Id(strconv.Itoa(id)).
			Script(script).Do(ctx)
		rch <- &Result{Res: res, Err: err}
	}

	return c.do(ctx, fn)
}

func (c *command) UpsertByScript(ctx context.Context, name string, id int, script string, params, upsert map[string]interface{}) (*Result, error) {
	fn := func(rch chan *Result) {
		script := elastic.NewScript(script).Params(params).Lang("painless")

		res, err := c.ESClient.Update().
			Pretty(c.Env.IsDebug()).Index(name).Type(name).Id(strconv.Itoa(id)).
			Script(script).ScriptedUpsert(true).Upsert(upsert).Do(ctx)
		rch <- &Result{Res: res, Err: err}
	}

	return c.do(ctx, fn)
}

func (c *command) DeleteDocument(ctx context.Context, name string, id int) (*Result, error) {
	fn := func(rch chan *Result) {
		res, err := c.ESClient.Delete().
			Pretty(c.Env.IsDebug()).
			Index(name).Type(name).Id(strconv.Itoa(id)).Do(ctx)
		rch <- &Result{Res: res, Err: err}
	}

	return c.do(ctx, fn)
}

func (c *command) CreateIndex(ctx context.Context, name string, index string) (*Result, error) {
	fn := func(rch chan *Result) {
		res, err := c.ESClient.CreateIndex(name).
			Pretty(c.Env.IsDebug()).Body(index).Do(ctx)
		rch <- &Result{Res: res, Err: err}
	}

	return c.do(ctx, fn)
}

func (c *command) DeleteIndex(ctx context.Context, name string) (*Result, error) {
	fn := func(rch chan *Result) {
		res, err := c.ESClient.DeleteIndex(name).
			Pretty(c.Env.IsDebug()).Do(ctx)
		rch <- &Result{Res: res, Err: err}
	}

	return c.do(ctx, fn)
}

func (c *command) Aliases(ctx context.Context, name string) (*Result, error) {
	fn := func(rch chan *Result) {
		res, err := c.ESClient.Aliases().
			Pretty(c.Env.IsDebug()).Index(name).Do(ctx)
		rch <- &Result{Res: res, Err: err}
	}

	return c.do(ctx, fn)
}

func (c *command) PutAlias(ctx context.Context, name, alias string) (*Result, error) {
	fn := func(rch chan *Result) {
		res, err := c.ESClient.Alias().
			Pretty(c.Env.IsDebug()).Add(name, alias).Do(ctx)
		rch <- &Result{Res: res, Err: err}
	}

	return c.do(ctx, fn)
}

func (c *command) UpdateAliases(ctx context.Context, name, old, new string) (*Result, error) {
	fn := func(rch chan *Result) {
		res, err := c.ESClient.Alias().
			Pretty(c.Env.IsDebug()).
			Action(elastic.NewAliasRemoveAction(name).Index(old)).
			Action(elastic.NewAliasAddAction(name).Index(new)).
			Do(ctx)
		rch <- &Result{Res: res, Err: err}
	}

	return c.do(ctx, fn)
}

func newCommand(env util.Environment, client *elastic.Client) Command {
	r := &command{
		Env:      env,
		ESClient: client,
	}

	return r
}
