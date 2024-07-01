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

package params

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetArchitecture(t *testing.T) {
	architecture := "x86"
	tests := []struct {
		name string
		r    *RunnerApplicationDownload
		want string
	}{
		{
			name: "nil RunnerApplicationDownload",
			r:    nil,
			want: "",
		},
		{
			name: "nil architecture",
			r:    &RunnerApplicationDownload{},
			want: "",
		},
		{
			name: "architecture",
			r: &RunnerApplicationDownload{
				Architecture: &architecture,
			},
			want: architecture,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			arch := tt.r.GetArchitecture()
			assert.Equal(t, tt.want, arch)
		})
	}
}

func TestGetFilename(t *testing.T) {
	filename := "filename"
	tests := []struct {
		name string
		r    *RunnerApplicationDownload
		want string
	}{
		{
			name: "nil RunnerApplicationDownload",
			r:    nil,
			want: "",
		},
		{
			name: "nil filename",
			r:    &RunnerApplicationDownload{},
			want: "",
		},
		{
			name: "filename",
			r: &RunnerApplicationDownload{
				Filename: &filename,
			},
			want: filename,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filename := tt.r.GetFilename()
			assert.Equal(t, tt.want, filename)
		})
	}
}

func TestGetOS(t *testing.T) {
	os := "linux"
	tests := []struct {
		name string
		r    *RunnerApplicationDownload
		want string
	}{
		{
			name: "nil RunnerApplicationDownload",
			r:    nil,
			want: "",
		},
		{
			name: "nil os",
			r:    &RunnerApplicationDownload{},
			want: "",
		},
		{
			name: "os",
			r: &RunnerApplicationDownload{
				OS: &os,
			},
			want: os,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os := tt.r.GetOS()
			assert.Equal(t, tt.want, os)
		})
	}
}

func TestGetSHA256Checksum(t *testing.T) {
	checksum := "test-checksum"
	tests := []struct {
		name string
		r    *RunnerApplicationDownload
		want string
	}{
		{
			name: "nil RunnerApplicationDownload",
			r:    nil,
			want: "",
		},
		{
			name: "nil checksum",
			r:    &RunnerApplicationDownload{},
			want: "",
		},
		{
			name: "checksum",
			r: &RunnerApplicationDownload{
				SHA256Checksum: &checksum,
			},
			want: checksum,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checksum := tt.r.GetSHA256Checksum()
			assert.Equal(t, tt.want, checksum)
		})
	}
}

func TestGetTempDownloadToken(t *testing.T) {
	token := "test-token"
	tests := []struct {
		name string
		r    *RunnerApplicationDownload
		want string
	}{
		{
			name: "nil RunnerApplicationDownload",
			r:    nil,
			want: "",
		},
		{
			name: "nil token",
			r:    &RunnerApplicationDownload{},
			want: "",
		},
		{
			name: "token",
			r: &RunnerApplicationDownload{
				TempDownloadToken: &token,
			},
			want: token,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token := tt.r.GetTempDownloadToken()
			assert.Equal(t, tt.want, token)
		})
	}
}
