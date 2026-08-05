// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package auto

import (
	"errors"
	"fmt"
	"regexp"

	git "github.com/go-git/go-git/v5"

	"github.com/siderolabs/kres/internal/project/common"
)

// DetectGit detects if current directory is git repository.
func (builder *builder) DetectGit() (bool, error) {
	repo, err := git.PlainOpenWithOptions(".", &git.PlainOpenOptions{
		EnableDotGitCommonDir: true,
	})
	if err != nil {
		// not a git repo, ignore
		return false, nil //nolint:nilerr
	}

	c, err := repo.Config()
	if err != nil {
		return true, fmt.Errorf("failed to get repository configuration: %w", err)
	}

	rawConfig := c.Raw

	const (
		main              = "main"
		branchSectionName = "branch"
	)

	if !rawConfig.HasSection(branchSectionName) {
		return true, fmt.Errorf("repository configuration section %q not found", branchSectionName)
	}

	branchSection := rawConfig.Section(branchSectionName)

	for _, b := range branchSection.Subsections {
		if b.Name == main {
			builder.meta.MainBranch = main

			break
		}

		remote := b.Option("remote")
		if remote == git.DefaultRemoteName {
			builder.meta.MainBranch = b.Name
		}
	}

	if builder.meta.MainBranch == "" {
		builder.meta.MainBranch = main
	}

	if head, headErr := repo.Head(); headErr == nil && head.Name().IsBranch() {
		builder.meta.CurrentBranch = head.Name().Short()
	}

	remotes, err := repo.Remotes()
	if err != nil {
		return true, err
	}

	// The repository's origin is the durable authority for generated project
	// settings. A fork commonly retains `upstream` for synchronization, but
	// using it here would make `make rekres` regenerate workflow/image settings
	// for somebody else's organization. Fall back to upstream only when origin
	// is genuinely absent.
	var sourceRemote *git.Remote

	for _, remote := range remotes {
		if remote.Config().Name == git.DefaultRemoteName {
			sourceRemote = remote

			break
		}
	}

	if sourceRemote == nil {
		for _, remote := range remotes {
			if remote.Config().Name == "upstream" {
				sourceRemote = remote

				break
			}
		}
	}

	if sourceRemote == nil {
		return true, errors.New("neither 'origin' nor 'upstream' remote found")
	}

	remoteURLregexp := `((?:[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,})[:/]+([^/:]+)/([^/]+)\.git$`
	for _, remoteURL := range sourceRemote.Config().URLs {
		matches := regexp.MustCompile(remoteURLregexp).FindStringSubmatch(remoteURL)
		if len(matches) == 4 {
			if matches[1] != "github.com" {
				return false, nil //nolint:nilerr
			}

			builder.meta.GitHubOrganization = matches[2]
			builder.meta.GitHubRepository = matches[3]

			return true, nil
		}
	}

	return true, fmt.Errorf("failed to parse remote URL: %s", sourceRemote)
}

// BuildGit builds steps for Git repository.
func (builder *builder) BuildGit() error {
	builder.commonInputs = append(
		builder.commonInputs,
		common.NewRepository(builder.meta),
		common.NewCheckDirty(builder.meta),
	)

	return nil
}
