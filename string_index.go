package gopandas

//Index transfer between string label and int index
type Index interface {
	Index(string) int
	String(int) string
	Length() int
}

//StringIndex implement index on the string
type StringIndex struct {
	mapping map[string]int
	array   []string
}

const (
	//NotExist String not exist int the StringIndex
	NotExist string = "NotExist"
)

//NewStringIndex is the constructor of StringIndex, duplicate record will be dropped
func NewStringIndex(array []string, needCopy bool) *StringIndex {
	mapping := make(map[string]int)
	idx := 0
	for _, value := range array {
		if mapping[value] > 0 {
			//duplicate record
			needCopy = true
		} else {
			mapping[value] = idx + 1 //Add one for convenient on Index
			idx++
		}
	}
	ret := &StringIndex{
		mapping: mapping,
		array:   array,
	}
	if needCopy {
		ret.array = make([]string, len(mapping))
		for value, idx := range mapping {
			ret.array[idx-1] = value
		}
	}
	return ret
}

//Index get idx, if not exist return -1
func (s *StringIndex) Index(ftr string) int {
	if s == nil {
		return -1
	}
	return s.mapping[ftr] - 1
}

//String get name for idx, if not exist, return ""
func (s *StringIndex) String(idx int) string {
	if s == nil {
		return ""
	}
	return StringArrayGetElse(s.array, idx, "")
}

//Length give the index length
func (s *StringIndex) Length() int {
	if s == nil {
		return 0
	}
	return len(s.array)
}

//Append add new string, if str is duplicate, nothing happens
func (s *StringIndex) Append(str string) {
	if s.Index(str) >= 0 {
		return
	}
	s.array = append(s.array, str)
	s.mapping[str] = len(s.array)
	return
}
