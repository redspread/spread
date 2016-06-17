package project

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	git "gopkg.in/libgit2/git2go.v23"

	pb "rsprd.com/spread/pkg/spreadproto"
)

func (p *Project) getObject(oid *git.Oid) (*pb.Object, error) {
	blob, err := p.repo.LookupBlob(oid)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve Object blob: %v", err)
	}

	obj := &pb.Object{}
	err = proto.Unmarshal(blob.Contents(), obj)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal object protobuf: %v", err)
	}
	return obj, nil
}
