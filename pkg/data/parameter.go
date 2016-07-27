package data

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"

	pb "rsprd.com/spread/pkg/spreadproto"
)

// InteractiveArgs hosts an interactive session using a Reader and Writer which prompts for input which is used to
// populate the provided field. If required is true then only fields without defaults will be prompted for.
func InteractiveArgs(r io.ReadCloser, w io.Writer, field *pb.Field, required bool) error {
	param := field.GetParam()

	defaultVal := param.GetDefault()
	// don't prompt if only checking for required and has default
	if required && defaultVal != nil {
		return nil
	}

	fmt.Fprintln(w, "Name: ", param.Name)
	fmt.Fprintln(w, "Prompt: ", param.Prompt)
	fmt.Fprintf(w, "Input [%s]: ", displayDefault(defaultVal))
	reader := bufio.NewReader(r)
	text, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	// use default if no input given
	args := []*pb.Argument{defaultVal}
	if len(text) > 1 {
		args, err = ParseArguments(text)
		if err != nil {
			return err
		}
	}

	return ApplyArguments(field, args...)
}

func displayDefault(d *pb.Argument) string {
	var out interface{}
	switch val := d.GetValue().(type) {
	case *pb.Argument_Number:
		out = val.Number
	case *pb.Argument_Str:
		out = val.Str
	case *pb.Argument_Boolean:
		out = val.Boolean
	}
	return fmt.Sprintf("%v", out)
}

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

// ParameterFields returns the fields with parameters contained within a Document.
func ParameterFields(docs map[string]*pb.Document) map[string]*pb.Field {
	fields := map[string]*pb.Field{}
	for _, doc := range docs {
		AddParameterFields(doc.GetRoot(), fields)
	}
	return fields
}

// AddParameterFields adds fields with parameters from the given field (and its subfields) to the map given. The name of the parameter is the key.
func AddParameterFields(field *pb.Field, params map[string]*pb.Field) {
	param := field.GetParam()
	// add to map if has parameter
	if param != nil {
		params[param.Name] = field
	}

	switch val := field.GetValue().(type) {
	case *pb.Field_Object:
		for _, objField := range val.Object.GetItems() {
			AddParameterFields(objField, params)
		}
	case *pb.Field_Array:
		for _, arrField := range val.Array.GetItems() {
			AddParameterFields(arrField, params)
		}
	}
}

// ParseArguments returns arguments parsed from JSON.
func ParseArguments(in string) (args []*pb.Argument, err error) {
	var data interface{}
	err = json.Unmarshal([]byte(in), &data)
	if err != nil {
		return
	}

	var dataArr []interface{}
	if arr, isArray := data.([]interface{}); isArray {
		dataArr = arr
	} else {
		dataArr = []interface{}{data}
	}

	args = make([]*pb.Argument, len(dataArr))
	for k, v := range dataArr {
		arg, err := argFromJSONType(v)
		if err != nil {
			return nil, err
		}
		args[k] = arg
	}
	return
}

func argFromJSONType(data interface{}) (*pb.Argument, error) {
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
		return nil, fmt.Errorf("unknown type: %+v", typedData)
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
