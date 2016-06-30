package data

import (
	"fmt"

	pb "rsprd.com/spread/pkg/spreadproto"
)

// Fields provides helper methods for working with protobuf object fields.
type Fields []*pb.Field

// ResolveFields returns a field based on the provided field path in the format used for SRIs.
// An error is returned if the given path doesn't exist.
func (f Fields) ResolveField(fieldpath string) (*pb.Field, error) {
	field, next, _ := nextField(fieldpath)
	if len(field) == 0 {
		return nil, fmt.Errorf("could not resolve fieldpath '%s'", fieldpath)
	} else if len(next) > 0 {
		fields := f.GetFields(field)
		if fields != nil {
			return fields.ResolveField(next)
		}
	} else if out := f.Get(field); out != nil {
		return out, nil
	}

	return nil, fmt.Errorf("could not find field '%s'", field)
}

// Get returns a field by name
func (f Fields) Get(name string) *pb.Field {
	for _, field := range f {
		if field.Key == name {
			return field
		}
	}
	return nil
}

// GetFields returns the sub-fields of a field by name through an O(n) operation. Nil is returned if no field exists.
func (f Fields) GetFields(name string) Fields {
	field := f.Get(name)
	if field == nil {
		return nil
	}

	if field.Fields != nil {
		return Fields(field.Fields)
	}
	return nil
}

// nextField returns the first field in a fieldpath and returns the remainder after removing the root element.
// It will return array as true if the field is an array. If there is no next field, an empty string in field will be returned.
func nextField(fieldpath string) (field, next string, array bool) {
	fieldpath, err := ValidateField(fieldpath)
	if err != nil {
		return
	}

	for i, c := range fieldpath {
		// check for end of field chars
		if (c == '.' || c == '(') && i > 0 {
			field = fieldpath[:i]
			if c == '.' {
				i++
			}

			if len(fieldpath) > i+1 {
				next = fieldpath[i:]
			}
			return
		} else if c == ')' {
			field = fieldpath[1:i]
			if len(fieldpath) > i+2 {
				next = fieldpath[i+1:]
			}
			array = true
			return
		} else if i+1 == len(fieldpath) {
			field = fieldpath
			return
		}
	}
	return
}
