package query

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"

	"golang.org/x/xerrors"

	fatihst "github.com/fatih/structs"
	"github.com/getsentry/sentry-go"
	"github.com/gigawattio/metaflector"
	"github.com/jinzhu/copier"
	"github.com/oleiade/reflections"
	"github.com/spf13/cast"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"github.com/eiicon-company/go-core/util"
	"github.com/eiicon-company/go-core/util/repo"
	"github.com/eiicon-company/go-core/util/slices"
	"github.com/eiicon-company/go-core/util/structs"
)

type (
	// BaseModel ...
	BaseModel interface {
		// So Specific below
		// *apimodel.Poster | *apimodel.Organization
		Insert(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) error
		Update(ctx context.Context, exec boil.ContextExecutor, columns boil.Columns) (int64, error)
		Upsert(ctx context.Context, exec boil.ContextExecutor, updateColumns, insertColumns boil.Columns) error
		// more method
		// more method
	}

	// BaseSlice ...
	// TODO: BaseSlice = []*BaseModel
	// BaseSlice = []BaseModel
	BaseSlice interface {
		// constraints.Ordered
		// more method
		// more method
	}

	// BaseGenerator is used as channel result
	BaseGenerator[S BaseSlice] struct {
		Rows S
		Err  error
	}

	// BaseQuery ...
	BaseQuery[M BaseModel, S BaseSlice] interface {
		One(ctx context.Context, exec boil.ContextExecutor) (M, error)
		All(ctx context.Context, exec boil.ContextExecutor) (S, error)
		Count(ctx context.Context, exec boil.ContextExecutor) (int64, error)
		Exists(ctx context.Context, exec boil.ContextExecutor) (bool, error)
		DeleteAll(ctx context.Context, exec boil.ContextExecutor) (int64, error)
		// UpdateAll(ctx context.Context, exec boil.ContextExecutor, cols apimodel.M) (int64, error)
		// more method
		// more method
	}

	// BaseRepo ...
	BaseRepo[M BaseModel, S BaseSlice] interface {
		// Returns Query TODO: remove someday
		Q(mods ...qm.QueryMod) BaseQuery[M, S]

		// Connection
		Conn() *sql.DB
		// Create a record
		Create(ctx context.Context, db boil.ContextExecutor, m M) error
		// Deprecated: Upsert a record. XXX: This is experimental!! Should be worried about that you want to use it.
		Upsert(ctx context.Context, db boil.ContextExecutor, m M) error
		// Update a racord by a model
		Update(ctx context.Context, db boil.ContextExecutor, m M) error
		// Delete deletes all of relevant records
		Delete(ctx context.Context, db boil.ContextExecutor, id int) error
		// Find is retriver for a model
		Find(context.Context, boil.ContextExecutor, int) (M, error)
		// FindPreload is retriver with eager preloading
		FindPreload(context.Context, boil.ContextExecutor, int, ...qm.QueryMod) (M, error)
		// FindBy is retriver with eager preloading
		FindBy(context.Context, boil.ContextExecutor, []qm.QueryMod, ...qm.QueryMod) (M, error)
		// LastBy picks a last row up
		LastBy(context.Context, boil.ContextExecutor, []qm.QueryMod, ...qm.QueryMod) (M, error)
		// FirstBy picks a first row up
		FirstBy(context.Context, boil.ContextExecutor, []qm.QueryMod, ...qm.QueryMod) (M, error)
		// All returns sort ordered records
		All(context.Context, boil.ContextExecutor) (S, error)
		// AllPreload returns sort ordered records
		AllPreload(context.Context, boil.ContextExecutor, ...qm.QueryMod) (S, error)
		// AllGenerator returns rows with go channel which is iterable
		AllGenerator(ctx context.Context, db boil.ContextExecutor, loads ...qm.QueryMod) chan *BaseGenerator[S]
		// ListBy is retriver with eager preloading
		ListBy(context.Context, boil.ContextExecutor, []qm.QueryMod, ...qm.QueryMod) (S, error)
		// ListByIDs is retriver with eager preloading
		ListByIDs(context.Context, boil.ContextExecutor, []int, ...qm.QueryMod) (S, error)
		// ListPagerBy is retriver with eager preloading
		ListPagerBy(context.Context, boil.ContextExecutor, []qm.QueryMod, int, int, ...qm.QueryMod) (S, int, error)
		// Exists a record
		Exists(context.Context, boil.ContextExecutor, int) (bool, error)
	}

	baseRepo[M BaseModel, S BaseSlice] struct {
		env util.Environment
		db  *sql.DB

		query any // TODO: func(mods ...qm.QueryMod) BaseQuery[M, S]
	}
)

var (
	// use this for updation
	updateColumns = boil.Blacklist("updated_at", "created_at")
)

func zeroT[T any]() T {
	return *new(T) // nolint:gocritic
}

func boolColumns[M any](m M) ([]string, error) {
	columns := []string{}
	for _, field := range fatihst.Names(m) {
		name, err := reflections.GetFieldTag(m, field, "boil")
		if err != nil {
			return nil, xerrors.Errorf("%#+v.%s boolColumns must be had boil tag: %w", zeroT[M](), field, err)

		}
		kind, err := reflections.GetFieldKind(m, field)
		if err != nil {
			return nil, xerrors.Errorf("%#+v.%s boolColumns something went wrong: %w", zeroT[M](), field, err)
		}
		if kind != reflect.Bool {
			continue
		}

		columns = append(columns, name)
	}

	return columns, nil
}

// TODO: must be removed this function, use just a r.query
func (r *baseRepo[M, S]) Q(mods ...qm.QueryMod) BaseQuery[M, S] {
	// nolint:gocritic
	switch fn := r.query.(type) {
	case func(mods ...qm.QueryMod) BaseQuery[M, S]:
		return fn(mods...)
	}
	// Imitate a above function behavior by reflection
	rr := reflect.ValueOf(r.query).CallSlice([]reflect.Value{reflect.ValueOf(mods)})
	// nolint:gocritic
	switch q := rr[0].Interface().(type) {
	case BaseQuery[M, S]:
		return q
	}

	panic("baseRepo imitation gives up!")
}

// Conn returns *sql.DB
func (r *baseRepo[M, S]) Conn() *sql.DB {
	return r.db
}

func (r *baseRepo[M, S]) Create(ctx context.Context, db boil.ContextExecutor, m M) error {
	columns, err := boolColumns(m)
	if err != nil {
		return xerrors.Errorf("%#+v Create() failure: %w", zeroT[M](), err)
	}

	in := boil.Infer()

	if len(columns) > 0 {
		in = boil.Greylist(columns...)
	}

	return m.Insert(ctx, db, in)
}

func (r *baseRepo[M, S]) Upsert(ctx context.Context, db boil.ContextExecutor, m M) error {
	columns, err := boolColumns(m)
	if err != nil {
		return xerrors.Errorf("%#+v Upsert() failure: %w", zeroT[M](), err)
	}

	up, in := updateColumns, boil.Infer()

	if len(columns) > 0 {
		in = boil.Greylist(columns...)
	}

	return m.Upsert(ctx, db, up, in)
}

func (r *baseRepo[M, S]) Update(ctx context.Context, db boil.ContextExecutor, m M) error {
	id, err := cast.ToIntE(metaflector.Get(m, "ID")) // TODO: generics
	if err != nil {
		return xerrors.Errorf("%#+v Update() missing ID(int) field: %w", zeroT[M](), err)
	}

	qs, err := r.FindPreload(ctx, db, id)
	if err != nil {
		return err
	}

	if err := structs.OverwriteMerge(qs, m); err != nil {
		return xerrors.Errorf("%#+v Update() merge failed pk=%d: %w", zeroT[M](), id, err)
	}

	if _, err := qs.Update(ctx, db, updateColumns); err != nil {
		return err
	}

	return copier.Copy(m, qs)
}

func (r *baseRepo[M, S]) Delete(ctx context.Context, db boil.ContextExecutor, id int) error {
	panic("need to be implemented")
}

func (r *baseRepo[M, S]) Find(ctx context.Context, db boil.ContextExecutor, id int) (M, error) {
	return r.Q(qm.Where("id = ?", id)).One(ctx, db)
}

func (r *baseRepo[M, S]) FindBy(ctx context.Context, db boil.ContextExecutor, where []qm.QueryMod, loads ...qm.QueryMod) (M, error) {
	mods, err := repo.PreloadBy(where, loads...)
	if err != nil {
		return zeroT[M](), xerrors.Errorf("poster FindBy where=%#+v: %w", where, err)
	}

	return r.Q(mods...).One(ctx, db)
}

func (r *baseRepo[M, S]) FirstBy(ctx context.Context, db boil.ContextExecutor, where []qm.QueryMod, loads ...qm.QueryMod) (M, error) {
	mods := []qm.QueryMod{repo.AscOrder, qm.Limit(1)}
	mods = append(mods, where...)
	return r.FindBy(ctx, db, mods, loads...)
}

func (r *baseRepo[M, S]) LastBy(ctx context.Context, db boil.ContextExecutor, where []qm.QueryMod, loads ...qm.QueryMod) (M, error) {
	mods := []qm.QueryMod{repo.DescOrder, qm.Limit(1)}
	mods = append(mods, where...)
	return r.FindBy(ctx, db, mods, loads...)
}

func (r *baseRepo[M, S]) FindPreload(ctx context.Context, db boil.ContextExecutor, id int, loads ...qm.QueryMod) (M, error) {
	return r.Q(repo.PreloadByID(id, loads...)...).One(ctx, db)
}

func (r *baseRepo[M, S]) All(ctx context.Context, db boil.ContextExecutor) (S, error) {
	span := sentry.StartSpan(ctx, fmt.Sprintf("db.repo.base.%#+v.All", zeroT[S]()))
	ctx = span.Context()
	defer span.Finish()

	return r.Q(repo.DescOrder).All(ctx, db)
}

func (r *baseRepo[M, S]) AllPreload(ctx context.Context, db boil.ContextExecutor, loads ...qm.QueryMod) (S, error) {
	span := sentry.StartSpan(ctx, fmt.Sprintf("db.repo.base.%#+v.AllPreload", zeroT[S]()))
	ctx = span.Context()
	defer span.Finish()

	return r.Q(repo.DescPreloads(loads...)...).All(ctx, db)
}

func (r *baseRepo[M, S]) AllGenerator(ctx context.Context, db boil.ContextExecutor, loads ...qm.QueryMod) chan *BaseGenerator[S] {
	iter := make(chan *BaseGenerator[S])

	go func() {
		defer close(iter)

		nextID, limit := -1, 1000

		for {
			select {
			case <-ctx.Done():
				return
			default:
				mods := []qm.QueryMod{}
				mods = append(mods, qm.Limit(limit))
				mods = append(mods, qm.Where("id > ?", nextID))
				mods = append(mods, qm.OrderBy("id ASC"))

				recs, err := r.ListBy(ctx, db, mods, loads...)
				if err != nil {
					iter <- &BaseGenerator[S]{Err: err}
					return
				}

				iter <- &BaseGenerator[S]{Rows: recs}

				rv := reflect.ValueOf(recs) // TODO: generics
				length := rv.Len()

				if length < limit {
					return
				}

				elm := rv.Index(length - 1).Interface()            // TODO: generics
				id, err := cast.ToIntE(metaflector.Get(elm, "ID")) // TODO: generics
				if err != nil {
					err = xerrors.Errorf("%#+v AllGenerator() missing ID(int) field: %w", zeroT[M](), err)
					iter <- &BaseGenerator[S]{Err: err}
					return
				}

				nextID = id
			}
		}
	}()

	return iter
}

func (r *baseRepo[M, S]) ListBy(ctx context.Context, db boil.ContextExecutor, where []qm.QueryMod, loads ...qm.QueryMod) (S, error) {
	mods, err := repo.DescPreloadBy(where, loads...)
	if err != nil {
		return zeroT[S](), xerrors.Errorf("%#+v ListBy where=%#+v: %w", zeroT[S](), where, err)
	}

	span := sentry.StartSpan(ctx, fmt.Sprintf("db.repo.base.%#+v.ListBy", zeroT[S]()))
	ctx = span.Context()
	defer span.Finish()

	return r.Q(append(mods, repo.DescOrder)...).All(ctx, db)
}

func (r *baseRepo[M, S]) ListByIDs(ctx context.Context, db boil.ContextExecutor, ids []int, loads ...qm.QueryMod) (S, error) {
	if len(ids) == 0 {
		return zeroT[S](), xerrors.Errorf("%#+v ListByIDs() no ids arg found ids=%v: %w", zeroT[S](), ids, sql.ErrNoRows)
	}

	ifaces, err := slices.Interfaces(ids)
	if err != nil {
		return zeroT[S](), xerrors.Errorf("%#+v ListByIDs() ids arg failure ids=%v: %w", zeroT[S](), ids, sql.ErrNoRows)
	}

	// NOTE: The IN-Operator often occurs ambiguous column name error in the JOIN+WHERE Clause although the JOIN+ORDER-BY Clause won't be often.
	mods := []qm.QueryMod{qm.WhereIn("id IN ?", ifaces...)} // NOTE: Watch Out!
	return r.ListBy(ctx, db, mods, loads...)
}

func (r *baseRepo[M, S]) ListPagerBy(ctx context.Context, db boil.ContextExecutor, where []qm.QueryMod, limit, offset int, loads ...qm.QueryMod) (S, int, error) {
	total, err := r.Q(where...).Count(ctx, db)
	if err != nil {
		return zeroT[S](), 0, err
	}

	mods := append([]qm.QueryMod{}, where...)
	mods = append(mods, qm.Limit(limit), qm.Offset(offset))

	records, err := r.ListBy(ctx, db, mods, loads...)
	if err != nil {
		return zeroT[S](), 0, err
	}

	return records, int(total), nil
}

func (r *baseRepo[M, S]) Exists(ctx context.Context, db boil.ContextExecutor, id int) (bool, error) {
	_, err := r.Q(qm.Select("id"), qm.Where("id = ?", id), qm.Limit(1)).One(ctx, db)

	if err != nil && !xerrors.Is(err, sql.ErrNoRows) {
		return false, xerrors.Errorf("%#+v Exists() failure: %w", zeroT[M](), err)
	}

	if xerrors.Is(err, sql.ErrNoRows) {
		return false, nil
	}

	return true, nil
}

// TODO: query any to be func(mods ...qm.QueryMod) BaseQuery[M, S],
func newBaseRepo[M BaseModel, S BaseSlice](env util.Environment, db *sql.DB, query any) BaseRepo[M, S] {
	return &baseRepo[M, S]{
		env:   env,
		db:    db,
		query: query,
	}
}
