package data

import (
	"fmt"

	pb "rsprd.com/spread/pkg/spreadproto"
)

func buildField(key string, data interface{}) (*pb.Field, error) {
	field := &pb.Field{
		Key: key,
	}

	// don't set Value for nil
	if data == nil {
		return field, nil
	}

	switch typedData := data.(type) {
	case bool:
		field.Value = &pb.Field_Boolean{
			Boolean: typedData,
		}
	case float64:
		field.Value = &pb.Field_Number{
			Number: typedData,
		}
	case string:
		field.Value = &pb.Field_Str{
			Str: typedData,
		}
	case []interface{}:
		return buildArray(key, typedData)
	case map[string]interface{}:
		return buildMap(key, typedData)
	}

	if field.Value == nil {
		return nil, fmt.Errorf("could not resolve type of %s (value=%+v)", key, data)
	}
	return field, nil
}

func buildArray(key string, data []interface{}) (*pb.Field, error) {
	arr := make([]*pb.Field, len(data))
	for k, v := range data {
		kStr := fmt.Sprintf("%d", k)
		field, err := buildField(kStr, v)
		if err != nil {
			return nil, err
		}
		arr[k] = field
	}

	return &pb.Field{
		Key: key,
		Value: &pb.Field_Array{
			Array: &pb.Array{
				Items: arr,
			},
		},
	}, nil
}

func buildMap(key string, data map[string]interface{}) (*pb.Field, error) {
	obj := make(map[string]*pb.Field, len(data))
	for k, v := range data {
		field, err := buildField(k, v)
		if err != nil {
			return nil, err
		}
		obj[k] = field
	}

	return &pb.Field{
		Key: key,
		Value: &pb.Field_Obj{
			Obj: &pb.Map{
				Item: obj,
			},
		},
	}, nil
}
