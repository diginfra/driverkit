#!/bin/bash
# SPDX-License-Identifier: Apache-2.0
#
# Copyright (C) 2023 The Diginfra Authors.
#
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
# Simple script that desperately tries to load the kernel instrumentation by
# looking for it in a bunch of ways. Convenient when running Diginfra inside
# a container or in other weird environments.
#
set -xeuo pipefail

export KBUILD_MODPOST_WARN=1

cd {{ .DriverBuildDir }}
mkdir -p build && cd build
{{ .CmakeCmd }}

{{ if .BuildModule }}
# Build the module
make CC=/usr/bin/gcc-{{ .GCCVersion }} driver
strip -g {{ .ModuleFullPath }}
# Print results
modinfo {{ .ModuleFullPath }}
{{ end }}

{{ if .BuildProbe }}
# Build the eBPF probe
make bpf
ls -l driver/bpf/probe.o
{{ end }}
