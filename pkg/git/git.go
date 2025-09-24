// Package git provides a Git client for performing Git operations such as
// cloning, pulling, committing, and pushing changes to a repository.
package git

import (
	"fmt"

	"golang.org/x/crypto/ssh/knownhosts"
	gogit "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

// Client contains an authentication method and a repository object.
// It is used for all Git related operations.
type Client struct {
	auth transport.AuthMethod
	repo *gogit.Repository
}

// New creates a new Git Client with an SSH key, the repository
// URL and an absolute path where to read/write the repository.
// SSH was chosen as it is available on all Git servers (Github, Gitea, Gitlab, etc.)
func New(sshKnownHostsPath, sshKeyPath, sshKeyPassword,
	url, absolutePath string,
) (client *Client, err error) {
	// Only PEM private keys supported
	auth, err := ssh.NewPublicKeysFromFile("git", sshKeyPath, sshKeyPassword)
	if err != nil {
		return nil, err
	}

	auth.HostKeyCallback, err = knownhosts.New(sshKnownHostsPath)
	if err != nil {
		return nil, err
	}

	repo, err := gogit.PlainOpen(absolutePath)
	if err != nil {
		repo, err = gogit.PlainClone(absolutePath, false, &gogit.CloneOptions{
			URL:      url,
			Progress: nil,
			Auth:     auth,
		})
		if err != nil {
			return nil, fmt.Errorf("cannot open or clone the repository for URL %q and path %q: %w", url, absolutePath, err)
		}
	}

	return &Client{
		auth: auth,
		repo: repo,
	}, nil
}
