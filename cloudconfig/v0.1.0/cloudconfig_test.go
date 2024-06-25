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

package cloudconfig

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/cloudbase/garm-provider-common/defaults"
	"github.com/stretchr/testify/require"
)

// helper function
func getCloudInit() *CloudInit {
	return &CloudInit{
		PackageUpgrade:    true,
		Packages:          []string{"curl"},
		SSHAuthorizedKeys: []string{"test-ssh-key"},
		SystemInfo: &SystemInfo{
			DefaultUser: DefaultUser{
				Name:   defaults.DefaultUser,
				Home:   fmt.Sprintf("/home/%s", defaults.DefaultUser),
				Shell:  defaults.DefaultUserShell,
				Groups: defaults.DefaultUserGroups,
				Sudo:   "ALL=(ALL) NOPASSWD:ALL",
			},
		},
		WriteFiles: []File{
			{
				Encoding:    "b64",
				Content:     base64.StdEncoding.EncodeToString([]byte("test")),
				Owner:       "test-owner",
				Path:        "path",
				Permissions: "test-permissions",
			},
		},
	}
}

var (
	cloudInit = getCloudInit()
)

func TestNewDefaultCloudInitConfig(t *testing.T) {
	cloudInitCfg := NewDefaultCloudInitConfig()
	require.Equal(t, cloudInitCfg.PackageUpgrade, true)
	require.Equal(t, cloudInitCfg.Packages, []string{"curl", "tar"})
	require.Equal(t, cloudInitCfg.SystemInfo.DefaultUser.Name, defaults.DefaultUser)
	require.Equal(t, cloudInitCfg.SystemInfo.DefaultUser.Home, fmt.Sprintf("/home/%s", defaults.DefaultUser))
	require.Equal(t, cloudInitCfg.SystemInfo.DefaultUser.Shell, defaults.DefaultUserShell)
	require.Equal(t, cloudInitCfg.SystemInfo.DefaultUser.Groups, defaults.DefaultUserGroups)
	require.Equal(t, cloudInitCfg.SystemInfo.DefaultUser.Sudo, "ALL=(ALL) NOPASSWD:ALL")
}

func TestAddCACertNil(t *testing.T) {
	err := cloudInit.AddCACert(nil)
	require.NoError(t, err)
}

func TestAddCACertFailed(t *testing.T) {
	err := cloudInit.AddCACert([]byte("unknown-cert"))
	require.EqualError(t, err, "failed to parse CA cert bundle")
}

func TestAddSSHKey(t *testing.T) {
	cloudInit.AddSSHKey(cloudInit.SSHAuthorizedKeys...)
	require.Equal(t, []string{"test-ssh-key"}, cloudInit.SSHAuthorizedKeys)
}

func TestAddSSHKeyNotFound(t *testing.T) {
	cloudInit.AddSSHKey("new-test-ssh-key")
	require.Equal(t, []string{"test-ssh-key", "new-test-ssh-key"}, cloudInit.SSHAuthorizedKeys)
}

func TestAddPackage(t *testing.T) {
	cloudInit.AddPackage(cloudInit.Packages...)
	require.Equal(t, []string{"curl"}, cloudInit.Packages)
}

func TestAddPackageNotFound(t *testing.T) {
	cloudInit.AddPackage("tar")
	require.Equal(t, []string{"curl", "tar"}, cloudInit.Packages)
}

func TestAddRunCmd(t *testing.T) {
	cloudInit.AddRunCmd("test-run-cmd")
	require.Equal(t, []string{"test-run-cmd"}, cloudInit.RunCmd)
}

func TestAddFile(t *testing.T) {
	cloudInit.AddFile([]byte("test"), "test-path", "test-owner", "test-permissions")
	require.Equal(t, "b64", cloudInit.WriteFiles[0].Encoding)
	require.Equal(t, "dGVzdA==", cloudInit.WriteFiles[0].Content)
	require.Equal(t, "test-owner", cloudInit.WriteFiles[0].Owner)
	require.Equal(t, "test-permissions", cloudInit.WriteFiles[0].Permissions)
}

func TestAddFilePath(t *testing.T) {
	cloudInit.AddFile([]byte("content"), "path", "test-owner", "test-permissions")
}

func TestSerialize(t *testing.T) {
	serialized, err := cloudInit.Serialize()
	require.NoError(t, err)
	require.Contains(t, serialized, "#cloud-config")
}
