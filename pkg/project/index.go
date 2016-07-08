package project

import (
	"errors"
	"fmt"
	"strings"

	git "gopkg.in/libgit2/git2go.v23"

	"rsprd.com/spread/pkg/deploy"
	pb "rsprd.com/spread/pkg/spreadproto"
)

func (p *Project) AddDocumentToIndex(doc *pb.Document) error {
	oid, size, err := p.createDocument(doc)
	if err != nil {
		return fmt.Errorf("object couldn't be created: %v", err)
	}

	info := doc.GetInfo()
	if info == nil {
		return ErrNilObjectInfo
	}

	entry := &git.IndexEntry{
		Mode: git.FilemodeBlob,
		Size: uint32(size),
		Id:   oid,
		Path: info.Path,
	}

	index, err := p.repo.Index()
	if err != nil {
		return fmt.Errorf("could not retreive index: %v", err)
	}

	err = index.Add(entry)
	if err != nil {
		return err
	}

	return index.Write()
}

func (p *Project) Index() (*deploy.Deployment, error) {
	index, err := p.repo.Index()
	if err != nil {
		return nil, fmt.Errorf("could not retrieve index: %v", err)
	}

	deployment := new(deploy.Deployment)
	indexSize := int(index.EntryCount())
	for i := 0; i < indexSize; i++ {
		entry, err := index.EntryByIndex(uint(i))
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve index entry: %v", err)
		}

		kubeObj, err := p.getKubeObject(entry.Id, entry.Path)
		if err != nil {
			return nil, err
		}

		err = deployment.Add(kubeObj)
		if err != nil {
			return nil, fmt.Errorf("could not add object to deployment: %v", err)
		}
	}
	return deployment, nil
}

func kindFromPath(path string) (string, error) {
	parts := strings.Split(path, "/")
	if len(parts) != 4 {
		return "", fmt.Errorf("path wrong length (is %d, expected 5)", len(parts))
	}
	return parts[2], nil
}

var (
	ErrNilObjectInfo = errors.New("an object's Info field cannot be nil")
)
