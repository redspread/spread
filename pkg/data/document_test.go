package data

import (
	"math/rand"

	pb "rsprd.com/spread/pkg/spreadproto"
)

func randomField() *pb.Field {
	key := randomString(rand.Intn(9))
	switch rand.Intn(5) {
	case 0: // number
		return &pb.Field{
			Key:   key,
			Value: &pb.Field_Number{Number: rand.Float64()},
		}
	case 1: // string
		return &pb.Field{
			Key:   key,
			Value: &pb.Field_Str{Str: randomString(40)},
		}
	case 2: // boolean
		b := rand.Intn(2) == 0
		return &pb.Field{
			Key:   key,
			Value: &pb.Field_Boolean{Boolean: b},
		}
	case 3: // object
		return &pb.Field{
			Key:   key,
			Value: &pb.Field_Object{Object: randomObject()},
		}
	case 4: // array
		return &pb.Field{
			Key:   key,
			Value: &pb.Field_Array{Array: randomArray()},
		}
	default:
		panic("illegal option")
	}
}

func randomObject() *pb.Object {
	numFields := rand.Intn(20)
	fields := make(map[string]*pb.Field, numFields)
	for i := 0; i < numFields; i++ {
		field := randomField()
		fields[field.Key] = field
	}
	return &pb.Object{
		Items: fields,
	}
}

func randomArray() *pb.Array {
	numFields := rand.Intn(20)
	fields := make([]*pb.Field, numFields)
	for i := 0; i < numFields; i++ {
		field := randomField()
		fields[i] = field
	}
	return &pb.Array{
		Items: fields,
	}
}

func randomString(strlen int) string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}
