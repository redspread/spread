package data

import (
	"encoding/json"
	"errors"
	"fmt"

	pb "rsprd.com/spread/pkg/spreadproto"
)

// CreateObject uses reflection to convert the data (usually a struct) into an Object.
func CreateObject(name, path string, ptr interface{}) (*pb.Object, error) {
	data, err := json.Marshal(ptr)
	if err != nil {
		return nil, fmt.Errorf("unable to generate JSON object: %v", err)
	}

	// TODO: use reflection to replace this
	var jsonData map[string]interface{}
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		return nil, err
	}

	return ObjectFromMap(name, path, jsonData)
}

// ObjectFromMap creates an Object, using the entries of a map as fields.
// This supports maps embedded as values. It is assumed that types are limited to JSON types.
func ObjectFromMap(name, path string, data map[string]interface{}) (*pb.Object, error) {
	obj := &pb.Object{
		Name: name,
		Info: &pb.ObjectInfo{
			Path: path,
		},
	}

	i := 0
	obj.Fields = make([]*pb.Field, len(data))
	for k, v := range data {
		field, err := buildField(k, v)
		if err != nil {
			return nil, err
		}
		obj.Fields[i] = field
		i++
	}
	return obj, nil
}

func Unmarshal(obj *pb.Object, ptr interface{}) error {
	fieldMap, err := MapFromObject(obj)
	if err != nil {
		return err
	}

	// TODO: use reflection to replace this
	jsonData, err := json.Marshal(&fieldMap)
	if err != nil {
		return fmt.Errorf("unable to generate JSON from object data: %v", err)
	}

	return json.Unmarshal(jsonData, ptr)
}

func MapFromObject(obj *pb.Object) (map[string]interface{}, error) {
	fields := obj.GetFields()
	if fields == nil {
		return nil, ErrObjectNilFields
	}

	out := make(map[string]interface{}, len(fields))
	for _, field := range fields {
		val, err := decodeField(field)
		if err != nil {
			return nil, fmt.Errorf("could not decode field '%s': %v", field.Key, err)
		}
		out[field.Key] = val
	}
	return out, nil
}

var (
	ErrObjectNilFields = errors.New("object had nil for Fields")
)
