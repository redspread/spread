package data

import (
	"fmt"

	pb "rsprd.com/spread/pkg/spreadproto"
)

func decodeField(field *pb.Field) (interface{}, error) {
	val := field.GetValue()
	if val == nil {
		return nil, nil
	}

	switch v := val.(type) {
	case *pb.Field_Number:
		return v.Number, nil
	case *pb.Field_Str:
		return v.Str, nil
	case *pb.Field_Boolean:
		return v.Boolean, nil
	case *pb.Field_Object:
		return decodeObject(v.Object.GetItems())
	case *pb.Field_Array:
		return decodeArray(v.Array.GetItems())
	case *pb.Field_Link:
		// TODO: IMPLEMENT FOLLOWING LINKS
		return nil, nil
	}

	return nil, fmt.Errorf("unknown type for Field '%s'", field.Key)
}

func decodeObject(fields map[string]*pb.Field) (map[string]interface{}, error) {
	out := make(map[string]interface{}, len(fields))
	for k, field := range fields {
		val, err := decodeField(field)
		if err != nil {
			return nil, fmt.Errorf("couldn't decode '%s': %v", field.Key, err)
		}
		out[k] = val
	}
	return out, nil
}

func decodeArray(fields []*pb.Field) ([]interface{}, error) {
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
