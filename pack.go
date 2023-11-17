package algopack

type PackService interface {
	// BeforeExecute /**
	BeforeExecute(params Params) Params

	// Execute /**
	Execute(params Params, imgExec ImageExec, annotations []Annotation, toVal bool, args ...any) (map[uint64]int, error)

	// GenerateConfigFile /**
	GenerateConfigFile(params Params)
}

type Params struct {
	Labels            []Label
	Relations         map[uint64]uint64
	TrainLabelName    []TrainLabelName
	TrainLabelIdIndex map[uint64]int
	Params            map[string]any
	BasePath          string
	YamlName          string
	Subs              []string
}

type Label struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`
}
