package data

import (
	"fmt"
	"strconv"

	pb "rsprd.com/spread/pkg/spreadproto"
)

func ResolveRelativeField(field *pb.Field, fieldpath string) (resolvedField *pb.Field, err error) {
	fieldKey, arrIndex, next := nextField(fieldpath)
	if arrIndex >= 0 {
		resolvedField, err = getFromArrayField(field, arrIndex)
	} else if len(fieldKey) > 0 {
		resolvedField, err = getFromMapField(field, fieldKey)
	} else {
		err = fmt.Errorf("could not resolve fieldpath '%s'", fieldpath)
	}

	// return if err or no more fields to traverse (end of field path)
	if err != nil || len(next) == 0 {
		return
	}
	return ResolveRelativeField(resolvedField, next)
}

func getFromArrayField(field *pb.Field, index int) (*pb.Field, error) {
	fieldArr := field.GetArray()
	if fieldArr == nil {
		return nil, fmt.Errorf("field '%s' isn't an array, cannot access %s[%d]", field.Key, field.Key, index)
	}

	items := fieldArr.GetItems()
	if items == nil {
		return nil, fmt.Errorf("the array wrapper struct for the value of field '%s' had nil for items, cannot access %s[%d]", field.Key, field.Key, index)
	} else if len(items)-1 < index {
		return nil, fmt.Errorf("could not access %s[%d], the size of '%s' is %d", field.Key, index, field.Key, len(items))
	}
	return items[index], nil
}

func getFromMapField(field *pb.Field, key string) (*pb.Field, error) {
	fieldMap := field.GetObject()
	if fieldMap == nil {
		return nil, fmt.Errorf("field '%s' isn't an object, cannot access %s['%s']", field.Key, field.Key, key)
	}

	items := fieldMap.GetItems()
	if items == nil {
		return nil, fmt.Errorf("the object wrapper struct for the value of field '%s' had nil for items, cannot access %s[%d]", field.Key, field.Key, key)
	}

	item, ok := items[key]
	if !ok {
		return nil, fmt.Errorf("no key '%s' in map for field '%s", key, field.Key)
	}
	return item, nil
}

// nextField returns the first field in a fieldpath and returns the remainder after removing the root element.
// It will return array as positive number or 0 if refers to array. If there is no next field, an empty string in field will be returned.
func nextField(fieldpath string) (field string, array int, next string) {
	array = -1
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
			indexStr := fieldpath[1:i]
			if len(fieldpath) > i+2 {
				next = fieldpath[i+1:]
			}

			if num, err := strconv.ParseInt(indexStr, 10, 64); err == nil {
				array = int(num)
			}
			return
		} else if i+1 == len(fieldpath) {
			field = fieldpath
			return
		}
	}
	return
}
