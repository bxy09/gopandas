package gopandas

import "math"

//Series is a one-dimensional ndarray with string series axes.
type Series struct {
	*StringIndex
	values []float64
}

//NewSeries construct serires by values and index
func NewSeries(values []float64, idx *StringIndex) *Series {
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
		StringIndex: idx,
		values:      values,
	}
}

//Get get the value by string
func (s *Series) Get(str string) float64 {
	return s.IGet(s.Index(str))
}

//IGet get by int idx
func (s *Series) IGet(i int) float64 {
	if i < 0 || i >= len(s.values) {
		return math.NaN()
	}
	return s.values[i]
}
