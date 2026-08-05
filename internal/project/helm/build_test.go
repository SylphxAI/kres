// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package helm_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/siderolabs/kres/internal/config"
	"github.com/siderolabs/kres/internal/output/makefile"
	"github.com/siderolabs/kres/internal/project/helm"
	"github.com/siderolabs/kres/internal/project/meta"
)

func TestBuildInterfaces(t *testing.T) {
	assert.Implements(t, (*makefile.Compiler)(nil), new(helm.Build))
}

func TestCompileMakefileInstallsPinnedHelmPluginsFailClosed(t *testing.T) {
	build := helm.NewBuild(&meta.Options{})
	output := makefile.NewOutput()

	require.NoError(t, output.Compile(build))

	var buf bytes.Buffer
	require.NoError(t, output.GenerateFile("Makefile", &buf))

	generated := buf.String()
	assert.Contains(t, generated, "helm plugin uninstall unittest")
	assert.Contains(t, generated, "helm plugin uninstall schema")
	assert.Contains(t, generated, "--version="+config.HelmUnitTestVersion)
	assert.Contains(t, generated, "--version="+config.HelmValuesSchemaJSONVersion)
	assert.NotContains(t, generated, "--verify=false")
	assert.NotContains(t, generated, "-helm plugin install")
}
