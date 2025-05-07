package structs

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/volatiletech/null/v8"
)

type testModel struct {
	ID          int          `bigquery:"id" csv:"id" boil:"id" json:"id" toml:"id" yaml:"id"`
	StoreID     int          `bigquery:"store_id" csv:"store_id" boil:"store_id" json:"store_id" toml:"store_id" yaml:"store_id"`
	ProvstoreID null.Int     `bigquery:"provstore_id" csv:"provstore_id" boil:"provstore_id" json:"provstore_id,omitempty" toml:"provstore_id" yaml:"provstore_id,omitempty"`
	Name        string       `bigquery:"name" csv:"name" boil:"name" json:"name" toml:"name" yaml:"name"`
	DeletedAt   null.Time    `bigquery:"deleted_at" csv:"deleted_at" boil:"deleted_at" json:"deleted_at" toml:"deleted_at" yaml:"deleted_at"`
	IsDeleted   bool         `bigquery:"is_deleted" csv:"is_deleted" boil:"is_deleted" json:"is_deleted" toml:"is_deleted" yaml:"is_deleted"`
	Score       null.Float32 `bigquery:"score" csv:"score" boil:"score" json:"score,omitempty" toml:"score" yaml:"score,omitempty"`
	Reviews     null.Int     `bigquery:"reviews" csv:"reviews" boil:"reviews" json:"reviews,omitempty" toml:"reviews" yaml:"reviews,omitempty"`
	CreatedAt   time.Time    `bigquery:"created_at" csv:"created_at" boil:"created_at" json:"created_at" toml:"created_at" yaml:"created_at"`
	UpdatedAt   time.Time    `bigquery:"updated_at" csv:"updated_at" boil:"updated_at" json:"updated_at" toml:"updated_at" yaml:"updated_at"`
}

func TestOverwriteMerge(t *testing.T) {
	t.Helper()
	t.Parallel()

	now := time.Now()

	dest := &testModel{
		ID:          1,
		StoreID:     1,
		ProvstoreID: null.IntFrom(1),
		Name:        "1",
		DeletedAt:   null.TimeFrom(now),
		IsDeleted:   true,
		Score:       null.Float32From(1),
		Reviews:     null.IntFrom(1),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	src1 := &testModel{
		StoreID:     2,
		ProvstoreID: null.NewInt(2, 2 > 0),
		DeletedAt:   null.TimeFrom(time.Now().Add(1 * time.Hour)),
		IsDeleted:   false,
		UpdatedAt:   time.Now().Add(1 * time.Hour),
	}

	err := OverwriteMerge(dest, src1)
	require.NoError(t, err)

	require.Equal(t, 0, dest.ID)
	require.Equal(t, 2, dest.StoreID)
	require.Equal(t, 2, dest.ProvstoreID.Int)
	require.False(t, dest.IsDeleted)
	require.Equal(t, src1.DeletedAt.Time, dest.DeletedAt.Time)
	require.False(t, dest.UpdatedAt.Equal(now))
}

func TestOverwriteMergeRestParameters(t *testing.T) {
	t.Helper()
	t.Parallel()

	now := time.Now()

	dest := &testModel{
		ID:          1,
		StoreID:     1,
		ProvstoreID: null.IntFrom(1),
		Name:        "1",
		DeletedAt:   null.TimeFrom(now),
		IsDeleted:   true,
		Score:       null.Float32From(1),
		Reviews:     null.IntFrom(1),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	src1 := &testModel{
		StoreID:     2,
		ProvstoreID: null.NewInt(2, 2 > 0),
		DeletedAt:   null.TimeFrom(time.Now().Add(1 * time.Hour)),
		IsDeleted:   false,
	}
	src2 := &testModel{
		Name: "src2",
	}

	err := OverwriteMerge(dest, src1, src2)
	require.NoError(t, err)

	require.Equal(t, 0, dest.ID)
	require.Equal(t, 2, dest.StoreID)
	require.Equal(t, 2, dest.ProvstoreID.Int)
	require.False(t, dest.IsDeleted)
	require.Equal(t, src1.DeletedAt.Time, dest.DeletedAt.Time)
	require.Equal(t, "src2", dest.Name)
}

func TestMerge(t *testing.T) {
	t.Helper()
	t.Parallel()

	now := time.Now()

	dest := &testModel{
		ID:          1,
		StoreID:     1,
		ProvstoreID: null.IntFrom(1),
		Name:        "1",
		DeletedAt:   null.TimeFrom(now),
		IsDeleted:   true,
		Score:       null.Float32From(1),
		Reviews:     null.IntFrom(1),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	src1 := &testModel{
		StoreID:     2,
		ProvstoreID: null.NewInt(2, 2 > 0),
		DeletedAt:   null.TimeFrom(time.Now().Add(1 * time.Hour)),
		IsDeleted:   false,
	}

	err := Merge(dest, src1)
	require.NoError(t, err)

	require.Equal(t, 1, dest.ID)
	require.Equal(t, 1, dest.StoreID)
	require.Equal(t, 1, dest.ProvstoreID.Int)
	require.True(t, dest.IsDeleted)
	require.False(t, dest.DeletedAt.Time.Equal(now))
}

func TestMergeRestParameters(t *testing.T) {
	t.Helper()
	t.Parallel()

	now := time.Now()

	dest := &testModel{
		ID:          1,
		StoreID:     1,
		ProvstoreID: null.IntFrom(1),
		Name:        "1",
		DeletedAt:   null.TimeFrom(now),
		IsDeleted:   true,
		Score:       null.Float32From(1),
		Reviews:     null.IntFrom(1),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	src1 := &testModel{
		StoreID:     2,
		ProvstoreID: null.NewInt(2, 2 > 0),
		DeletedAt:   null.TimeFrom(time.Now().Add(1 * time.Hour)),
		IsDeleted:   false,
	}
	src2 := &testModel{
		Name: "src2",
	}

	err := Merge(dest, src1, src2)
	require.NoError(t, err)

	require.Equal(t, 1, dest.ID)
	require.Equal(t, 1, dest.StoreID)
	require.Equal(t, 1, dest.ProvstoreID.Int)
	require.True(t, dest.IsDeleted)
	require.False(t, dest.DeletedAt.Time.Equal(now))
	require.Equal(t, "1", dest.Name)
}

func TestMergeNil(t *testing.T) {
	t.Helper()
	t.Parallel()

	now := time.Now()

	dest := &testModel{
		ID:          1,
		StoreID:     1,
		ProvstoreID: null.IntFrom(1),
		Name:        "1",
		DeletedAt:   null.TimeFrom(now),
		IsDeleted:   true,
		Score:       null.Float32From(1),
		Reviews:     null.IntFrom(1),
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	src1 := &testModel{
		StoreID:     2,
		ProvstoreID: null.NewInt(2, 2 > 0),
		DeletedAt:   null.TimeFrom(time.Now().Add(1 * time.Hour)),
		IsDeleted:   false,
	}
	src2 := &testModel{
		Name: "src2",
	}

	err := Merge(dest, src1, nil, src2)
	require.NoError(t, err)

	require.Equal(t, 1, dest.ID)
	require.Equal(t, 1, dest.StoreID)
	require.Equal(t, 1, dest.ProvstoreID.Int)
	require.True(t, dest.IsDeleted)
	require.False(t, dest.DeletedAt.Time.Equal(now))
	require.Equal(t, "1", dest.Name)
	require.Equal(t, now, dest.UpdatedAt)
}
