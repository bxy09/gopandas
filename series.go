package gopandas

import (
	"math"
)

//SeriesRO series for read only
type SeriesRO interface {
	Index(string) int
	String(int) string
	Length() int
	Get(string) float64
	IGet(int) float64
}

type SeriesRW interface {
	SeriesRO
	ISet(int, float64)
	Set(string, float64)
}

//Series is a one-dimensional ndarray with string series axes.
type Series struct {
	idx    Index
	values []float64
}

//NewSeries construct serires by values and index
func NewSeries(values []float64, idx Index) *Series {
	length := idx.Length()
	if len(values) != length {
		newValues := make([]float64, length)
		copied := copy(newValues, values)
		for i := copied; i < length; i++ {
			newValues[i] = math.NaN()
		}
		values = newValues
	}
	return &Series{
		idx:    idx,
		values: values,
	}
}

//Index find idx for label
func (s *Series) Index(name string) int {
	return s.idx.Index(name)
}

//String get label for idx
func (s *Series) String(idx int) string {
	return s.idx.String(idx)
}

//Length get len of labels
func (s *Series) Length() int {
	if s == nil {
		return 0
	}
	return s.idx.Length()
}

//Get get the value by string
func (s *Series) Get(str string) float64 {
	if s == nil {
		return math.NaN()
	}
	return s.IGet(s.idx.Index(str))
}

//IGet get by int idx
func (s *Series) IGet(i int) float64 {
	if s == nil {
		return math.NaN()
	}
	if i < 0 || i >= len(s.values) {
		return math.NaN()
	}
	return s.values[i]
}

//ISet set by int idx
func (s *Series) ISet(i int, value float64) {
	if s == nil {
		return
	}
	if i < 0 || i >= len(s.values) {
		return
	}
	s.values[i] = value
}

//Set set by name
func (s *Series) Set(str string, value float64) {
	if s == nil {
		return
	}
	s.ISet(s.Index(str), value)
}
