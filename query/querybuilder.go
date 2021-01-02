package query

import (
	"fmt"
	"strings"

	"github.com/przebro/databazaar/selector"
)

type couchFormatter struct {
}

func (f *couchFormatter) Format(fld, op, val string) string {
	return fmt.Sprintf(`{"%s":{"%s":%s}}`, fld, op, val)
}
func (f *couchFormatter) FormatArray(op string, val ...string) string {

	result := strings.Join(val, ",")
	return fmt.Sprintf(`{"%s":[%s]}`, op, result)
}

type couchQueryBuilder struct {
	formatter selector.Formatter
}

func NewBuilder() selector.DataSelectorBuilder {

	return &couchQueryBuilder{formatter: &couchFormatter{}}
}

func (qb *couchQueryBuilder) Build(expr selector.Expr) string {

	epxand := expr.Expand(qb.formatter)
	result := fmt.Sprintf(`%s`, epxand)
	return result
}
