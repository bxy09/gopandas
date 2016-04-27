package gopandas_test

import (
	"github.com/bxy09/gopandas"
	"testing"
	"time"
)

func TestTimePanel(t *testing.T) {
	secondIdx := gopandas.NewStringIndex([]string{"a", "b"}, false)
	thirdIdx := gopandas.NewStringIndex([]string{"1a", "2b"}, false)
	now := time.Now()
	dates := []time.Time{now.Add(-time.Hour), now, now, time.Now().Add(time.Hour)}
	tmp := 0.0
	values := make([][][]float64, 4)
	for i := range values {
		values[i] = make([][]float64, 2)
		for j := range values[i] {
			values[i][j] = make([]float64, 2)
			for k := range values[i][j] {
				values[i][j][k] = tmp
				tmp += 1.0
			}

		}
	}
	var tp *gopandas.TimePanel
	checkAdd := func(rank []int) {
		tp = gopandas.NewTimePanel(secondIdx, thirdIdx)
		for _, i := range rank {
			tp.AddMat(dates[i], values[i])
		}
		CheckSame(dates, values, tp, t)
		if t.Failed() {
			t.Fatal("Wrong at ", rank)
		}
	}
	checkAdd([]int{0, 1, 2, 3})
	checkAdd([]int{3, 1, 2, 0})
	checkAdd([]int{1, 2, 3, 0})
	checkAdd([]int{1, 3, 2, 0})

	//Check Get
	checkGet := func(rank int, date time.Time) {
		df, dateOut := tp.Get(date)
		if !dateOut.Equal(dates[rank]) {
			t.Errorf("Wrong date on %s:%s:%s", dateOut, dates[rank], date)
		}
		DFSame(values[rank], df, t)
		if t.Failed() {
			t.Fatalf("Wrong at %d:%s", rank, date)
		}
	}
	checkGet(1, now)
	checkGet(3, now.Add(time.Hour))
	checkGet(3, now.Add(time.Second))
	checkGet(1, now.Add(-time.Second))
	checkGet(0, now.Add(-time.Hour))
	checkGet(0, now.Add(-time.Hour-time.Second))
	df, _ := tp.Get(now.Add(time.Hour + time.Second))
	if df != nil {
		t.Fatal("Should be nil")
	}
	//Check Slice
	checkSlice := func(from time.Time, to time.Time, fromI int, toI int) {
		rtp := tp.Slice(from, to)
		if rtp.Length() != toI-fromI {
			t.Fatalf("Length wrong for %s:%s %d:%d length:%d", from, to, fromI, toI, rtp.Length())
		}
		for i := 0; i < rtp.Length(); i++ {
			if !rtp.IDate(i).Equal(dates[fromI+i]) {
				t.Errorf("Wrong on #%d:%s:%s", i, rtp.IDate(i), dates[fromI+i])
			}
		}
		if t.Failed() {
			t.Fatalf("Wrong for %s:%s %d:%d", from, to, fromI, toI)
		}
	}
	checkSlice(now, now, 1, 1)
	checkSlice(now, now.Add(time.Second), 1, 3)
	checkSlice(now, now.Add(time.Hour), 1, 3)
	checkSlice(now, now.Add(time.Hour+time.Second), 1, 4)
	checkSlice(now.Add(time.Second), now.Add(time.Hour+time.Second), 3, 4)
	checkSlice(now.Add(-time.Second), now.Add(time.Hour), 1, 3)
	checkSlice(now.Add(-time.Hour), now.Add(time.Hour), 0, 3)
	checkSlice(now.Add(-time.Hour-time.Second), now.Add(time.Hour), 0, 3)
	//Check CutHead
	checkCutHead := func(until time.Time, len int) {
		tp.CutHead(until)
		if tp.Length() != len {
			t.Fatalf("Wrong to cut: %s,%d", until.String(), len)
		}
	}
	checkCutHead(now.Add(-time.Hour-time.Second), 4)
	checkCutHead(now.Add(-time.Hour), 4)
	checkCutHead(now.Add(-time.Second), 3)
	checkCutHead(now, 3)
	checkCutHead(now.Add(time.Second), 1)
	checkCutHead(now.Add(time.Hour-time.Second), 1)
	checkCutHead(now.Add(time.Hour), 1)
	checkCutHead(now.Add(time.Hour+time.Second), 0)
	checkCutHead(now.Add(2*time.Hour), 0)
}

func CheckSame(dates []time.Time, values [][][]float64, tp *gopandas.TimePanel, t *testing.T) {
	for i := 0; i < tp.Length(); i++ {
		if !tp.IDate(i).Equal(dates[i]) {
			t.Error("Wrong date on %d as %s:%s", i, tp.IDate(i), dates[i])
		}
		DFSame(values[i], tp.IGet(i), t)
	}
}

func DFSame(values [][]float64, df *gopandas.DataFrame, t *testing.T) {
	major := df.Major()
	for i := 0; i < major.Length(); i++ {
		series := df.IGet(i)
		for j := 0; j < series.Length(); j++ {
			if series.IGet(j) != values[i][j] {
				t.Errorf("Wrong on (%d,%d) as %f:%f", i, j, series.IGet(j), values[i][j])
			}
		}
	}
}
