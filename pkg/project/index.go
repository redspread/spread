package project

import (
	"errors"
	"fmt"

	git "gopkg.in/libgit2/git2go.v23"

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

func (p *Project) Index() (docs map[string]*pb.Document, err error) {
	index, err := p.repo.Index()
	if err != nil {
		return nil, fmt.Errorf("could not retrieve index: %v", err)
	}
	indexSize := int(index.EntryCount())
	docs = make(map[string]*pb.Document, indexSize)
	for i := 0; i < indexSize; i++ {
		entry, err := index.EntryByIndex(uint(i))
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve index entry: %v", err)
		}

		doc, err := p.getDocument(entry.Id)
		if err != nil {
			return nil, err
		}
		docs[entry.Path] = doc
	}
	return
}

func (p *Project) DocFromIndex(path string) (*pb.Document, error) {
	index, err := p.repo.Index()
	if err != nil {
		return nil, fmt.Errorf("could not retrieve index: %v", err)
	}

	indexSize := int(index.EntryCount())
	for i := 0; i < indexSize; i++ {
		entry, err := index.EntryByIndex(uint(i))
		if err != nil {
			return nil, fmt.Errorf("failed to retrieve index entry: %v", err)
		}

		if entry.Path == path {
			doc, err := p.getDocument(entry.Id)
			if err != nil {
				return nil, fmt.Errorf("failed to read object from Git repository: %v", err)
			}
			return doc, nil
		}
	}
	return nil, fmt.Errorf("could not find document with path '%s'", path)
}

var (
	ErrNilObjectInfo = errors.New("an object's Info field cannot be nil")
)
