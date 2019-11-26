package repo

import (
	"github.com/eiicon-company/go-core/util/priv"
	"github.com/volatiletech/sqlboiler/queries/qm"
	"golang.org/x/xerrors"
)

const (
	// DBFormat translates as common database datetime format
	DBFormat = "2006-01-02 15:04:05"
)

var (
	// ErrExists a record already exists
	ErrExists = xerrors.New("already exists")
)

var (
	// AscOrder defines ORDER BY query as ascending
	AscOrder = qm.OrderBy("id ASC")
	// DescOrder defines ORDER BY query as descending
	DescOrder = qm.OrderBy("id DESC")
)

// Fuzzy this is so fuzzzy
func Fuzzy(v interface{}, arr []interface{}) bool {
	for _, i := range arr {
		if priv.MustString(i) == priv.MustString(v) {
			return true
		}
	}

	return false
}

// PreloadBy assembles QueryMod with where statements
func PreloadBy(where []qm.QueryMod, loads ...string) ([]qm.QueryMod, error) {
	if len(where) <= 0 {
		return nil, xerrors.New("no queries")
	}

	return append(Preloads(loads...), where...), nil
}

// Preload assembles QueryMod with primary id
func Preload(id int, loads ...string) []qm.QueryMod {
	mods := []qm.QueryMod{qm.Where("id = ?", id)}
	for _, load := range loads {
		mods = append(mods, qm.Load(load))
	}

	return mods
}

// Preloads assembles loads
func Preloads(loads ...string) []qm.QueryMod {
	mods := []qm.QueryMod{qm.OrderBy("id ASC")}
	for _, load := range loads {
		mods = append(mods, qm.Load(load))
	}

	return mods
}
