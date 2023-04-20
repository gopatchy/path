package path_test

import (
	"reflect"
	"testing"
	"time"

	"cloud.google.com/go/civil"
	"github.com/gopatchy/path"
	"github.com/stretchr/testify/require"
)

type testType1 struct {
	Int     int
	Int64   int64
	UInt    uint
	UInt64  uint64
	Float32 float32
	Float64 float64
	String  string `json:"string2,omitempty"`
	Bool    bool   `json:"bool2"`
	BoolP   *bool

	Ints     []int
	Int64s   []int64
	UInts    []uint
	UInt64s  []uint64
	Float32s []float32
	Float64s []float64
	Strings  []string
	Bools    []bool

	Time  time.Time
	Times []time.Time
	Date  civil.Date
	Dates []civil.Date

	TimeP  *time.Time
	TimesP []*time.Time
}

type testType2 struct {
	Tt1  testType1
	Tt1p *testType1
}

type testType3 struct {
	testType1
}

type testType4 struct {
	Foo *testType5 `json:"foo"`
	testType5
}

type testType5 struct {
	String string `json:"string2,omitempty"`
	Bool   bool   `json:"bool2"`
	UInt   uint   `xtest:"foo,bar,zig"`
}

func TestSet(t *testing.T) {
	t.Parallel()

	tt1 := &testType1{}
	err := path.Set(tt1, "int64", "1234")
	require.NoError(t, err)
	require.Equal(t, int64(1234), tt1.Int64)

	get, err := path.Get(tt1, "int64")
	require.NoError(t, err)
	require.Equal(t, int64(1234), get)

	err = path.Set(tt1, "time", "2022-11-01-08:00")
	require.NoError(t, err)
	require.Equal(t, int64(1667289600), tt1.Time.Unix())

	tt2 := &testType2{}
	err = path.Set(tt2, "tt1p.bool2", "true")
	require.NoError(t, err)
	require.Equal(t, true, tt2.Tt1p.Bool)

	err = path.Set(tt2, "tt1p.string2", "foo")
	require.NoError(t, err)
	require.Equal(t, "foo", tt2.Tt1p.String)

	err = path.Set(tt2, "tt1.boolp", "true")
	require.NoError(t, err)
	require.Equal(t, true, *tt2.Tt1.BoolP)
}

func TestEmbed(t *testing.T) {
	t.Parallel()

	tt3 := &testType3{}
	err := path.Set(tt3, "int", "1234")
	require.NoError(t, err)
	require.Equal(t, 1234, tt3.Int)
}

func TestList(t *testing.T) {
	t.Parallel()

	list := path.List(&testType4{})
	require.Equal(t, []string{
		"UInt",
		"bool2",
		"foo.UInt",
		"foo.bool2",
		"foo.string2",
		"string2",
	}, list)
}

func TestGetFieldType(t *testing.T) {
	t.Parallel()

	typ := reflect.TypeOf(&testType4{})

	typ2 := path.GetFieldType(typ, "bool2")
	require.NotNil(t, typ2)
	require.Equal(t, reflect.TypeOf(true), typ2)

	typ2 = path.GetFieldType(typ, "foo.UInt")
	require.NotNil(t, typ2)
	require.Equal(t, reflect.TypeOf(uint(1)), typ2)
}

func TestFindTagValueType(t *testing.T) {
	t.Parallel()

	typ := reflect.TypeOf(&testType4{})

	p, ok := path.FindTagValueType(typ, "xtest", "foo")
	require.True(t, ok)
	require.Equal(t, "UInt", p)

	p, ok = path.FindTagValueType(typ, "xtest", "zag")
	require.False(t, ok)
	require.Empty(t, p)
}
