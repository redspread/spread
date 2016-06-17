package data

import (
	"fmt"

	pb "rsprd.com/spread/pkg/spreadproto"
)

func buildField(key string, data interface{}) (*pb.Field, error) {
	if data == nil {
		return buildNil(key)
	}

	switch typedData := data.(type) {
	case bool:
		return buildBool(key, typedData)
	case float64:
		return buildNumber(key, typedData)
	case string:
		return buildString(key, typedData)
	case []interface{}:
		return buildArray(key, typedData)
	case map[string]interface{}:
		return buildMap(key, typedData)
	}

	return nil, fmt.Errorf("could not resolve type of %s (value=%+v)", key, data)
}

func buildNil(key string) (*pb.Field, error) {
	return &pb.Field{
		Key: key,
		Value: &pb.FieldValue{
			Type: pb.FieldValue_NULL,
		},
	}, nil
}

func buildBool(key string, data bool) (*pb.Field, error) {
	return &pb.Field{
		Key: key,
		Value: &pb.FieldValue{
			Type:  pb.FieldValue_BOOL,
			Value: fmt.Sprintf("%t", data),
		},
	}, nil
}

func buildNumber(key string, data float64) (*pb.Field, error) {
	return &pb.Field{
		Key: key,
		Value: &pb.FieldValue{
			Type:  pb.FieldValue_NUMBER,
			Value: fmt.Sprintf("%g", data),
		},
	}, nil
}

func buildString(key string, data string) (*pb.Field, error) {
	return &pb.Field{
		Key: key,
		Value: &pb.FieldValue{
			Type:  pb.FieldValue_STRING,
			Value: data,
		},
	}, nil
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
		Value: &pb.FieldValue{
			Type: pb.FieldValue_ARRAY,
		},
		Fields: arr,
	}, nil
}

func buildMap(key string, data map[string]interface{}) (*pb.Field, error) {
	arr := make([]*pb.Field, len(data))
	i := 0
	for k, v := range data {
		field, err := buildField(k, v)
		if err != nil {
			return nil, err
		}
		arr[i] = field
		i++
	}

	return &pb.Field{
		Key: key,
		Value: &pb.FieldValue{
			Type: pb.FieldValue_MAP,
		},
		Fields: arr,
	}, nil
}
