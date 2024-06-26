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

package cmd

import flag "github.com/spf13/pflag"

var kubernetesOptions = &KubeOptions{}

type KubeOptions struct {
	RunAsUser       int64  `json:"runAsUser,omitempty" protobuf:"varint,2,opt,name=runAsUser" default:"0"`
	Namespace       string `validate:"required" name:"namespace" default:"default"`
	ImagePullSecret string `validate:"omitempty" name:"image-pull-secret" default:""`
}

func addKubernetesFlags(flags *flag.FlagSet) {
	flags.StringVarP(&kubernetesOptions.Namespace, "namespace", "n", "default", "If present, the namespace scope for the pods and its config ")
	flags.Int64Var(&kubernetesOptions.RunAsUser, "run-as-user", 0, "Pods runner user")
	flags.StringVar(&kubernetesOptions.ImagePullSecret, "image-pull-secret", "", "ImagePullSecret")
}
