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
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExecSuccess(t *testing.T) {
	ctx := context.Background()
	providerBin := "gofmt"
	stdinData := []byte("")
	environ := []string{"TEST=1"}

	_, err := Exec(ctx, providerBin, stdinData, environ)
	require.NoError(t, err)
}

func TestExecFail(t *testing.T) {
	ctx := context.Background()
	providerBin := "garm-provider-test"
	stdinData := []byte("test")
	environ := []string{"TEST=1"}

	_, err := Exec(ctx, providerBin, stdinData, environ)
	require.ErrorContains(t, err, "provider binary failed with stdout")
}
