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

package builder

import (
	_ "embed"
	"fmt"

	"github.com/blang/semver"
	"github.com/diginfra/driverkit/pkg/kernelrelease"
)

//go:embed templates/centos_kernel.sh
var centosKernelTemplate string

//go:embed templates/centos.sh
var centosTemplate string

// TargetTypeCentos identifies the Centos target.
const TargetTypeCentos Type = "centos"

func init() {
	byTarget[TargetTypeCentos] = &centos{}
}

// centos is a driverkit target.
type centos struct {
}

type centosTemplateData struct {
	KernelDownloadURL string
}

func (c *centos) Name() string {
	return TargetTypeCentos.String()
}

func (c *centos) TemplateKernelUrlsScript() string { return centosKernelTemplate }

func (c *centos) TemplateScript() string {
	return centosTemplate
}

func (c *centos) URLs(kr kernelrelease.KernelRelease) ([]string, error) {
	vaultReleases := []string{
		"6.0/os",
		"6.0/updates",
		"6.1/os",
		"6.1/updates",
		"6.2/os",
		"6.2/updates",
		"6.3/os",
		"6.3/updates",
		"6.4/os",
		"6.4/updates",
		"6.5/os",
		"6.5/updates",
		"6.6/os",
		"6.6/updates",
		"6.7/os",
		"6.7/updates",
		"6.8/os",
		"6.8/updates",
		"6.9/os",
		"6.9/updates",
		"6.10/os",
		"6.10/updates",
		"7.0.1406/os",
		"7.0.1406/updates",
		"7.1.1503/os",
		"7.1.1503/updates",
		"7.2.1511/os",
		"7.2.1511/updates",
		"7.3.1611/os",
		"7.3.1611/updates",
		"7.4.1708/os",
		"7.4.1708/updates",
		"7.5.1804/os",
		"7.5.1804/updates",
		"7.6.1810/os",
		"7.6.1810/updates",
		"7.7.1908/os",
		"7.7.1908/updates",
		"7.8.2003/os",
		"7.8.2003/updates",
		"7.9.2009/os",
		"7.9.2009/updates",
		"8.0.1905/os",
		"8.0.1905/updates",
		"8.1.1911/os",
		"8.1.1911/updates",
	}

	centos8VaultReleases := []string{
		"8.0.1905/BaseOS",
		"8.1.1911/BaseOS",
		"8.2.2004/BaseOS",
		"8.3.2011/BaseOS",
		"8.4.2105/BaseOS",
		"8.5.2111/BaseOS",
	}

	edgeReleases := []string{
		"6/os",
		"6/updates",
		"7/os",
		"7/updates",
	}

	streamReleases := []string{
		"8/BaseOS",
		"8-stream/BaseOS",
	}

	stream9Releases := []string{
		"9-stream/AppStream",
		"9-stream/BaseOS",
	}

	urls := []string{}
	for _, r := range edgeReleases {
		urls = append(urls, fmt.Sprintf(
			"https://mirrors.edge.kernel.org/centos/%s/%s/Packages/kernel-devel-%s%s.rpm",
			r,
			kr.Architecture.ToNonDeb(),
			kr.Fullversion,
			kr.FullExtraversion,
		))
	}
	for _, r := range streamReleases {
		urls = append(urls, fmt.Sprintf(
			"https://mirrors.edge.kernel.org/centos/%s/%s/os/Packages/kernel-devel-%s%s.rpm",
			r,
			kr.Architecture.ToNonDeb(),
			kr.Fullversion,
			kr.FullExtraversion,
		))
	}
	for _, r := range vaultReleases {
		urls = append(urls, fmt.Sprintf(
			"http://vault.centos.org/%s/%s/Packages/kernel-devel-%s%s.rpm",
			r,
			kr.Architecture.ToNonDeb(),
			kr.Fullversion,
			kr.FullExtraversion,
		))
	}
	for _, r := range centos8VaultReleases {
		urls = append(urls, fmt.Sprintf(
			"http://vault.centos.org/%s/%s/os/Packages/kernel-devel-%s%s.rpm",
			r,
			kr.Architecture.ToNonDeb(),
			kr.Fullversion,
			kr.FullExtraversion,
		))
	}

	for _, r := range stream9Releases {
		urls = append(urls, fmt.Sprintf(
			"http://mirror.stream.centos.org/%s/%s/os/Packages/kernel-devel-%s%s.rpm",
			r,
			kr.Architecture.ToNonDeb(),
			kr.Fullversion,
			kr.FullExtraversion,
		))
	}

	return urls, nil
}

func (c *centos) KernelTemplateData(_ kernelrelease.KernelRelease, urls []string) interface{} {
	return centosTemplateData{
		KernelDownloadURL: urls[0],
	}
}

func (c *centos) GCCVersion(kr kernelrelease.KernelRelease) semver.Version {
	// 4.18+ centos 8 kernels need gcc 9
	if kr.Major == 4 && kr.Minor >= 18 {
		return semver.Version{Major: 9}
	}
	// 3.10.X kernels need 4.8.5 gcc version; see:
	// https://github.com/diginfra/driverkit/issues/236
	if kr.Major == 3 && kr.Minor == 10 {
		return semver.Version{Major: 4, Minor: 8, Patch: 5}
	}
	return semver.Version{}
}
