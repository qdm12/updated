package git

import (
	"context"
	"errors"
	"fmt"
	"time"

	gogit "gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

var (
	ErrHead         = errors.New("cannot get the HEAD")
	ErrSetReference = errors.New("cannot set reference")
	ErrWorkTree     = errors.New("cannot get work tree")
	ErrCheckout     = errors.New("cannot checkout")
	ErrStatus       = errors.New("cannot get status")
	ErrAdd          = errors.New("cannot add file")
	ErrCommit       = errors.New("cannot commit")
	ErrPush         = errors.New("cannot push")
)

// Branch creates a new branch from the current head.
func (c *Client) Branch(branchName string) (err error) {
	headRef, err := c.repo.Head()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrHead, err)
	}

	refName := plumbing.NewBranchReferenceName(branchName)
	ref := plumbing.NewHashReference(refName, headRef.Hash())
	err = c.repo.Storer.SetReference(ref)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrSetReference, err)
	}

	return nil
}

// CheckoutBranch force checkout to an existing branch.
func (c *Client) CheckoutBranch(branchName string) (err error) {
	workTree, err := c.repo.Worktree()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrWorkTree, err)
	}

	err = workTree.Checkout(&gogit.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName(branchName),
		Force:  true,
	})
	if err != nil {
		return fmt.Errorf("for branch: %s: %w", branchName, err)
	}

	return nil
}

// Pull pulls changes from the repository to the local directory.
// It does not support merge conflicts and will return an error in this case.
func (c *Client) Pull(ctx context.Context) (err error) {
	workTree, err := c.repo.Worktree()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrWorkTree, err)
	}

	options := &gogit.PullOptions{
		RemoteName: "origin",
		Auth:       c.auth,
	}

	err = workTree.PullContext(ctx, options)
	if err != nil && !errors.Is(err, gogit.NoErrAlreadyUpToDate) {
		return fmt.Errorf("cannot pull from %s: %w",
			options.RemoteName, err)
	}

	return nil
}

// Status returns the git status of the working tree.
func (c *Client) Status() (statusString string, err error) {
	workTree, err := c.repo.Worktree()
	if err != nil {
		return "", fmt.Errorf("%w: %w", ErrWorkTree, err)
	}

	status, err := workTree.Status()
	if err != nil {
		return "", err
	}

	return status.String(), nil
}

// IsClean checks if the working tree is clean (no uncommitted changes).
func (c *Client) IsClean() (clean bool, err error) {
	workTree, err := c.repo.Worktree()
	if err != nil {
		return false, fmt.Errorf("%w: %w", ErrWorkTree, err)
	}

	status, err := workTree.Status()
	if err != nil {
		return false, fmt.Errorf("%w: %w", ErrStatus, err)
	}

	return status.IsClean(), err
}

// Add stages a file or files.
func (c *Client) Add(filename string) (err error) {
	workTree, err := c.repo.Worktree()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrWorkTree, err)
	}

	_, err = workTree.Add(filename)
	if err != nil {
		return err
	}

	return nil
}

// Commit commits changes to the local repository with the given message.
func (c *Client) Commit(message string) (err error) {
	workTree, err := c.repo.Worktree()
	if err != nil {
		return fmt.Errorf("%w: %w", ErrWorkTree, err)
	}

	options := &gogit.CommitOptions{
		Author: &object.Signature{
			Name: "updated",
			When: time.Now(),
		},
	}

	_, err = workTree.Commit(message, options)
	if err != nil {
		return err
	}

	return nil
}

// Push pushes local changes to the remote repository.
func (c *Client) Push(ctx context.Context) (err error) {
	options := &gogit.PushOptions{
		Auth:     c.auth,
		Progress: nil,
	}

	return c.repo.PushContext(ctx, options)
}

// UploadAllChanges adds all changes, commits them with the given message and pushes them.
func (c *Client) UploadAllChanges(ctx context.Context,
	message string,
) (err error) {
	err = c.Add(".")
	if err != nil {
		return fmt.Errorf("%w: %w", ErrAdd, err)
	}

	err = c.Commit(message)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrCommit, err)
	}

	err = c.Push(ctx)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrPush, err)
	}

	return nil
}
