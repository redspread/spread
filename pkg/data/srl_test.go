package data

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

type SRLTest struct {
	in     string
	out    *SRL   // nil if error
	outStr string // expected string for success, prefix of error for failure
}

var goodSRLs = []SRLTest{
	// treeish only
	{
		"a434f0ba11e6ec04ca640f90b854dddcecd0c8d9",
		&SRL{
			Treeish: "a434f0ba11e6ec04ca640f90b854dddcecd0c8d9",
		},
		"a434f0ba11e6ec04ca640f90b854dddcecd0c8d9",
	},
	// treeish only, with slash
	{
		"a434f0ba11e6ec04ca640f90b854dddcecd0c8d9/",
		&SRL{
			Treeish: "a434f0ba11e6ec04ca640f90b854dddcecd0c8d9",
		},
		"a434f0ba11e6ec04ca640f90b854dddcecd0c8d9",
	},
	// treeish only, with double slash
	{
		"a434f0ba11e6ec04ca640f90b854dddcecd0c8d9//",
		&SRL{
			Treeish: "a434f0ba11e6ec04ca640f90b854dddcecd0c8d9",
		},
		"a434f0ba11e6ec04ca640f90b854dddcecd0c8d9",
	},
	// shortened treeish
	{
		"a434f0b",
		&SRL{
			Treeish: "a434f0b",
		},
		"a434f0b",
	},
	// treeish & path
	{
		"e8f3ab9/default/replicationcontroller/web/",
		&SRL{
			Treeish: "e8f3ab9",
			Path:    "default/replicationcontroller/web",
		},
		"e8f3ab9/default/replicationcontroller/web",
	},
	// full SRL (no Field)
	{
		"e8f3ab9/default/replicationcontroller/web/?",
		&SRL{
			Treeish: "e8f3ab9",
			Path:    "default/replicationcontroller/web",
		},
		"e8f3ab9/default/replicationcontroller/web",
	},
	// full SRL
	{
		"e8f3ab9/default/replicationcontroller/web?spec.template.spec.containers(0)",
		&SRL{
			Treeish: "e8f3ab9",
			Path:    "default/replicationcontroller/web",
			Field:   "spec.template.spec.containers(0)",
		},
		"e8f3ab9/default/replicationcontroller/web?spec.template.spec.containers(0)",
	},
	{
		"a434f0ba11e6ec04ca640f90b854dddcecd0c8d9/default/replicationcontroller/web/?spec.template.spec.containers(0)",
		&SRL{
			Treeish: "a434f0ba11e6ec04ca640f90b854dddcecd0c8d9",
			Path:    "default/replicationcontroller/web",
			Field:   "spec.template.spec.containers(0)",
		},
		"a434f0ba11e6ec04ca640f90b854dddcecd0c8d9/default/replicationcontroller/web?spec.template.spec.containers(0)",
	},
	{
		"a434f0ba11e6ec04ca640f90b854dddcecd0c8d9//default//replicationcontroller//web//?spec.template.spec.containers(0)",
		&SRL{
			Treeish: "a434f0ba11e6ec04ca640f90b854dddcecd0c8d9",
			Path:    "default/replicationcontroller/web",
			Field:   "spec.template.spec.containers(0)",
		},
		"a434f0ba11e6ec04ca640f90b854dddcecd0c8d9/default/replicationcontroller/web?spec.template.spec.containers(0)",
	},
	{
		"a434f0ba11e6ec04ca640f90b854dddcecd0c8d9/default/replicationcontroller/web/?spec.template.spec.containers(0)(1)",
		&SRL{
			Treeish: "a434f0ba11e6ec04ca640f90b854dddcecd0c8d9",
			Path:    "default/replicationcontroller/web",
			Field:   "spec.template.spec.containers(0)(1)",
		},
		"a434f0ba11e6ec04ca640f90b854dddcecd0c8d9/default/replicationcontroller/web?spec.template.spec.containers(0)(1)",
	},
}

func sfmt(s *SRL) string {
	if s == nil {
		s = new(SRL)
	}
	return fmt.Sprintf("treeish=%s, path=%s, field=%s", s.Treeish, s.Path, s.Field)
}

func TestParseGoodSRLs(t *testing.T) {
	for i, test := range goodSRLs {
		srl, err := ParseSRL(test.in)
		if err != nil {
			t.Errorf("%s(%d) failed with: %v", test.in, i, err)
		} else if !reflect.DeepEqual(srl, test.out) {
			t.Errorf("%s(%d):\n\thave %v\n\twant %v\n", test.in, i, sfmt(srl), sfmt(test.out))
		} else if srl.String() != test.outStr {
			t.Errorf("%s(%d) - Bad Serialization:\n\thave %s\n\twant %s\n", test.in, i, srl.String(), test.outStr)
		}
	}
}

var badSRLs = []SRLTest{
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

func TestParseBadSRLs(t *testing.T) {
	for i, test := range badSRLs {
		_, err := ParseSRL(test.in)
		if err == nil {
			t.Errorf("%s(%d) did not return error (expected error prefix: %s)", test.in, i, test.outStr)
		} else if !strings.HasPrefix(err.Error(), test.outStr) {
			t.Errorf("%s(%d) wrong error: '%v' (expected error prefix: %s)", test.in, i, err.Error(), test.outStr)
		}
	}
}

// PartTest checks if parts are being properly created for SRLs
// rawsrl is the input and the remaining fields are the output. Empty fields mean the related element was missing.
type PartTest struct {
	rawsrl string
	oid    string
	path   string
	field  string
}

func (t PartTest) String() string {
	return fmt.Sprintf("rawsrl=%s, oid=%s, path=%s, field=%s", t.rawsrl, t.oid, t.path, t.field)
}

var partTests = []PartTest{
	{
		rawsrl: "oid",
		oid:    "oid",
	},
	{
		rawsrl: "oid/",
		oid:    "oid",
	},
	{
		rawsrl: "oid/?",
		oid:    "oid",
	},
	{
		rawsrl: "oid//////",
		oid:    "oid",
		path:   "/////",
	},
	{
		rawsrl: "oid//////?",
		oid:    "oid",
		path:   "/////",
	},
	{
		rawsrl: "oid//////?**",
		oid:    "oid",
		path:   "/////",
		field:  "**",
	},
	{
		rawsrl: "oid??//////?**",
		oid:    "oid??",
		path:   "/////",
		field:  "**",
	},
	{
		rawsrl: "oid//////?//",
		oid:    "oid",
		path:   "/////",
		field:  "//",
	},
	{
		rawsrl: "oid???",
		oid:    "oid???",
	},
}

func TestParts(t *testing.T) {
	for i, expected := range partTests {
		input := expected.rawsrl
		actual := PartTest{rawsrl: input}
		actual.oid, actual.path, actual.field = parts(input)

		if !reflect.DeepEqual(actual, expected) {
			t.Errorf("Part %d:\n\thave %v\n\twant %v\n", i, actual, expected)
		}
	}
}
