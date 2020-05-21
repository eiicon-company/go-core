package repo

import (
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

// PreloadByID assembles QueryMod with primary ID
//
// XXX: There is this code for compatibility
//
func PreloadByID(id int, loads ...qm.QueryMod) []qm.QueryMod {
	mods := []qm.QueryMod{qm.Where("id = ?", id)}
	mods = append(mods, loads...)

	return mods
}

// PreloadBy assembles QueryMod with where statements
//
// XXX: There is this code for compatibility
//
func PreloadBy(where []qm.QueryMod, loads ...qm.QueryMod) ([]qm.QueryMod, error) {
	if len(where) <= 0 {
		return nil, xerrors.New("no queries")
	}

	return append(where, Preloads(loads...)...), nil
}

// Preloads assembles loads
//
// XXX: There is this code for compatibility
//
func Preloads(loads ...qm.QueryMod) []qm.QueryMod {
	mods := []qm.QueryMod{}
	mods = append(mods, loads...)

	return mods
}

// DescPreloadBy assembles QueryMod with where statements
//
// XXX: There is this code for compatibility
//
func DescPreloadBy(where []qm.QueryMod, loads ...qm.QueryMod) ([]qm.QueryMod, error) {
	if len(where) <= 0 {
		return nil, xerrors.New("no queries")
	}

	return append(where, DescPreloads(loads...)...), nil
}

// DescPreloads assembles loads
//
// XXX: There is this code for compatibility
//
func DescPreloads(loads ...qm.QueryMod) []qm.QueryMod {
	mods := []qm.QueryMod{DescOrder}
	mods = append(mods, loads...)

	return mods
}
