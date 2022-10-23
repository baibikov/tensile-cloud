package types

type Sort struct {
	Name string
	Type SortType

	empty bool
}

func (s Sort) Empty() bool {
	return s.empty == true
}

func NewSort(name, typ *string) Sort {
	if name == nil {
		return Sort{
			empty: true,
		}
	}

	return Sort{
		Name: *name,
		Type: NewSortType(typ),
	}
}

type SortType string

func (s SortType) String() string {
	return string(s)
}

const (
	SortDesc SortType = "desc"
	SortAsc  SortType = "asc"
)

func NewSortType(s *string) SortType {
	if s == nil {
		return SortAsc
	}

	switch *s {
	default:
		fallthrough
	case "asc":
		return SortAsc
	case "desc":
		return SortDesc
	}
}
