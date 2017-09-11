package gopandas_test

import (
	"testing"
	"time"

	"github.com/bxy09/gopandas"
	"github.com/golang/protobuf/proto"
)

func TestProtoTimePanel(t *testing.T) {
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
	tp = gopandas.NewTimePanel(secondIdx, thirdIdx)
	for i := range dates {
		tp.AddMat(dates[i], values[i])
	}
	bytes, err := proto.Marshal(tp)
	if err != nil {
		t.Fatal(err)
	}
	other := new(gopandas.TimePanel)
	err = proto.Unmarshal(bytes, other)
	if err != nil {
		t.Fatal(err)
	}
	CheckSame(dates, values, other, t)
}
