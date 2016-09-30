package gopandas

//DataFrame is a two-dimensional size-mutable data structure (time*string) with labeled axes (rows and columns).
type DataFrame struct {
	values      [][]float64
	majorIndex  Index
	secondIndex Index
}

func NewDataFrame(values [][]float64, majorIndex Index, secondIndex Index) *DataFrame {
	return &DataFrame{values, majorIndex, secondIndex}
}

func (d *DataFrame) Get(str string) *Series {
	if d == nil {
		return nil
	}
	return d.IGet(d.majorIndex.Index(str))
}

func (d *DataFrame) IGet(i int) *Series {
	if d == nil {
		return nil
	}
	if i >= len(d.values) {
		if i >= d.majorIndex.Length() {
			return nil
		}
		i = len(d.values) - 1
	}
	if i < 0 {
		return nil
	}
	return NewSeries(d.values[i], d.secondIndex)
}

func (d *DataFrame) Major() Index {
	return d.majorIndex
}

func (d *DataFrame) Secondary() Index {
	return d.secondIndex
}
