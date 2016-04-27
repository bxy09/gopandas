package gopandas

//DataFrame is a two-dimensional size-mutable data structure (time*string) with labeled axes (rows and columns).
type DataFrame struct {
	values      [][]float64
	majorIndex  *StringIndex
	secondIndex *StringIndex
}

func NewDataFrame(values [][]float64, majorIndex *StringIndex, secondIndex *StringIndex) *DataFrame {
	return &DataFrame{values, majorIndex, secondIndex}
}

func (d *DataFrame) Get(str string) *Series {
	return d.IGet(d.majorIndex.Index(str))
}

func (d *DataFrame) IGet(i int) *Series {
	if i < 0 || i >= len(d.values) {
		return nil
	}
	return NewSeries(d.values[i], d.secondIndex)
}

func (d *DataFrame) Major() *StringIndex {
	return d.majorIndex
}

func (d *DataFrame) Secondary() *StringIndex {
	return d.secondIndex
}
