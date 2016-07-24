package data

import (
	"testing"

	pb "rsprd.com/spread/pkg/spreadproto"
)

type ArgTest struct {
	In    *pb.Field
	Args  []*pb.Argument
	Out   *pb.Field
	Error bool
}

var argTests = []ArgTest{
	{ // nil field
		In:    nil,
		Error: true,
	},
	{ // no args
		In:    randomField(),
		Args:  []*pb.Argument{},
		Error: true,
	},
	{ // no parameter
		In: &pb.Field{
			Key:   "testField",
			Value: &pb.Field_Str{Str: "original"},
		},
		Args: []*pb.Argument{
			{Value: &pb.Argument_Str{Str: "test"}},
			{Value: &pb.Argument_Boolean{Boolean: true}},
		},
		Error: true,
	},
	{ // default value
		In: &pb.Field{
			Key:   "testField",
			Value: &pb.Field_Str{Str: "original"},
			Param: &pb.Parameter{
				Default: &pb.Argument{
					Value: &pb.Argument_Str{Str: "haldo"},
				},
			},
		},
		Args: []*pb.Argument{},
		Out: &pb.Field{
			Key:   "testField",
			Value: &pb.Field_Str{Str: "haldo"},
		},
	},
	{ // simple sub
		In: &pb.Field{
			Key:   "testField",
			Value: &pb.Field_Str{Str: "original"},
			Param: &pb.Parameter{
				Default: &pb.Argument{
					Value: &pb.Argument_Str{Str: "haldo"},
				},
			},
		},
		Args: []*pb.Argument{
			{
				Value: &pb.Argument_Number{Number: 3.34334},
			},
		},
		Out: &pb.Field{
			Key:   "testField",
			Value: &pb.Field_Number{Number: 3.34334},
		},
	},
	{ // sub -- too many args
		In: &pb.Field{
			Key:   "testField",
			Value: &pb.Field_Str{Str: "original"},
			Param: &pb.Parameter{
				Default: &pb.Argument{
					Value: &pb.Argument_Str{Str: "haldo"},
				},
			},
		},
		Args: []*pb.Argument{
			{
				Value: &pb.Argument_Number{Number: 3.34334},
			},
			{
				Value: &pb.Argument_Str{Str: "HELO"},
			},
		},
		Error: true,
	},
	{ // formatted string
		In: &pb.Field{
			Key:   "testField",
			Value: &pb.Field_Str{Str: "original"},
			Param: &pb.Parameter{
				Pattern: "we are %s",
			},
		},
		Args: []*pb.Argument{
			{
				Value: &pb.Argument_Str{Str: "maryland"},
			},
		},
		Out: &pb.Field{
			Key:   "testField",
			Value: &pb.Field_Str{Str: "we are maryland"},
		},
	},
}

func TestApplyArguments(t *testing.T) {
	for i, test := range argTests {
		field := test.In
		err := ApplyArguments(field, test.Args...)
		hasErr := err != nil
		if !hasErr && test.Error {
			t.Errorf("test %d: should have returned error", i)
		} else if hasErr && !test.Error {
			t.Errorf("test %d: shouldn't have errored. Error: %v", i, err)
		} else if !test.Error && !FieldValueEquals(field, test.Out) {
			t.Errorf("test %d: field values don't match. expected: '%+v', actual: '%+v'", i, test.Out.GetValue(), test.In.GetValue())
		}
	}
}
