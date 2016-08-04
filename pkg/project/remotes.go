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
		if cert.Kind == git.CertificateHostkey {
			return git.ErrOk
		} else if valid {
			return git.ErrOk
		}
		return git.ErrAuth
	},
}

func (p *Project) Remotes() *git.RemoteCollection {
	return &p.repo.Remotes
}

func (p *Project) Push(remoteName string, refspecs ...string) error {
	remote, err := p.Remotes().Lookup(remoteName)
	if err != nil {
		return fmt.Errorf("Failed to lookup branch: %v", err)
	}

	opts := &git.PushOptions{
		RemoteCallbacks: remoteCallbacks,
	}
	err = remote.Push(refspecs, opts)
	if err != nil {
		return fmt.Errorf("Failed to push: %v", err)
	}
	return nil
}

func (p *Project) Fetch(remoteName string, refspecs ...string) error {
	remote, err := p.Remotes().Lookup(remoteName)
	if err != nil {
		return fmt.Errorf("Failed to lookup remote: %v", err)
	}

	return p.fetch(remote, refspecs...)
}

func (p *Project) FetchAnonymous(url string, refspecs ...string) error {
	remote, err := p.Remotes().CreateAnonymous(url)
	if err != nil {
		return fmt.Errorf("Failed to create anonymous remote for '%s': %v", url, err)
	}

	return p.fetch(remote, refspecs...)
}

func (p *Project) fetch(remote *git.Remote, refspecs ...string) (err error) {
	opts := &git.FetchOptions{
		RemoteCallbacks: remoteCallbacks,
	}

	// fetch with default reflog message
	err = remote.Fetch(refspecs, opts, "")
	if err != nil {
		return fmt.Errorf("Failed to fetch: %v", err)
	}
	return
}
