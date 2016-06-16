package project

import (
	"errors"
	"fmt"

	"github.com/golang/protobuf/proto"
	git "gopkg.in/libgit2/git2go.v23"

	pb "rsprd.com/spread/pkg/spreadproto"
)

func (p *Project) AddObjectToIndex(obj *pb.Object) error {
	info := obj.GetInfo()
	if info == nil {
		return ErrNilObjectInfo
	}

	data, err := proto.Marshal(obj)
	if err != nil {
		return fmt.Errorf("could not encode object: %v", err)
	}

	oid, err := p.repo.CreateBlobFromBuffer(data)
	if err != nil {
		return fmt.Errorf("could not write Object as blob in Git repo: %v", err)
	}

	entry := &git.IndexEntry{
		Mode: git.FilemodeBlob,
		Size: uint32(len(data)),
		Id:   oid,
		Path: info.Path,
	}

	index, err := p.repo.Index()
	if err != nil {
		return fmt.Errorf("could not retreive index: %v", err)
	}

	return index.Add(entry)
}

var (
	ErrNilObjectInfo = errors.New("an object's Info field cannot be nil")
)
