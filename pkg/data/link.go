package data

import (
	pb "rsprd.com/spread/pkg/spreadproto"
)

// NewLink creates a new link from with the given details.
func NewLink(packageName string, target *SRI, override bool) *pb.Link {
	pbTarget := pb.SRI(*target)
	return &pb.Link{
		PackageName: packageName,
		Target:      &pbTarget,
		Override:    override,
	}
}

// CreateLinkInObject creates a link from source to target with object.
func CreateLinkInObject(obj *pb.Object, target *pb.Link, source *SRI) error {
	field, err := GetFieldFromObject(obj, source)
	if err != nil {
		return err
	}

	field.Link = target
	return nil
}
