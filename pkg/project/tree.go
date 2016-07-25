package project

import (
	"fmt"

	git "gopkg.in/libgit2/git2go.v23"

	pb "rsprd.com/spread/pkg/spreadproto"
)

func (p *Project) mapFromTree(tree *git.Tree) (docs map[string]*pb.Document, err error) {
	var walkErr error
	err = tree.Walk(func(path string, entry *git.TreeEntry) int {
		// add objects to deployment
		if entry.Type == git.ObjectBlob {
			doc, err := p.getDocument(entry.Id)
			if err != nil {
				walkErr = err
				return -1
			}

			docs[path] = doc
			if walkErr != nil {
				return -1
			}
		}
		return 0
	})

	if err != nil {
		err = fmt.Errorf("error starting walk: %v", err)
	} else if walkErr != nil {
		err = fmt.Errorf("error during walk: %v", walkErr)
	}
	return
}
