package project

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	git "gopkg.in/libgit2/git2go.v23"

	pb "rsprd.com/spread/pkg/spreadproto"
)

func (p *Project) GetDocument(revision, path string) (*pb.Document, error) {
	docs, err := p.ResolveCommit(revision)
	if err != nil {
		return nil, err
	}

	doc, has := docs[path]
	if !has {
		return nil, fmt.Errorf("the path '%s' does not exist in doc", path)
	}
	return doc, nil
}

func (p *Project) getDocument(oid *git.Oid) (*pb.Document, error) {
	blob, err := p.repo.LookupBlob(oid)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve Document blob: %v", err)
	}

	doc := &pb.Document{}
	err = proto.Unmarshal(blob.Contents(), doc)
	if err != nil {
		return nil, fmt.Errorf("unable to unmarshal document protobuf: %v", err)
	}
	return doc, nil
}

func (p *Project) createDocument(obj *pb.Document) (oid *git.Oid, size int, err error) {
	data, err := proto.Marshal(obj)
	if err != nil {
		err = fmt.Errorf("could not encode document: %v", err)
		return
	}
	size = len(data)

	oid, err = p.repo.CreateBlobFromBuffer(data)
	if err != nil {
		err = fmt.Errorf("could not write Document as blob in Git repo: %v", err)
		return
	}
	return
}
