package data

import (
	"encoding/json"
	"fmt"

	pb "rsprd.com/spread/pkg/spreadproto"
)

// GetFieldFromDocument returns a pointer to field in doc based on the given field path.
func GetFieldFromDocument(doc *pb.Document, fieldpath string) (*pb.Field, error) {
	root := doc.GetRoot()
	if root == nil {
		return nil, fmt.Errorf("root for document '%s' was nil", doc.Name)
	}

	field, err := ResolveRelativeField(root, fieldpath)
	if err != nil {
		return nil, fmt.Errorf("could not resolve '%s': %v", fieldpath, err)
	}
	return field, err
}

// CreateDocument creates a Document with an Object as it's value using CreateObject.
func CreateDocument(name, path string, ptr interface{}) (*pb.Document, error) {
	obj, err := CreateObject("", ptr)
	if err != nil {
		return nil, fmt.Errorf("could not create object for document: %v", err)
	}
	return &pb.Document{
		Name: name,
		Info: &pb.DocumentInfo{
			Path: path,
		},
		Root: &pb.Field{
			Value: &pb.Field_Object{
				Object: obj,
			},
		},
	}, nil
}

func Unmarshal(doc *pb.Document, ptr interface{}) error {
	fieldMap, err := MapFromDocument(doc)
	if err != nil {
		return err
	}

	// TODO: use reflection to replace this
	jsonData, err := json.Marshal(&fieldMap)
	if err != nil {
		return fmt.Errorf("unable to generate JSON from document data: %v", err)
	}

	return json.Unmarshal(jsonData, ptr)
}

func MapFromDocument(doc *pb.Document) (map[string]interface{}, error) {
	obj, err := getObjectFromDoc(doc)
	if err != nil {
		return nil, err
	}

	return MapFromObject(obj)
}

func getObjectFromDoc(doc *pb.Document) (*pb.Object, error) {
	root := doc.GetRoot()
	if root == nil {
		return nil, fmt.Errorf("document '%s' does not have a root", doc.Name)
	}

	obj := root.GetObject()
	if obj == nil {
		return nil, fmt.Errorf("root field of document '%s' does not have an object as it's value", doc.Name)
	}
	return obj, nil
}
