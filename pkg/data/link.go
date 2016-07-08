package data

import (
	"errors"

	pb "rsprd.com/spread/pkg/spreadproto"
)

// NewLink creates a new link from with the given details.
func NewLink(packageName string, target *SRI, override bool) *pb.Link {
	return &pb.Link{
		PackageName: packageName,
		Target:      target.Proto(),
		Override:    override,
	}
}

// CreateLinkInDocument creates a link from source to target with document.
func CreateLinkInDocument(doc *pb.Document, target *pb.Link, source *SRI) error {
	if !source.IsField() {
		return errors.New("passed SRI is not a field")
	}

	field, err := GetFieldFromDocument(doc, source.Field)
	if err != nil {
		return err
	}

	field.Value = &pb.Field_Link{
		Link: target,
	}
	return nil
}
