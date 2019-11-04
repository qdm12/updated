package git

import (
	"fmt"

	gogit "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

// Branch creates a new branch from the current head
func (c *Client) Branch(branchName string) error {
	headRef, err := c.repo.Head()
	if err != nil {
		return fmt.Errorf("cannot branch: %w", err)
	}
	refName := plumbing.NewBranchReferenceName(branchName)
	ref := plumbing.NewHashReference(refName, headRef.Hash())
	err = c.repo.Storer.SetReference(ref)
	return err
}

// CheckoutBranch force checkout to an existing branch
func (c *Client) CheckoutBranch(branchName string) error {
	workTree, err := c.repo.Worktree()
	if err != nil {
		return fmt.Errorf("cannot checkout to branch %q: %w", branchName, err)
	}
	err = workTree.Checkout(&gogit.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branchName),
		Force:  true,
	})
	return err
}

// Pull pulls changes from the repository to the local directory.
// It does not support merge conflicts and will return an error in this case.
func (c *Client) Pull() error {
	workTree, err := c.repo.Worktree()
	if err != nil {
		return fmt.Errorf("cannot pull: %w", err)
	}
	err = workTree.Pull(&gogit.PullOptions{RemoteName: "origin"})
	if err != nil {
		return fmt.Errorf("cannot pull: %w", err)
	}
	return nil
}
