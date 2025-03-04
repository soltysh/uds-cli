// Copyright 2024 Defense Unicorns
// SPDX-License-Identifier: AGPL-3.0-or-later OR LicenseRef-Defense-Unicorns-Commercial

// Package test provides e2e tests for UDS.
package test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/zarf-dev/zarf/src/pkg/utils/exec"
)

// NOTE: These tests test that the embedded `zarf` commands are imported properly and function as expected

// TestZarfLint tests to ensure that the `zarf dev lint` command functions (which requires the zarf schema to be embedded in main.go)
func TestZarfLint(t *testing.T) {
	_, stderr := runCmd(t, "zarf dev lint src/test/packages/podinfo")
	require.Contains(t, stderr, "Image not pinned with digest - ghcr.io/stefanprodan/podinfo:6.4.0")
}

func TestZarfToolsVersions(t *testing.T) {
	type args struct {
		tool     string
		toolRepo string
	}
	tests := []struct {
		name        string
		description string
		args        args
	}{
		{
			name:        "HelmVersion",
			description: "zarf tools helm version",
			args:        args{tool: "helm", toolRepo: "helm.sh/helm/v3"},
		},
		{
			name:        "CraneVersion",
			description: "zarf tools crane version",
			args:        args{tool: "crane", toolRepo: "github.com/google/go-containerregistry"},
		},
		{
			name:        "SyftVersion",
			description: "zarf tools syft version",
			args:        args{tool: "syft", toolRepo: "github.com/anchore/syft"},
		},
		{
			name:        "ArchiverVersion",
			description: "zarf tools archiver version",
			args:        args{tool: "archiver", toolRepo: "github.com/mholt/archiver/v3"},
		},
	}

	for _, tt := range tests {
		res, stderr := runCmd(t, fmt.Sprintf("zarf tools %s version", tt.args.tool))

		toolRepoVerArgs := fmt.Sprintf("list -f '{{.Version}}' -m %s", tt.args.toolRepo)
		verRes, _, verErr := exec.Cmd("go", strings.Split(toolRepoVerArgs, " ")...)
		require.NoError(t, verErr)

		toolVersion := strings.Split(verRes, "'")[1]
		output := res
		if res == "" {
			output = stderr
		}
		require.Contains(t, output, toolVersion)
	}
}
