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

package execution

import (
	semver "github.com/Masterminds/semver/v3"
	executionv010 "github.com/cloudbase/garm-provider-common/execution/v0.1.0"
	executionv011 "github.com/cloudbase/garm-provider-common/execution/v0.1.1"
)

type Environment struct {
	EnvironmentV010  executionv010.EnvironmentV010
	EnvironmentV011  executionv011.EnvironmentV011
	InterfaceVersion semver.Version
}
