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
	array bool
}

var nextFieldTests = []NextFieldTest{
	{
		"spec.template.spec.containers(0)",
		[]NextField{
			{name: "spec"}, {name: "template"}, {name: "spec"}, {name: "containers"},
			{name: "0", array: true},
		},
	},
	{
		"(0)(1)(2)(3)(4)(5)ape",
		[]NextField{
			{name: "0", array: true}, {name: "1", array: true}, {name: "2", array: true},
			{name: "3", array: true}, {name: "4", array: true}, {name: "5", array: true},
			{name: "ape"},
		},
	},
	{
		"spec(0)helo.eat.cheese",
		[]NextField{
			{name: "spec"}, {name: "0", array: true}, {name: "helo"}, {name: "eat"}, {name: "cheese"},
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
		field, next, array := "", test.FieldStr, false
		for fNum, curField := range test.NextFields {
			field, next, array = nextField(next)

			if len(test.NextFields) == 0 && (len(field)+len(next) != 0 || array) {
				t.Errorf("test %d: should have returned nothing", i)
			}

			if field != curField.name {
				t.Errorf("test %d: field for step %d does not match. expected: %s, actual: %s", i, fNum,
					curField.name, field)
			}

			if array != curField.array {
				expected := "field"
				if curField.array {
					expected = "array"
				}
				t.Errorf("test %d: expecting %s for step %d", i, expected, fNum)
			}
		}
	}
}
