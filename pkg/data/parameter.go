package data

import (
	"encoding/json"
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

// ApplyArguments takes the given arguments and uses them to satisfy a field parameter. If a single argument and no
// formatting pattern are given the single argument is used as the field value. Otherwise the arguments will be used as
// arguments to Printf with the formatting string as the pattern.
func ApplyArguments(field *pb.Field, args ...*pb.Argument) error {
	if field == nil {
		return errors.New("field was nil")
	} else if field.GetParam() == nil {
		return fmt.Errorf("field %s does not have a parameter", field.Key)
	} else if len(args) < 1 && field.GetParam().GetDefault() == nil {
		return errors.New("an argument must be specified if no default is given")
	} else if len(args) < 1 {
		return applyDefault(field)
	} else if len(args) == 1 && len(field.GetParam().Pattern) == 0 {
		return simpleArgApply(field, args[0])
	} else if len(args) > 1 && len(field.GetParam().Pattern) == 0 {
		return errors.New("may only use multiple arguments if a string template is provided")
	}

	argVals := make([]interface{}, len(args))
	for i, v := range args {
		switch val := v.GetValue().(type) {
		case *pb.Argument_Number:
			argVals[i] = val.Number
		case *pb.Argument_Str:
			argVals[i] = val.Str
		case *pb.Argument_Boolean:
			argVals[i] = val.Boolean
		}
	}

	val := fmt.Sprintf(field.GetParam().Pattern, argVals...)
	field.Value = &pb.Field_Str{Str: val}
	return nil
}

// ParseArgument returns an argument parsed from JSON.
func ParseArgument(in string) (*pb.Argument, error) {
	var data interface{}
	err := json.Unmarshal([]byte(in), &data)
	if err != nil {
		return nil, err
	}

	arg := &pb.Argument{}
	switch typedData := data.(type) {
	case bool:
		arg.Value = &pb.Argument_Boolean{
			Boolean: typedData,
		}
	case float64:
		arg.Value = &pb.Argument_Number{
			Number: typedData,
		}
	case string:
		arg.Value = &pb.Argument_Str{
			Str: typedData,
		}
	default:
		return nil, errors.New("unknown type")
	}
	return arg, nil
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
	default:
		field.Value = nil
	}

	return nil
}

func applyDefault(field *pb.Field) error {
	if field == nil {
		return errors.New("field cannot be nil")
	} else if field.GetParam() == nil {
		return errors.New("field does not have parameters")
	} else if field.GetParam().GetDefault() == nil {
		return errors.New("fields has paramaters but default was nil")
	}

	switch d := field.GetParam().GetDefault().GetValue().(type) {
	case *pb.Argument_Number:
		field.Value = &pb.Field_Number{Number: d.Number}
	case *pb.Argument_Str:
		field.Value = &pb.Field_Str{Str: d.Str}
	case *pb.Argument_Boolean:
		field.Value = &pb.Field_Boolean{Boolean: d.Boolean}
	}
	return nil
}
