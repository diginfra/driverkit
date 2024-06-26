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

# Fetch the kernel
mkdir /tmp/kernel-download
cd /tmp/kernel-download
{{ range $url := .KernelDownloadURLS }}
curl --silent -o kernel.deb -SL {{ $url }}
ar x kernel.deb
tar -xf data.tar.xz
{{ end }}

cd usr/src/
sourcedir=$(find . -type d -name "{{ .KernelHeadersPattern }}" | head -n 1 | xargs readlink -f)

# Patch makefile to avoid using absolute `/usr/src` path; instead use `..` relative one.
sed -i 's/\/usr\/src/../g' $sourcedir/Makefile

# exit value
export KERNELDIR=$sourcedir
