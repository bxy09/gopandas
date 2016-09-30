package gopandas

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/golang/protobuf/proto"
	"math"
	"sort"
	"strconv"
	"strings"
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
	Index(time.Time) int
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
	df := NewDataFrame(value, p.secondIndex, p.thirdIndex)
	if dateLen == 0 || !date.Before(p.dates[dateLen-1]) {
		p.dates = append(p.dates, date)
		p.values = append(p.values, value)
	} else {
		insertP := sort.Search(dateLen, func(i int) bool {
			return p.dates[i].After(date)
		})
		newDates := make([]time.Time, dateLen+1)
		newValues := make([][][]float64, dateLen+1)
		newDfs := make([]DataFrame, dateLen+1)
		copy(newDates, p.dates[:insertP])
		copy(newValues, p.values[:insertP])
		newDates[insertP] = date
		newValues[insertP] = value
		newDfs[insertP] = *df

		copy(newDates[insertP+1:], p.dates[insertP:])
		copy(newValues[insertP+1:], p.values[insertP:])
		p.dates = newDates
		p.values = newValues
	}
}

func (p *TimePanel) IAddDataFrame(date time.Time, df *DataFrame) {
	p.AddMat(date, df.values)
}

// SecondaryLeftReplace 替换Panel Secondary 索引中部分列的值, 时间索引以本地为主, 发生替换时,数值取新数值中该时间之前最接近时间点的数值
// 该 API 目前并不是经济的, 时间开销并没有进行优化
// 如果目标的第三索引与本地不相同,则不进行处理
func (p *TimePanel) SecondaryLeftReplace(target TimePanelRO) {
	if p.Length() == 0 {
		return
	}
	tThird := target.Thirdly()
	if p.thirdIndex.Length() != tThird.Length() {
		return
	}
	for i := 0; i < p.thirdIndex.Length(); i++ {
		if p.thirdIndex.String(i) != tThird.String(i) {
			return
		}
	}
	targetIdx := target.Index(p.dates[0])
	if targetIdx == target.Length() || target.IDate(targetIdx).After(p.dates[0]) {
		targetIdx--
	}
	for i := 0; i < p.Length(); i++ {
		df := p.IGet(i)
		localDate := p.IDate(i)
		for targetIdx+1 < target.Length() && !target.IDate(targetIdx+1).After(localDate) {
			targetIdx++
		}
		for j := 0; j < target.Secondary().Length(); j++ {
			name := target.Secondary().String(j)
			localJ := df.Major().Index(name)
			if localJ < 0 {
				// this column do not exist
				continue
			}
			localSeries := df.IGet(localJ)
			if targetIdx < 0 {
				// invalid target, fill with nan
				for k := range localSeries.values {
					localSeries.values[k] = math.NaN()
				}
			} else {
				targetValue := target.IGet(targetIdx).IGet(j).values
				if len(localSeries.values) < len(targetValue) {
					localSeries.values = make([]float64, len(targetValue))
				}
				copy(localSeries.values, targetValue)
			}
		}
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
	for i := 0; i < until; i++ {
		p.values[i] = nil
		p.dates[i] = time.Time{}
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

//Index give the first idx big or equal to date
func (p *TimePanel) Index(date time.Time) int {
	return sort.Search(len(p.dates), func(i int) bool {
		return !p.dates[i].Before(date)
	})
}

//SetMajor replace the index with input
func (p *TimePanel) SetMajor(dates []time.Time) {
	p.dates = dates
}

//SetSecond replace the index with input
func (p *TimePanel) SetSecond(index Index) {
	p.secondIndex = index
}

//SetThird replace the index with input
func (p *TimePanel) SetThird(index Index) {
	p.thirdIndex = index
}

//Length return length of data
func (p *TimePanel) Length() int {
	if p == nil {
		return 0
	}
	return len(p.dates)
}

//IGet return data on index
func (p *TimePanel) IGet(i int) *DataFrame {
	if p == nil {
		return nil
	}
	if i < 0 {
		i += p.Length()
	}
	if i < 0 || i >= len(p.dates) {
		return nil
	}
	return NewDataFrame(p.values[i], p.secondIndex, p.thirdIndex)
}

//IDate return date on index
func (p *TimePanel) IDate(i int) time.Time {
	if p == nil {
		return AncientTime
	}
	if i < 0 {
		i += p.Length()
	}
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

// ImportTimePanelFromCSV import time panel from csv format string, using the first column for time, second column for
// secondary keys.
func ImportTimePanelFromCSV(str string) (*TimePanel, error) {
	reader := csv.NewReader(strings.NewReader(str))
	reader.TrimLeadingSpace = true
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	headLine := records[0]
	if len(headLine) < 3 {
		return nil, errors.New("Need more columns (dateKey, secondKey, thirdKeys ...)")
	}
	thirdIndex := NewStringIndex(headLine[2:], true)
	secondMap := make(map[string]bool)
	records = records[1:]
	for i := range records {
		secondMap[records[i][1]] = true
	}
	secondKeys := make([]string, 0, len(secondMap))
	for k := range secondMap {
		secondKeys = append(secondKeys, k)
	}
	sort.Strings(secondKeys)
	secondIndex := NewStringIndex(secondKeys, false)
	data := make(map[int64][][]float64)
	for _, record := range records {
		date, err := time.Parse(time.RFC3339, record[0])
		if err != nil {
			date, err = time.Parse("2006-01-02", record[0])
			if err != nil {
				return nil, errors.New("Unknown timeformat")
			}
		}
		unix := date.Unix()
		key := record[1]
		keyIdx := secondIndex.Index(key)
		df := data[unix]
		if len(df) == 0 {
			df = make([][]float64, secondIndex.Length())
		}
		values := make([]float64, thirdIndex.Length())
		for j := range values {
			values[j], err = strconv.ParseFloat(record[j+2], 64)
			if err != nil {
				return nil, err
			}
		}
		df[keyIdx] = values
		data[unix] = df
	}
	dates := make([]int, 0, len(data))
	for k := range data {
		dates = append(dates, int(k))
	}
	sort.Ints(dates)
	panel := NewTimePanel(secondIndex, thirdIndex)
	for _, date := range dates {
		panel.AddMat(time.Unix(int64(date), 0), data[int64(date)])
	}
	return panel, nil
}

func DebugString(panel TimePanelRO) string {
	str := ``
	for i := 0; i < panel.Length(); i++ {
		str += panel.IDate(i).String() + ":\n"
		df := panel.IGet(i)
		for j := 0; j < df.Major().Length(); j++ {
			str += df.Major().String(j) + " ::: "
			series := df.IGet(j)
			for k := 0; k < df.Secondary().Length(); k++ {
				str += fmt.Sprintf("%s:%f, ", series.String(k), series.IGet(k))
			}
			str += "\n"
		}
	}
	return str
}
