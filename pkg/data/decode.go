package data

import (
	"fmt"
	"strconv"

	pb "rsprd.com/spread/pkg/spreadproto"
)

func decodeField(field *pb.Field) (interface{}, error) {
	val := field.GetValue()
	if val == nil {
		return nil, fmt.Errorf("value for '%s' was nil", field.Key)
	}

	switch val.Type {
	case pb.FieldValue_NUMBER:
		return strconv.ParseFloat(val.Value, 64)
	case pb.FieldValue_STRING:
		return val.Value, nil
	case pb.FieldValue_BOOL:
		return strconv.ParseBool(val.Value)
	case pb.FieldValue_NULL:
		return nil, nil
	case pb.FieldValue_MAP:
		return decodeMapField(field)
	case pb.FieldValue_ARRAY:
		return decodeArrayField(field)
	}

	return nil, fmt.Errorf("unknown type for Field '%s'", field.Key)
}

func decodeMapField(root *pb.Field) (map[string]interface{}, error) {
	fields := root.GetFields()
	if fields == nil {
		return nil, nil
	}

	out := make(map[string]interface{}, len(fields))
	for _, field := range fields {
		val, err := decodeField(field)
		if err != nil {
			return nil, fmt.Errorf("couldn't decode '%s': %v", field.Key, err)
		}
		out[field.Key] = val
	}
	return out, nil
}

func decodeArrayField(root *pb.Field) ([]interface{}, error) {
	fields := root.GetFields()
	if fields == nil {
		return nil, nil
	}

	out := make([]interface{}, len(fields))
	for i, field := range fields {
		val, err := decodeField(field)
		if err != nil {
			return nil, fmt.Errorf("couldn't decode '%s': %v", field.Key, err)
		}
		out[i] = val
	}
	return out, nil
}
