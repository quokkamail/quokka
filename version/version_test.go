// Copyright 2023 Quokka Contributors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package version_test

import (
	"testing"

	"github.com/quokkamail/quokka/version"
)

func TestVersion_String(t *testing.T) {
	tests := []struct {
		name    string
		version *version.Version
		want    string
	}{
		{
			name: "0.0.0",
			version: &version.Version{
				Major: 0,
				Minor: 0,
				Patch: 0,
			},
			want: "0.0.0",
		},
		{
			name: "1.0.0",
			version: &version.Version{
				Major: 1,
				Minor: 0,
				Patch: 0,
			},
			want: "1.0.0",
		},
		{
			name: "0.1.0",
			version: &version.Version{
				Major: 0,
				Minor: 1,
				Patch: 0,
			},
			want: "0.1.0",
		},
		{
			name: "0.0.1",
			version: &version.Version{
				Major: 0,
				Minor: 0,
				Patch: 1,
			},
			want: "0.0.1",
		},
		{
			name: "1.2.3",
			version: &version.Version{
				Major: 01,
				Minor: 02,
				Patch: 03,
			},
			want: "1.2.3",
		},
		{
			name: "10.20.30",
			version: &version.Version{
				Major: 10,
				Minor: 20,
				Patch: 30,
			},
			want: "10.20.30",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.version.String(); got != tt.want {
				t.Errorf("Version.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
