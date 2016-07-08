package data

import (
	"testing"
)

type NextFieldTest struct {
	FieldStr   string
	NextFields []NextField
}

type NextField struct {
	name  string
	array int
}

var nextFieldTests = []NextFieldTest{
	{
		"spec.template.spec.containers(0)",
		[]NextField{
			{name: "spec", array: -1}, {name: "template", array: -1}, {name: "spec", array: -1},
			{name: "containers", array: -1}, {array: 0},
		},
	},
	{
		"(0)(1)(2)(3)(4)(5)ape",
		[]NextField{
			{array: 0}, {array: 1}, {array: 2}, {array: 3}, {array: 4}, {array: 5}, {name: "ape", array: -1},
		},
	},
	{
		"spec(0)helo.eat.cheese",
		[]NextField{
			{name: "spec", array: -1}, {array: 0}, {name: "helo", array: -1}, {name: "eat", array: -1},
			{name: "cheese", array: -1},
		},
	},
	// invalid (two dots + dot at beginning)
	{
		FieldStr: "..spec.template.spec.containers(0)",
	},
	// invalid (no index specified for array)
	{
		FieldStr: "spec.template.spec.containers()",
	},
}

func TestNextField(t *testing.T) {
	for i, test := range nextFieldTests {
		field, array, next := "", -1, test.FieldStr
		for fNum, curField := range test.NextFields {
			field, array, next = nextField(next)

			if len(test.NextFields) == 0 && (len(field)+len(next) != 0 || array < 0) {
				t.Errorf("test %d: should have returned nothing", i)
			}

			if field != curField.name {
				t.Errorf("test %d: field for step %d does not match. expected: %s, actual: %s", i, fNum,
					curField.name, field)
			}

			if array != curField.array {
				if curField.array != -1 {
					t.Errorf("test %d: expecting array for step %d", i, fNum)
				} else {
					t.Errorf("test %d: expecting %d (got %d) for step %d", i, curField.array, array, fNum)
				}
			}
		}
	}
}
