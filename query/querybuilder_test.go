package query

import (
	"testing"

	s "github.com/przebro/databazaar/selector"
)

var eqExect = `{"FieldName":{"$eq":"Some value"}}`
var eqExectInt = `{"FieldName":{"$eq":20}}`
var eqExectBool = `{"FieldName":{"$eq":true}}`

var andExpect = `{"$and":[{"Group":{"$eq":"group_1"}},{"Age":{"$eq":20}}]}`
var andOrExpect = `{"$and":[{"Group":{"$eq":"group_1"}},{"$or":[{"Age":{"$eq":20}},{"Age":{"$eq":21}}]}]}`

func TestQueryBuilder(t *testing.T) {

	builder := NewBuilder()
	expr := s.Eq("FieldName", s.String("Some value"))

	result := builder.Build(expr)

	if result != eqExect {
		t.Error("Unexpected value, expected:", eqExect, "actual:", result)
	}

	expr = s.Eq("FieldName", s.Int(20))
	result = builder.Build(expr)

	if result != eqExectInt {
		t.Error("Unexpected value, expected:", eqExectInt, "actual:", result)
	}

	expr = s.Eq("FieldName", s.Bool(true))
	result = builder.Build(expr)

	if result != eqExectBool {
		t.Error("Unexpected value, expected:", eqExectBool, "actual:", result)
	}

}

func TestAndExpr(t *testing.T) {
	builder := NewBuilder()
	expr := s.And(
		s.Eq("Group", s.String("group_1")),
		s.Eq("Age", s.Int(20)),
	)
	result := builder.Build(expr)

	if result != andExpect {
		t.Error("Unexpected value, expected:", andExpect, "actual:", result)
	}

	expr = s.And(
		s.Eq("Group", s.String("group_1")),
		s.Or(s.Eq("Age", s.Int(20)), s.Eq("Age", s.Int(21))),
	)

	result = builder.Build(expr)

	if result != andOrExpect {
		t.Error("Unexpected value, expected:", andExpect, "actual:", result)
	}

}
