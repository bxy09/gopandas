package gopandas

import (
	"math"
)

// JoinedSeriesFactory 连接型数据向量工厂
type JoinedSeriesFactory interface {
	LevelIndex(int) (int, int)
	//JoinedIndex(int, int) int
	Index(string) int
	String(int) string
	Length() int

	NewSeries(data []*Series) *JoinedSeries
}

// JoinedSeries 连接型数据向量
type JoinedSeries struct {
	JoinedSeriesFactory
	data []*Series
}

// Get 获取数据
func (s *JoinedSeries) Get(str string) float64 {
	idx := s.JoinedSeriesFactory.Index(str)
	if idx < 0 {
		return math.NaN()
	}
	idx1, idx2 := s.JoinedSeriesFactory.LevelIndex(idx)
	return s.data[idx1].IGet(idx2)
}

// IGet 获取数据
func (s *JoinedSeries) IGet(idx int) float64 {
	if idx < 0 || idx >= s.Length() {
		return math.NaN()
	}
	idx1, idx2 := s.JoinedSeriesFactory.LevelIndex(idx)
	return s.data[idx1].IGet(idx2)
}

// Set 设置数据
func (s *JoinedSeries) Set(str string, value float64) {
	idx := s.JoinedSeriesFactory.Index(str)
	if idx < 0 {
		return
	}
	idx1, idx2 := s.JoinedSeriesFactory.LevelIndex(idx)
	s.data[idx1].ISet(idx2, value)
}

// ISet 设置数据
func (s *JoinedSeries) ISet(idx int, value float64) {
	if idx < 0 || idx >= s.Length() {
		return
	}
	idx1, idx2 := s.JoinedSeriesFactory.LevelIndex(idx)
	s.data[idx1].ISet(idx2, value)
}

// NewJoinedSeriesFactory 从数据构造两级数据的索引
func NewJoinedSeriesFactory(data []Index) JoinedSeriesFactory {
	var names []string
	var levels [][2]int
	for i := range data {
		if data[i] == nil {
			continue
		}
		for j := 0; j < data[i].Length(); j++ {
			names = append(names, data[i].String(j))
			levels = append(levels, [2]int{i, j})
		}
	}
	index := NewStringIndex(names, false)
	return &joinedSeriesFactory{
		StringIndex: index,
		levels:      levels,
	}
}

type joinedSeriesFactory struct {
	*StringIndex
	levels [][2]int
}

// LevelIndex 给出两级数据的索引
func (f *joinedSeriesFactory) LevelIndex(idx int) (int, int) {
	if idx < 0 || idx >= len(f.levels) {
		return -1, -1
	}
	levels := f.levels[idx]
	return levels[0], levels[1]

}

func (f *joinedSeriesFactory) NewSeries(data []*Series) *JoinedSeries {
	return &JoinedSeries{
		JoinedSeriesFactory: f,
		data:                data,
	}
}
