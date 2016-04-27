package gopandas

import (
	"math"
	"sort"
	"time"
)

//TimePanel represents wide format panel data, stored as 3-dimensional array with sorted date major index
type TimePanel struct {
	values      [][][]float64
	dates       []time.Time
	secondIndex *StringIndex
	thirdIndex  *StringIndex
}

//TimePanelRO Read only TimePanel
type TimePanelRO interface {
	Length() int
	Slice(from, to time.Time) TimePanelRO
	ISlice(from, to int) TimePanelRO
	Get(time.Time) (*DataFrame, time.Time)
	IGet(int) *DataFrame
	IDate(int) time.Time
}

//NewTimePanel new time panel from second index and third index
func NewTimePanel(secondIndex *StringIndex, thirdIndex *StringIndex) *TimePanel {
	return &TimePanel{
		secondIndex: secondIndex,
		thirdIndex:  thirdIndex,
	}
}

//AddMat is try to push or add mat in data
func (p *TimePanel) AddMat(date time.Time, value [][]float64) {
	len2 := p.secondIndex.Length()
	len3 := p.thirdIndex.Length()
	if len(value) != len2 {
		newValue := make([][]float64, len2)
		copy(newValue, value)
		value = newValue
	}
	for i := range value {
		if len(value[i]) != len3 {
			newArray := make([]float64, len3)
			copied := copy(newArray, value[i])
			for j := copied; j < len3; j++ {
				newArray[j] = math.NaN()
			}
			value[i] = newArray
		}
	}
	dateLen := len(p.dates)
	if dateLen == 0 || !date.Before(p.dates[dateLen-1]) {
		p.dates = append(p.dates, date)
		p.values = append(p.values, value)
	} else {
		insertP := sort.Search(dateLen, func(i int) bool {
			return p.dates[i].After(date)
		})
		newDates := make([]time.Time, dateLen+1)
		newValues := make([][][]float64, dateLen+1)
		copy(newDates, p.dates[:insertP])
		copy(newValues, p.values[:insertP])
		newDates[insertP] = date
		newValues[insertP] = value
		copy(newDates[insertP+1:], p.dates[insertP:])
		copy(newValues[insertP+1:], p.values[insertP:])
		p.dates = newDates
		p.values = newValues
	}
}

//Slice get part of TimePanel result >= from < to
func (p *TimePanel) Slice(from, to time.Time) TimePanelRO {
	i := sort.Search(len(p.dates), func(i int) bool {
		return !p.dates[i].Before(from)
	})
	j := sort.Search(len(p.dates), func(i int) bool {
		return !p.dates[i].Before(to)
	})
	return p.ISlice(i, j)
}

//ISlice get port of TimePanel result >=i <j
func (p *TimePanel) ISlice(i, j int) TimePanelRO {
	if i < 0 {
		i = 0
	}
	if j > len(p.dates) {
		j = len(p.dates)
	}
	if i > j {
		i = j
	}
	return &TimePanel{
		values:      p.values[i:j],
		dates:       p.dates[i:j],
		secondIndex: p.secondIndex,
		thirdIndex:  p.thirdIndex,
	}
}

//CutHead cut value head until
func (p *TimePanel) CutHead(until time.Time) {
	i := sort.Search(len(p.dates), func(i int) bool {
		return !p.dates[i].Before(until)
	})
	p.ICutHead(i)
}

//ICutHead cut value head until
func (p *TimePanel) ICutHead(until int) {
	if until < 0 {
		until = 0
	}
	if until > p.Length() {
		until = p.Length()
	}
	p.values = p.values[until:]
	p.dates = p.dates[until:]
}

//Get get the first DataFrame one  big or equal to date
func (p *TimePanel) Get(date time.Time) (*DataFrame, time.Time) {
	i := sort.Search(len(p.dates), func(i int) bool {
		return !p.dates[i].Before(date)
	})
	df := p.IGet(i)
	date = p.IDate(i)
	return df, date
}

func (p *TimePanel) Length() int {
	return len(p.dates)
}

func (p *TimePanel) IGet(i int) *DataFrame {
	if i < 0 || i >= len(p.dates) {
		return nil
	}
	return NewDataFrame(p.values[i], p.secondIndex, p.thirdIndex)
}

func (p *TimePanel) IDate(i int) time.Time {
	if i < 0 || i >= len(p.dates) {
		return AncientTime
	}
	return p.dates[i]
}
