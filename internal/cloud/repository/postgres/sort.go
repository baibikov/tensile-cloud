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

	if _, ok := s.def[s.Name]; !ok {
		return q
	}

	return fmt.Sprintf("select t.* from(%s) t order by %s %s", q, s.Name, s.Type)
}

type columnsDef map[string]struct{}

func newColumnsDef(cols ...string) columnsDef {
	cc := make(columnsDef, len(cols))
	for _, c := range cols {
		cc[c] = struct{}{}
	}

	return cc
}
