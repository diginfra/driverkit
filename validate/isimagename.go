// SPDX-License-Identifier: Apache-2.0
/*
Copyright (C) 2023 The Diginfra Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package validate

import (
	"strings"
	"unicode"

	"github.com/go-playground/validator/v10"
)

const letters = "abcdefghijklmnopqrstuvwxyz"
const digits = "0123456789"
const separators = "/.-@_:"
const alphabet = letters + digits + separators

func isImageName(fl validator.FieldLevel) bool {
	name := fl.Field().String()

	for _, c := range name {
		if !strings.ContainsRune(alphabet, unicode.ToLower(c)) {
			return false
		}
	}

	for _, component := range strings.Split(name, "/") {
		// a component may not be empty (i.e. double slashes are not allowed)
		if len(component) == 0 {
			return false
		}

		// a component may not start or end with a separator
		if strings.Contains(separators, component[0:1]) || strings.Contains(separators, component[len(component)-1:]) {
			return false
		}
	}

	return true
}
