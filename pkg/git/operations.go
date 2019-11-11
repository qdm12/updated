package git

import (
	"fmt"
	"time"

	gogit "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
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
	err = workTree.Pull(&gogit.PullOptions{
		RemoteName: "origin",
		Auth:       c.auth,
	})
	if err != nil && err.Error() != "already up-to-date" {
		return fmt.Errorf("cannot pull: %w", err)
	}
	return nil
}

func (c *Client) Status() (string, error) {
	workTree, err := c.repo.Worktree()
	if err != nil {
		return "", err
	}
	status, err := workTree.Status()
	return status.String(), err
}

func (c *Client) IsClean() (bool, error) {
	workTree, err := c.repo.Worktree()
	if err != nil {
		return false, err
	}
	status, err := workTree.Status()
	return status.IsClean(), err
}

func (c *Client) Add(filename string) error {
	workTree, err := c.repo.Worktree()
	if err != nil {
		return err
	}
	_, err = workTree.Add(filename)
	return err
}

func (c *Client) Commit(message string) error {
	workTree, err := c.repo.Worktree()
	if err != nil {
		return err
	}
	_, err = workTree.Commit(message, &gogit.CommitOptions{
		Author: &object.Signature{
			Name: "updated",
			When: time.Now(),
		}})
	return err
}

func (c *Client) Push() error {
	return c.repo.Push(&gogit.PushOptions{
		Auth:     c.auth,
		Progress: nil,
	})
}

func (c *Client) UploadAllChanges(message string) (err error) {
	err = c.Add(".")
	if err != nil {
		return err
	}
	err = c.Commit(message)
	if err != nil {
		return err
	}
	err = c.Push()
	if err != nil {
		return err
	}
	return nil
}
