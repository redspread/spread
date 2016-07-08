package data

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

type SRITest struct {
	in     string
	out    *SRI   // nil if error
	outStr string // expected string for success, prefix of error for failure
}

var goodSRIs = []SRITest{
	// treeish only
	{
		"a434f0ba11e6ec04ca640f90b854dddcecd0c8d9",
		&SRI{
			Treeish: "a434f0ba11e6ec04ca640f90b854dddcecd0c8d9",
		},
		"a434f0ba11e6ec04ca640f90b854dddcecd0c8d9",
	},
	// treeish only, with slash
	{
		"a434f0ba11e6ec04ca640f90b854dddcecd0c8d9/",
		&SRI{
			Treeish: "a434f0ba11e6ec04ca640f90b854dddcecd0c8d9",
		},
		"a434f0ba11e6ec04ca640f90b854dddcecd0c8d9",
	},
	// treeish only, with double slash
	{
		"a434f0ba11e6ec04ca640f90b854dddcecd0c8d9//",
		&SRI{
			Treeish: "a434f0ba11e6ec04ca640f90b854dddcecd0c8d9",
		},
		"a434f0ba11e6ec04ca640f90b854dddcecd0c8d9",
	},
	// shortened treeish
	{
		"a434f0b",
		&SRI{
			Treeish: "a434f0b",
		},
		"a434f0b",
	},
	// treeish & path
	{
		"e8f3ab9/default/replicationcontroller/web/",
		&SRI{
			Treeish: "e8f3ab9",
			Path:    "default/replicationcontroller/web",
		},
		"e8f3ab9/default/replicationcontroller/web",
	},
	// full SRI (no Field)
	{
		"e8f3ab9/default/replicationcontroller/web/?",
		&SRI{
			Treeish: "e8f3ab9",
			Path:    "default/replicationcontroller/web",
		},
		"e8f3ab9/default/replicationcontroller/web",
	},
	// full SRI
	{
		"e8f3ab9/default/replicationcontroller/web?spec.template.spec.containers(0)",
		&SRI{
			Treeish: "e8f3ab9",
			Path:    "default/replicationcontroller/web",
			Field:   "spec.template.spec.containers(0)",
		},
		"e8f3ab9/default/replicationcontroller/web?spec.template.spec.containers(0)",
	},
	{
		"a434f0ba11e6ec04ca640f90b854dddcecd0c8d9/default/replicationcontroller/web/?spec.template.spec.containers(0)",
		&SRI{
			Treeish: "a434f0ba11e6ec04ca640f90b854dddcecd0c8d9",
			Path:    "default/replicationcontroller/web",
			Field:   "spec.template.spec.containers(0)",
		},
		"a434f0ba11e6ec04ca640f90b854dddcecd0c8d9/default/replicationcontroller/web?spec.template.spec.containers(0)",
	},
	{
		"a434f0ba11e6ec04ca640f90b854dddcecd0c8d9//default//replicationcontroller//web//?spec.template.spec.containers(0)",
		&SRI{
			Treeish: "a434f0ba11e6ec04ca640f90b854dddcecd0c8d9",
			Path:    "default/replicationcontroller/web",
			Field:   "spec.template.spec.containers(0)",
		},
		"a434f0ba11e6ec04ca640f90b854dddcecd0c8d9/default/replicationcontroller/web?spec.template.spec.containers(0)",
	},
	{
		"a434f0ba11e6ec04ca640f90b854dddcecd0c8d9/default/replicationcontroller/web/?spec.template.spec.containers(0)(1)",
		&SRI{
			Treeish: "a434f0ba11e6ec04ca640f90b854dddcecd0c8d9",
			Path:    "default/replicationcontroller/web",
			Field:   "spec.template.spec.containers(0)(1)",
		},
		"a434f0ba11e6ec04ca640f90b854dddcecd0c8d9/default/replicationcontroller/web?spec.template.spec.containers(0)(1)",
	},
	// no treeish SRI
	{
		"*/default/replicationcontroller/web/?spec.template.spec.containers(1)",
		&SRI{
			Treeish: "*",
			Path:    "default/replicationcontroller/web",
			Field:   "spec.template.spec.containers(1)",
		},
		"*/default/replicationcontroller/web?spec.template.spec.containers(1)",
	},
}

func sfmt(s *SRI) string {
	if s == nil {
		s = new(SRI)
	}
	return fmt.Sprintf("treeish=%s, path=%s, field=%s", s.Treeish, s.Path, s.Field)
}

func TestParseGoodSRIs(t *testing.T) {
	for i, test := range goodSRIs {
		sri, err := ParseSRI(test.in)
		if err != nil {
			t.Errorf("%s(%d) failed with: %v", test.in, i, err)
		} else if !reflect.DeepEqual(sri, test.out) {
			t.Errorf("%s(%d):\n\thave %v\n\twant %v\n", test.in, i, sfmt(sri), sfmt(test.out))
		} else if sri.String() != test.outStr {
			t.Errorf("%s(%d) - Bad Serialization:\n\thave %s\n\twant %s\n", test.in, i, sri.String(), test.outStr)
		}
	}
}

var badSRIs = []SRITest{
	// empty string
	{
		"",
		nil,
		"git object ID was too short",
	},
	// short OID
	{
		"a434f",
		nil,
		"git object ID was too short",
	},
	// long oid
	// short OID
	{
		"a434f0ba11e6ec04ca640f90b854dddcecd0c8d9d",
		nil,
		"git object ID was too long",
	},
	// invalid characters in ID
	{
		"a343invalidID",
		nil,
		"invalid Treeish",
	},
	// invalid characters in path
	{
		"e8f3ab9/default/replication  controller/web?spec.template.spec.containers(0)",
		nil,
		"invalid Path",
	},
	// invalid characters in field
	{
		"e8f3ab9/default/replicationcontroller/web?spec.tem&&&&plate.spec.containers(0)",
		nil,
		"invalid Field",
	},
	// character in field array access
	{
		"e8f3ab9/default/replicationcontroller/web/?spec.template.spec.containers(d)",
		nil,
		"invalid Field",
	},
	// double field dot
	{
		"e8f3ab9/default/replicationcontroller/web/?spec..template.spec.containers(0)",
		nil,
		"invalid Field",
	},
	// start field with dot
	{
		"e8f3ab9/default/replicationcontroller/web?.spec.template.spec.containers(0)",
		nil,
		"invalid Field",
	},
	// unclosed parentheses
	{
		"e8f3ab9/default/replicationcontroller/web?spec.template.spec.containers(",
		nil,
		"invalid Field",
	},
	// unopened parentheses
	{
		"e8f3ab9/default/replicationcontroller/web?spec.template).spec.containers",
		nil,
		"invalid Field",
	},
	// already open parentheses
	{
		"e8f3ab9/default/replicationcontroller/web?spec.template.spec.co(ntainers(0)",
		nil,
		"invalid Field",
	},
}

func TestParseBadSRIs(t *testing.T) {
	for i, test := range badSRIs {
		_, err := ParseSRI(test.in)
		if err == nil {
			t.Errorf("%s(%d) did not return error (expected error prefix: %s)", test.in, i, test.outStr)
		} else if !strings.HasPrefix(err.Error(), test.outStr) {
			t.Errorf("%s(%d) wrong error: '%v' (expected error prefix: %s)", test.in, i, err.Error(), test.outStr)
		}
	}
}

// PartTest checks if parts are being properly created for SRIs
// rawsri is the input and the remaining fields are the output. Empty fields mean the related element was missing.
type PartTest struct {
	rawsri string
	oid    string
	path   string
	field  string
}

func (t PartTest) String() string {
	return fmt.Sprintf("rawsri=%s, oid=%s, path=%s, field=%s", t.rawsri, t.oid, t.path, t.field)
}

var partTests = []PartTest{
	{
		rawsri: "oid",
		oid:    "oid",
	},
	{
		rawsri: "oid/",
		oid:    "oid",
	},
	{
		rawsri: "oid/?",
		oid:    "oid",
	},
	{
		rawsri: "oid//////",
		oid:    "oid",
		path:   "/////",
	},
	{
		rawsri: "oid//////?",
		oid:    "oid",
		path:   "/////",
	},
	{
		rawsri: "oid//////?**",
		oid:    "oid",
		path:   "/////",
		field:  "**",
	},
	{
		rawsri: "oid??//////?**",
		oid:    "oid??",
		path:   "/////",
		field:  "**",
	},
	{
		rawsri: "oid//////?//",
		oid:    "oid",
		path:   "/////",
		field:  "//",
	},
	{
		rawsri: "oid???",
		oid:    "oid???",
	},
}

func TestParts(t *testing.T) {
	for i, expected := range partTests {
		input := expected.rawsri
		actual := PartTest{rawsri: input}
		actual.oid, actual.path, actual.field = parts(input)

		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Part %d:\n\thave %v\n\twant %v\n", i, actual, expected)
		}
	}
}
