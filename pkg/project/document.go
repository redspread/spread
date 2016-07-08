package project

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	git "gopkg.in/libgit2/git2go.v23"

	"rsprd.com/spread/pkg/deploy"
	pb "rsprd.com/spread/pkg/spreadproto"
)

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

func (p *Project) getKubeObject(oid *git.Oid, path string) (deploy.KubeObject, error) {
	doc, err := p.getDocument(oid)
	if err != nil {
		return nil, fmt.Errorf("failed to read object from Git repository: %v", err)
	}

	kind, err := kindFromPath(path)
	if err != nil {
		return nil, err
	}

	kubeObj, err := deploy.KubeObjectFromDocument(kind, doc)
	if err != nil {
		return nil, err
	}
	return kubeObj, nil
}
