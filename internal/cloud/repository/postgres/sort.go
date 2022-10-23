package postgres

import (
	"fmt"

	"github.com/baibikov/tensile-cloud/internal/cloud/types"
)

type Sort struct {
	types.Sort
	def columnsDef
}

func NewSort(ss types.Sort, def columnsDef) Sort {
	return Sort{
		Sort: ss,
		def:  def,
	}
}

func (s Sort) OrderQuery(q string) string {
	if s.Empty() || s.Name == "" {
		return q
	}

	col, ok := s.def[s.Name]
	if !ok {
		return q
	}

	return fmt.Sprintf("select t.* from(%s) t order by t.%s %s", q, col, s.Type)
}

// columnsDef set the mapping of input to column
// example:
// "createdAt": "created_at",
// "parentId":  "parent_id",
type columnsDef map[string]string
