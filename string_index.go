package gopandas

//StringIndex implement index on the string
type StringIndex struct {
	mapping map[string]int
	array   []string
}

const (
	//NotExist String not exist int the StringIndex
	NotExist string = "NotExist"
)

//NewStringIndex is the constructor of StringIndex
func NewStringIndex(array []string, needCopy bool) StringIndex {
	mapping := make(map[string]int)
	for idx, value := range array {
		mapping[value] = idx + 1 //Add one for convenient on Index
	}
	ret := StringIndex{
		mapping: mapping,
		array:   array,
	}
	if needCopy {
		ret.array = make([]string, len(array))
		copy(ret.array, array)
	}
	return ret
}

//Index get idx, if not exist return -1
func (s StringIndex) Index(ftr string) int {
	return s.mapping[ftr] - 1
}

//String get name for idx, if not exist, return ""
func (s StringIndex) String(idx int) string {
	return StringArrayGetElse(s.array, idx, "")
}

//Length give the index length
func (s StringIndex) Length() int {
	return len(s.array)
}
