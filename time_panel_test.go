package gopandas_test

import (
	"github.com/bxy09/gopandas"
	"math"
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

func CheckTPIsSame(a, b *gopandas.TimePanel, t *testing.T) {
	if a.Length() != b.Length() {
		t.Error("tp length not same")
		return
	}
	if a.Secondary().Length() != b.Secondary().Length() {
		t.Error("tp second length not same")
		return
	}
	for i := 0; i < a.Secondary().Length(); i++ {
		if a.Secondary().String(i) != b.Secondary().String(i) {
			t.Error("secondary column wrong")
			return
		}
	}
	if a.Thirdly().Length() != b.Thirdly().Length() {
		t.Error("tp third length not same")
		return
	}
	for i := 0; i < a.Thirdly().Length(); i++ {
		if a.Thirdly().String(i) != b.Thirdly().String(i) {
			t.Error("thirdly column wrong")
			return
		}
	}
	for i := 0; i < a.Length(); i++ {
		if !a.IDate(i).Equal(b.IDate(i)) {
			t.Error("Wrong date on %d as %s:%s", i, a.IDate(i), b.IDate(i))
		}
		for j := 0; j < a.Secondary().Length(); j++ {
			seriesA := a.IGet(i).IGet(j)
			seriesB := b.IGet(i).IGet(j)
			for k := 0; k < seriesA.Length(); k++ {
				aValue := seriesA.IGet(k)
				bValue := seriesB.IGet(k)
				if aValue != bValue && !(math.IsNaN(aValue) && math.IsNaN(bValue)) {
					t.Errorf("Wrong on (%s,%s,%s) as %f:%f",
						a.IDate(i).String(), a.Secondary().String(j), a.Thirdly().String(k), aValue, bValue)
				}
			}
		}
	}
}

func DFSame(values [][]float64, df *gopandas.DataFrame, t *testing.T) {
	major := df.Major()
	for i := 0; i < major.Length(); i++ {
		series := df.IGet(i)
		for j := 0; j < series.Length(); j++ {
			if series.IGet(j) != values[i][j] && !(math.IsNaN(series.IGet(j)) && math.IsNaN(values[i][j])) {
				t.Errorf("Wrong on (%d,%d) as %f:%f", i, j, series.IGet(j), values[i][j])
			}
		}
	}
}

func TestImport(t *testing.T) {
	testTarget :=
		`date,key,a,b,c,d
	2016-02-01,000001.SZ,0.19,0.2,0.3,0.4
	2016-01-01,000001.SZ,0.19,0.2,0.3,0.4
	2016-01-01,000003.SZ,0.11,0.2,0.3,0.2`
	panel, err := gopandas.ImportTimePanelFromCSV(testTarget)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("\n" + gopandas.DebugString(panel))
	compile := func(str string) time.Time {
		date, err := time.Parse("2006-01-02", str)
		if err != nil {
			t.Fatal(err)
		}
		return date
	}
	CheckSame(
		[]time.Time{compile("2016-01-01"), compile("2016-02-01")},
		[][][]float64{
			[][]float64{
				[]float64{0.19, 0.2, 0.3, 0.4},
				[]float64{0.11, 0.2, 0.3, 0.2},
			},
			[][]float64{
				[]float64{0.19, 0.2, 0.3, 0.4},
				[]float64{math.NaN(), math.NaN(), math.NaN(), math.NaN()},
			},
		},
		panel,
		t,
	)
}

func TestSecondaryReplace(t *testing.T) {
	aTarget :=
		`date,key,a,b,c,d
2016-01-01,000001.SZ,0.19,0.2,0.3,0.4
2016-02-01,000001.SZ,0.19,0.2,0.3,0.5
2016-01-01,000003.SZ,0.11,0.2,0.3,0.2`
	bTarget :=
		`date,key,a,b,c,f
2016-02-01,000001.SZ,0.19,0.2,0.3,0.4
2016-01-01,000001.SZ,0.19,0.2,0.3,0.4
2016-01-01,000003.SZ,0.11,0.2,0.3,0.2`
	aPanel, _ := gopandas.ImportTimePanelFromCSV(aTarget)
	bPanel, _ := gopandas.ImportTimePanelFromCSV(bTarget)
	aPanel.SecondaryLeftReplace(bPanel)
	oldAPanel, _ := gopandas.ImportTimePanelFromCSV(aTarget)
	t.Log("Test join b")
	CheckTPIsSame(aPanel, oldAPanel, t)
	if t.Failed() {
		t.Fatal()
	}
	cTarget :=
		`date,key,a,b,c,d
2016-01-01,000003.SZ,0.12,0.2,0.4,0.6`
	cPanel, _ := gopandas.ImportTimePanelFromCSV(cTarget)
	aPanel.SecondaryLeftReplace(cPanel)
	dTarget :=
		`date,key,a,b,c,d
2016-01-01,000001.SZ,0.19,0.2,0.3,0.4
2016-02-01,000001.SZ,0.19,0.2,0.3,0.5
2016-01-01,000003.SZ,0.12,0.2,0.4,0.6
2016-02-01,000003.SZ,0.12,0.2,0.4,0.6`
	dPanel, _ := gopandas.ImportTimePanelFromCSV(dTarget)
	t.Log("Test join c")
	CheckTPIsSame(aPanel, dPanel, t)
	if t.Failed() {
		t.Fatal()
	}
	fTarget :=
		`date,key,a,b,c,d
2016-02-01,000001.SZ,0.12,0.2,0.4,0.6`
	fPanel, _ := gopandas.ImportTimePanelFromCSV(fTarget)
	aPanel.SecondaryLeftReplace(fPanel)
	gTarget :=
		`date,key,a,b,c,d
2016-02-01,000001.SZ,0.12,0.2,0.4,0.6
2016-01-01,000003.SZ,0.12,0.2,0.4,0.6
2016-02-01,000003.SZ,0.12,0.2,0.4,0.6`
	gPanel, _ := gopandas.ImportTimePanelFromCSV(gTarget)
	t.Log("Test join f")
	CheckTPIsSame(aPanel, gPanel, t)
	if t.Failed() {
		t.Fatal()
	}
}
