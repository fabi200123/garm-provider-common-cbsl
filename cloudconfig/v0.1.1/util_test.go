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
	"fmt"
	"testing"

	params "github.com/cloudbase/garm-provider-common/params/v0.1.1"
	"github.com/stretchr/testify/require"
)

var (
	fileName          = "test-filename"
	downloadURL       = "https://example.com/test.zip"
	tempDownloadToken = "test-token"
	extraSpecsJson    = `{
        "runner_install_template": "dGVzdF90ZW1wbGF0ZToge3sgLlJ1bm5lck5hbWUgfX0gLSB7eyAuRXh0cmFDb250ZXh0LnRlc3RfdmFyMSB9fSAtIHt7IC5FeHRyYUNvbnRleHQudGVzdF92YXIyIH19",
        "extra_context": {
            "test_var1": "bogus-value1",
            "test_var2": "bogus-value2"
        },
        "pre_install_scripts": {
            "test-script": "dGVzdC1zY3JpcHQtY29udGVudA=="
        }
    }`

	newCloudCfg     = NewDefaultCloudInitConfig()
	bootstrapParams = params.BootstrapInstance{
		ExtraSpecs: []byte(extraSpecsJson),
		UserDataOptions: params.UserDataOptions{
			DisableUpdatesOnBoot: newCloudCfg.PackageUpgrade,
			ExtraPackages:        newCloudCfg.Packages,
		},
	}
	tools = params.RunnerApplicationDownload{
		Filename:          &fileName,
		DownloadURL:       &downloadURL,
		TempDownloadToken: &tempDownloadToken,
	}
)

func TestGetSpecs(t *testing.T) {
	expectedSpecs := CloudConfigSpec{
		RunnerInstallTemplate: []byte("test_template: {{ .RunnerName }} - {{ .ExtraContext.test_var1 }} - {{ .ExtraContext.test_var2 }}"),
		ExtraContext: map[string]string{
			"test_var1": "bogus-value1",
			"test_var2": "bogus-value2",
		},
		PreInstallScripts: map[string][]byte{
			"test-script": []byte("test-script-content"),
		},
	}

	specs, err := GetSpecs(bootstrapParams)
	require.NoError(t, err)
	require.Equal(t, expectedSpecs, specs)
}

func TestGetSpecsExtraSpecsNil(t *testing.T) {
	bootstrapParams := params.BootstrapInstance{
		ExtraSpecs: []byte(`{"runner_install_template": "dGVzdA=="}`),
	}

	specs, err := GetSpecs(bootstrapParams)
	require.NoError(t, err)
	require.Equal(t, "test", string(specs.RunnerInstallTemplate))
}

func TestGetSpecsUnmarshalFailed(t *testing.T) {
	bootstrapParams := params.BootstrapInstance{
		ExtraSpecs: []byte("invalid-json"),
	}

	_, err := GetSpecs(bootstrapParams)
	require.Error(t, err)
	require.Contains(t, err.Error(), "unmarshaling extra specs: invalid character")
}

func TestGetRunnerInstallScript(t *testing.T) {
	script, err := GetRunnerInstallScript(bootstrapParams, tools, "test-runner-name")
	require.NoError(t, err)
	require.Equal(t, "test_template: test-runner-name - bogus-value1 - bogus-value2", string(script))
}

func TestGetRunnerInstallScriptMissingFilename(t *testing.T) {
	_, err := GetRunnerInstallScript(bootstrapParams, params.RunnerApplicationDownload{}, "test-runner-name")
	require.Error(t, err)
	require.EqualError(t, err, "missing tools filename")
}

func TestGetRunnerInstallScriptMissingURL(t *testing.T) {
	tools := params.RunnerApplicationDownload{
		Filename: &fileName,
	}

	_, err := GetRunnerInstallScript(params.BootstrapInstance{}, tools, "test-runner-name")
	require.Error(t, err)
	require.EqualError(t, err, "missing tools download URL")
}

func TestGetRunnerInstallScriptGettingSpecsFailed(t *testing.T) {
	bootstrapParams := params.BootstrapInstance{
		ExtraSpecs: []byte("invalid-json"),
	}

	_, err := GetRunnerInstallScript(bootstrapParams, tools, "test-runner-name")
	require.Error(t, err)
	require.Contains(t, err.Error(), "getting specs: unmarshaling extra specs: invalid character")
}

func TestGetRunnerInstallScriptFailed(t *testing.T) {
	_, err := GetRunnerInstallScript(params.BootstrapInstance{}, tools, "test-runner-name")
	require.Error(t, err)
	require.Contains(t, err.Error(), "generating script: unsupported os type: ")
}

func TestGetCloudInitConfig(t *testing.T) {
	cloudInitCfg, err := GetCloudInitConfig(bootstrapParams, []byte("test-install-script"))
	require.NoError(t, err)
	require.Contains(t, cloudInitCfg, `#cloud-config`)
}

func TestGetCloudInitConfigGetSpecsFailed(t *testing.T) {
	bootstrapParams = params.BootstrapInstance{
		ExtraSpecs: []byte("invalid-json"),
	}

	_, err := GetCloudInitConfig(bootstrapParams, []byte("test-install-script"))
	require.Error(t, err)
	require.Contains(t, err.Error(), "getting specs: unmarshaling extra specs: invalid character")
}

func TestGetCloudInitConfigAddCACertBundleFailed(t *testing.T) {
	bootstrapParams = params.BootstrapInstance{
		CACertBundle: []byte(`dummy-ca-cert-bundle`),
	}

	_, err := GetCloudInitConfig(bootstrapParams, []byte("test-install-script"))
	require.Error(t, err)
	require.Contains(t, err.Error(), "adding CA cert bundle: failed to parse CA cert bundle")
}

func TestGetCloudConfigForLinux(t *testing.T) {
	bootstrapParams = params.BootstrapInstance{
		OSType: "linux",
		UserDataOptions: params.UserDataOptions{
			DisableUpdatesOnBoot: newCloudCfg.PackageUpgrade,
			ExtraPackages:        newCloudCfg.Packages,
		},
	}

	cloudCfg, err := GetCloudConfig(bootstrapParams, tools, "test-runner-name")
	require.NoError(t, err)
	require.Contains(t, cloudCfg, `#cloud-config`)
}

func TestGetCloudConfigForWindows(t *testing.T) {
	bootstrapParams = params.BootstrapInstance{
		OSType: "windows",
		UserDataOptions: params.UserDataOptions{
			DisableUpdatesOnBoot: newCloudCfg.PackageUpgrade,
			ExtraPackages:        newCloudCfg.Packages,
		},
	}

	cloudCfg, err := GetCloudConfig(bootstrapParams, tools, "test-runner-name")
	require.NoError(t, err)
	require.Contains(t, cloudCfg, `#ps1_sysnative`)
}

func TestGetCloudConfigGeneratingScriptFailed(t *testing.T) {
	_, err := GetCloudConfig(bootstrapParams, params.RunnerApplicationDownload{}, "test-runner-name")
	require.Error(t, err)
	require.Contains(t, err.Error(), "generating script: missing tools filename")
}

func TestGetCloudConfigInitFailed(t *testing.T) {
	bootstrapParams = params.BootstrapInstance{
		CACertBundle: []byte(`dummy-ca-cert-bundle`),
		OSType:       "linux",
	}

	_, err := GetCloudConfig(bootstrapParams, tools, "test-runner-name")
	require.Error(t, err)
	require.Contains(t, err.Error(), "getting cloud init config: adding CA cert bundle: failed to parse CA cert bundle")
}

func TestGetCloudConfigUnknownOSType(t *testing.T) {
	bootstrapParams = params.BootstrapInstance{
		ExtraSpecs: []byte(extraSpecsJson),
		OSType:     "dummy-os-type",
	}

	_, err := GetCloudConfig(bootstrapParams, tools, "test-runner-name")
	require.Error(t, err)
	require.EqualError(t, err, fmt.Sprintf("unknown os type: %s", bootstrapParams.OSType))
}
