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

package smtp

import "time"

const (
	DefaultDataBlockTimeout         = 3 * time.Minute
	DefaultDataInitiationTimeout    = 2 * time.Minute
	DefaultDataTerminationTimeout   = 10 * time.Minute
	DefaultInitial220MessageTimeout = 5 * time.Minute
	DefaultMAILCommandTimeout       = 5 * time.Minute
	DefaultRCPTCommandTimeout       = 5 * time.Minute
	DefaultServerTimeout            = 5 * time.Minute
)
