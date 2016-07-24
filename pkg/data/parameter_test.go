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
}

func TestApplyArguments(t *testing.T) {
	for i, test := range argTests {
		err := ApplyArguments(test.In, test.Args...)
		hasErr := err != nil
		if hasErr != test.Error {
			if test.Error {
				t.Errorf("test %d: should have returned error", i)
			} else {
				t.Errorf("test %d: shouldn't have errored. Error: %v", i, err)
			}
		}

		// TODO: check that value was correctly found
	}
}
