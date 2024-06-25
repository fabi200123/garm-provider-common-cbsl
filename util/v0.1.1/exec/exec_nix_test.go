//go:build !windows
// +build !windows

// Copyright 2023 Cloudbase Solutions SRL
//
//    Licensed under the Apache License, Version 2.0 (the "License"); you may
//    not use this file except in compliance with the License. You may obtain
//    a copy of the License at
//
//         http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//    WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//    License for the specific language governing permissions and limitations
//    under the License.

package exec

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsExecutable(t *testing.T) {
	// Create a temporary file with executable permission.
	err := os.WriteFile("file.sh", []byte(`!#/bin/sh`), 0o755)
	if err != nil {
		t.Fatalf("failed to create temporary file: %s", err)
	}
	// Remove the temporary file when the test finishes.
	defer os.Remove("file.sh")

	// Test that the function returns true for an executable file.
	require.True(t, IsExecutable("file.sh"))

	// Create a temporary file without executable permission.
	err = os.WriteFile("file2.txt", []byte(`!#/bin/sh`), 0o644)
	if err != nil {
		t.Fatalf("failed to create temporary file: %s", err)
	}
	// Remove the temporary file when the test finishes.
	defer os.Remove("file2.txt")

	// Test that the function returns false for a non-executable file.
	require.False(t, IsExecutable("file2.txt"))
}
