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
	tests := []struct {
		path       string
		pathExt    string
		executable bool
	}{
		{
			path:       "file.exe",
			pathExt:    ".exe",
			executable: true,
		},
		{
			path:       "file.bat",
			pathExt:    ".exe;.bat",
			executable: true,
		},
		{
			path:       "file.sh",
			pathExt:    ".exe;.bat",
			executable: false,
		},
		{
			path:       "file",
			pathExt:    ".exe;.bat",
			executable: false,
		},
	}

	for _, test := range tests {
		os.Setenv("PATHEXT", test.pathExt)
		executable := IsExecutable(test.path)
		require.Equal(t, test.executable, executable)
	}
}
