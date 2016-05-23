package gopandas

import (
	"github.com/golang/protobuf/proto"
	"math"
	"sort"
	"time"
)

//TimePanel represents wide format panel data, stored as 3-dimensional array with sorted date major index
type TimePanel struct {
	values      [][][]float64
	dates       []time.Time
	secondIndex Index
	thirdIndex  Index
}

//TimePanelRO Read only TimePanel
type TimePanelRO interface {
	Length() int
	Slice(from, to time.Time) TimePanelRO
	ISlice(from, to int) TimePanelRO
	Get(time.Time) (*DataFrame, time.Time)
	IGet(int) *DataFrame
	IDate(int) time.Time
	Secondary() Index
	Thirdly() Index
	ToProtoBuf() ([]byte, error)
}

//NewTimePanel new time panel from second index and third index
func NewTimePanel(secondIndex Index, thirdIndex Index) *TimePanel {
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

//Length return length of data
func (p *TimePanel) Length() int {
	return len(p.dates)
}

//IGet return data on index
func (p *TimePanel) IGet(i int) *DataFrame {
	if i < 0 || i >= len(p.dates) {
		return nil
	}
	return NewDataFrame(p.values[i], p.secondIndex, p.thirdIndex)
}

//IDate return date on index
func (p *TimePanel) IDate(i int) time.Time {
	if i < 0 || i >= len(p.dates) {
		return AncientTime
	}
	return p.dates[i]
}

//Secondary return second index
func (p *TimePanel) Secondary() Index {
	return p.secondIndex
}

//Thirdly return third index
func (p *TimePanel) Thirdly() Index {
	return p.thirdIndex
}

//ToProtoBuf transfer to ProtoBuf
func (p *TimePanel) ToProtoBuf() ([]byte, error) {
	data := make([]float64, p.Length()*p.secondIndex.Length()*p.thirdIndex.Length())
	dates := make([]uint64, p.Length())
	secondary := make([]string, p.secondIndex.Length())
	thirdly := make([]string, p.thirdIndex.Length())

	idx := 0
	thirdLen := p.thirdIndex.Length()
	for i := range p.values {
		for j := range p.values[i] {
			copy(data[idx:], p.values[i][j])
			idx += thirdLen
		}
	}
	for i := range dates {
		dates[i] = uint64(p.dates[i].UnixNano())
	}
	for i := range secondary {
		secondary[i] = p.secondIndex.String(i)
	}
	for i := range thirdly {
		thirdly[i] = p.thirdIndex.String(i)
	}
	return proto.Marshal(&FlyTimePanel{
		Data:      data,
		Dates:     dates,
		Secondary: secondary,
		Thirdly:   thirdly,
	})
}

//FromProtoBuf transfer from ProtoBuf
func (p *TimePanel) FromProtoBuf(bytes []byte) error {
	fp := &FlyTimePanel{}
	err := proto.Unmarshal(bytes, fp)
	if err != nil {
		return err
	}
	p.dates = make([]time.Time, len(fp.Dates))
	for i := range fp.Dates {
		p.dates[i] = time.Unix(0, int64(fp.Dates[i]))
	}
	p.secondIndex = NewStringIndex(fp.Secondary, true)
	p.thirdIndex = NewStringIndex(fp.Thirdly, true)
	p.values = make([][][]float64, len(p.dates))
	idx := 0
	for i := range p.dates {
		p.values[i] = make([][]float64, p.secondIndex.Length())
		for j := range p.values[i] {
			p.values[i][j] = make([]float64, p.thirdIndex.Length())
			copied := 0
			if idx < len(fp.Data) {
				copied = copy(p.values[i][j], fp.Data[idx:])
				idx += copied
			}
			for ; copied < p.thirdIndex.Length(); copied++ {
				p.values[i][j][copied] = math.NaN()
			}
		}
	}
	return nil
}
