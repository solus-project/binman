//
// Copyright © 2016 Ikey Doherty <ikey@solus-project.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//

// Package libeopkg provides Go-native access to `.eopkg` files, allowing
// binman to read and manipulate them without having a host-side eopkg
// tool.
//
// It should also be noted that `eopkg` is implemented in Python, so calling
// out to the host-side tool just isn't acceptable for the performance we
// require.
// In time, `sol` will replace eopkg and it is very likely that we'll base
// the new `libsol` component on the C library using cgo.
package libeopkg

// Package represents a binary .eopkg file
type Package struct {
	Path string // Path to this .eopkg file
}