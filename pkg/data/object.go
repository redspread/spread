package data

import (
	"encoding/json"
	"fmt"

	pb "rsprd.com/spread/pkg/spreadproto"
)

// CreateObject uses reflection to convert the data (usually a struct) into an Object.
func CreateObject(name, path string, ptr interface{}) (*pb.Object, error) {
	data, err := json.Marshal(ptr)
	if err != nil {
		return nil, fmt.Errorf("unable to generate JSON object: %v", err)
	}

	// this is a bit hacky but not sure of a better way to ensure proper tagging
	var jsonData map[string]interface{}
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		return nil, err
	}

	return CreateObjectFromMap(name, path, jsonData)
}

// CreateObjectFromMap creates an Object, using the entries of a map as fields.
// This supports maps embedded as values. It is assumed that types are limited to JSON types.
func CreateObjectFromMap(name, path string, data map[string]interface{}) (*pb.Object, error) {
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
