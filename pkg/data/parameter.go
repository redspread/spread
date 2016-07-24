package data

import (
	"errors"
	"fmt"

	pb "rsprd.com/spread/pkg/spreadproto"
)

// AddParamToDoc adds the given parameter to the
func AddParamToDoc(doc *pb.Document, target *SRI, param *pb.Parameter) error {
	if !target.IsField() {
		return errors.New("passed SRI is not a field")
	}

	field, err := GetFieldFromDocument(doc, target.Field)
	if err != nil {
		return err
	}

	field.Param = param
	return nil
}

// ApplyArguments takes the given arguments and uses them to satisfy a field parameter.
func ApplyArguments(field *pb.Field, args ...*pb.Argument) error {
	if field == nil {
		return errors.New("field was nil")
	} else if field.GetParam() == nil {
		return fmt.Errorf("field %s does not have a parameter", field.Key)
	} else if len(args) < 1 {
		return errors.New("an argument must be specified")
	} else if len(args) == 1 && len(field.GetParam().Pattern) == 0 {
		return simpleArgApply(field, args[0])
	}
	// TODO: complete string formatting based apply
	return nil
}

// simpleArgApply is used when no formatting template string is given.
func simpleArgApply(field *pb.Field, arg *pb.Argument) error {
	switch val := arg.GetValue().(type) {
	case *pb.Argument_Number:
		field.Value = &pb.Field_Number{Number: val.Number}
	case *pb.Argument_Str:
		field.Value = &pb.Field_Str{Str: val.Str}
	case *pb.Argument_Boolean:
		field.Value = &pb.Field_Boolean{Boolean: val.Boolean}
	case *pb.Argument_Object:
		field.Value = &pb.Field_Object{Object: val.Object}
	case *pb.Argument_Array:
		field.Value = &pb.Field_Array{Array: val.Array}
	default:
		field.Value = nil
	}

	return nil
}
