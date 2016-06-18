package project

import (
	"fmt"

	"github.com/golang/protobuf/proto"
	git "gopkg.in/libgit2/git2go.v23"

	"rsprd.com/spread/pkg/deploy"
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

func (p *Project) getKubeObject(oid *git.Oid, path string) (deploy.KubeObject, error) {
	obj, err := p.getObject(oid)
	if err != nil {
		return nil, fmt.Errorf("failed to read object from Git repository: %v", err)
	}

	kind, err := kindFromPath(path)
	if err != nil {
		return nil, err
	}

	kubeObj, err := deploy.KubeObjectFromObject(kind, obj)
	if err != nil {
		return nil, err
	}
	return kubeObj, nil
}
