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
	"fmt"
	"reflect"

	"github.com/diginfra/driverkit/pkg/kernelrelease"
	"github.com/go-playground/validator/v10"
)

func isArchitectureSupported(fl validator.FieldLevel) bool {
	field := fl.Field()

	switch field.Kind() {
	case reflect.String:
		for arch := range kernelrelease.SupportedArchs {
			if arch.String() == field.String() {
				return true
			}
		}
		return false
	}

	panic(fmt.Sprintf("Bad field type %T", field.Interface()))
}
