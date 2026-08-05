// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package common_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/siderolabs/kres/internal/config"
	"github.com/siderolabs/kres/internal/output/makefile"
	"github.com/siderolabs/kres/internal/project/common"
	"github.com/siderolabs/kres/internal/project/meta"
)

func TestDockerInterfaces(t *testing.T) {
	assert.Implements(t, (*makefile.Compiler)(nil), new(common.Docker))
}

func TestDockerUsesLowercaseOrganizationForOCIRegistry(t *testing.T) {
	docker := common.NewDocker(&meta.Options{
		ContainerImageFrontend: config.ContainerImageFrontendDockerfile,
		GitHubOrganization:     "SylphxAI",
	})
	output := makefile.NewOutput()

	require.NoError(t, docker.CompileMakefile(output))

	var rendered bytes.Buffer

	require.NoError(t, output.GenerateFile("Makefile", &rendered))
	assert.Contains(t, rendered.String(), "USERNAME ?= sylphxai\n")
}
