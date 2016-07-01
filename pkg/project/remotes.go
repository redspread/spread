package project

import (
	"fmt"

	"github.com/mitchellh/go-homedir"

	git "gopkg.in/libgit2/git2go.v23"
)

var remoteCallbacks = git.RemoteCallbacks{
	CredentialsCallback: func(url string, username_from_url string, allowed_types git.CredType) (git.ErrorCode, *git.Cred) {
		pubKey, err := homedir.Expand("~/.ssh/id_rsa.pub")
		if err != nil {
			return git.ErrAuth, nil
		}

		privKey, err := homedir.Expand("~/.ssh/id_rsa")
		if err != nil {
			return git.ErrAuth, nil
		}

		code, key := git.NewCredSshKey("git", pubKey, privKey, "")
		return git.ErrorCode(code), &key
	},
	CertificateCheckCallback: func(cert *git.Certificate, valid bool, hostname string) git.ErrorCode {
		// TODO: MAJOR SECURITY VULNERABILITY!!!!, resolve ASAP
		return git.ErrOk
	},
}

func (p *Project) Remotes() *git.RemoteCollection {
	return &p.repo.Remotes
}

func (p *Project) Push(remoteName string, refspec ...string) error {
	remote, err := p.Remotes().Lookup(remoteName)
	if err != nil {
		return fmt.Errorf("Failed to lookup branch: %v", err)
	}

	pushOpts := &git.PushOptions{
		RemoteCallbacks: remoteCallbacks,
	}
	err = remote.Push(refspec, pushOpts)
	if err != nil {
		return fmt.Errorf("Failed to push: %v", err)
	}
	return nil
}
