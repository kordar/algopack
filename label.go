package algopack

type TrainLabelName struct {
	LabelId uint64
	Index   int
	Name    string
}

func GenLabelNameIndex(names []TrainLabelName) map[uint64]int {
	data := make(map[uint64]int)
	for i, labelName := range names {
		labelId := labelName.LabelId
		data[labelId] = i
	}
	return data
}

func MergeLabelId[T comparable](list []T, t []T) []T {
	mm := map[T]bool{}
	for _, v1 := range list {
		mm[v1] = true
	}
	for _, v2 := range t {
		mm[v2] = true
	}
	return MapKeys[T, bool](mm)
}

func GetRealLabelId(labelId uint64, relation map[uint64]uint64) uint64 {
	if relation != nil && relation[labelId] > 0 {
		return relation[labelId]
	}
	return labelId
}
