package gopandas

import "time"

//TimeSeries is a one-dimensional ndarray with time series axes.
type TimeSeries struct {
	date   []time.Time
	values []float64
}

//NewTimeSeries give TimeSeries, if values and date doesn't have same length, we drop the tail of long one
func NewTimeSeries(values []float64, date []time.Time, isCopy bool) TimeSeries {
	if len(date) > len(values) {
		date = date[0:len(values)]
	}
	if len(date) < len(values) {
		values = values[0:len(date)]
	}
	ret := TimeSeries{
		date:   date,
		values: values,
	}
	if isCopy {
		ret.date = make([]time.Time, len(date))
		copy(ret.date, date)
		ret.values = make([]float64, len(values))
		copy(ret.values, values)
	}
	return ret
}
