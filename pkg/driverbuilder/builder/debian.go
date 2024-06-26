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
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"github.com/diginfra/driverkit/pkg/kernelrelease"
)

//go:embed templates/debian_kernel.sh
var debianKernelTemplate string

//go:embed templates/debian.sh
var debianTemplate string

// TargetTypeDebian identifies the Debian target.
const TargetTypeDebian Type = "debian"

// We need:
// kernel devel
// kernel devel common
// kbuild package
const debianRequiredURLs = 3

func init() {
	byTarget[TargetTypeDebian] = &debian{}
}

type debianTemplateData struct {
	KernelDownloadURLS   []string
	KernelLocalVersion   string
	KernelHeadersPattern string
}

// debian is a driverkit target.
type debian struct {
}

func (v *debian) Name() string {
	return TargetTypeDebian.String()
}

func (v *debian) TemplateKernelUrlsScript() string { return debianKernelTemplate }

func (v *debian) TemplateScript() string {
	return debianTemplate
}

func (v *debian) URLs(kr kernelrelease.KernelRelease) ([]string, error) {
	return fetchDebianKernelURLs(kr)
}

func (v *debian) KernelTemplateData(kr kernelrelease.KernelRelease, urls []string) interface{} {
	var KernelHeadersPattern string
	if strings.HasSuffix(kr.Extraversion, "pve") {
		KernelHeadersPattern = "linux-headers-*pve"
	} else if strings.Contains(kr.FullExtraversion, "rpi") {
		KernelHeadersPattern = "linux-headers-*-rpi-v*"
	} else {
		KernelHeadersPattern = "linux-headers-*" + kr.Architecture.String()
	}

	return debianTemplateData{
		KernelDownloadURLS:   urls,
		KernelLocalVersion:   kr.FullExtraversion,
		KernelHeadersPattern: KernelHeadersPattern,
	}
}

func (v *debian) MinimumURLs() int {
	return debianRequiredURLs
}

func fetchDebianKernelURLs(kr kernelrelease.KernelRelease) ([]string, error) {
	kbuildURL, err := debianKbuildURLFromRelease(kr)
	if err != nil {
		return nil, err
	}

	urls, err := debianHeadersURLFromRelease(kr)
	if err != nil {
		return nil, err
	}
	urls = append(urls, kbuildURL)

	return urls, nil
}

func debianHeadersURLFromRelease(kr kernelrelease.KernelRelease) ([]string, error) {
	baseURLS := []string{
		"http://security-cdn.debian.org/pool/main/l/linux/",
		"http://security-cdn.debian.org/pool/updates/main/l/linux/",
		"https://mirrors.edge.kernel.org/debian/pool/main/l/linux/",
	}

	for _, u := range baseURLS {
		urls, err := fetchDebianHeadersURLFromRelease(u, kr)

		if err == nil {
			return urls, err
		}
	}

	return nil, HeadersNotFoundErr
}

func fetchDebianHeadersURLFromRelease(baseURL string, kr kernelrelease.KernelRelease) ([]string, error) {
	extraVersionPartial := strings.TrimSuffix(kr.FullExtraversion, "-"+kr.Architecture.String())
	matchExtraGroup := kr.Architecture.String()
	rmatch := `href="(linux-headers-%d\.%d\.%d%s-(%s)_.*(%s|all)\.deb)"`

	// For urls like: http://security.debian.org/pool/updates/main/l/linux/linux-headers-5.10.0-12-amd64_5.10.103-1_amd64.deb
	// when 5.10.103-1 is passed as kernel version
	rmatchNew := `href="(linux-headers-[0-9]+\.[0-9]+\.[0-9]+-[0-9]+-(%s)_%d\.%d\.%d%s_(%s|all)\.deb)"`

	matchExtraGroupCommon := "common"

	// match for kernel versions like 4.19.0-6-cloud-amd64
	supportedExtraFlavors := []string{"cloud", "rt", "rpi"}
	for _, supportedExtraFlavor := range supportedExtraFlavors {
		if strings.Contains(kr.FullExtraversion, "-"+supportedExtraFlavor) {
			extraVersionPartial = strings.TrimSuffix(extraVersionPartial, "-"+supportedExtraFlavor)
			matchExtraGroup = supportedExtraFlavor + "-" + matchExtraGroup

			// rpi and rt have a different common package, named `common-{rt,rpi}`
			if supportedExtraFlavor == "rt" || supportedExtraFlavor == "rpi" {
				matchExtraGroupCommon += "-" + supportedExtraFlavor
			}
			break
		}
	}

	// download index
	resp, err := http.Get(baseURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	bodyStr := string(body)

	// look for kernel headers
	fullregex := fmt.Sprintf(rmatch, kr.Major, kr.Minor, kr.Patch,
		extraVersionPartial, matchExtraGroup, kr.Architecture.String())
	pattern := regexp.MustCompile(fullregex)
	matches := pattern.FindStringSubmatch(bodyStr)
	if len(matches) < 1 {
		fullregex = fmt.Sprintf(rmatchNew, matchExtraGroup, kr.Major, kr.Minor, kr.Patch,
			extraVersionPartial, kr.Architecture.String())
		pattern = regexp.MustCompile(fullregex)
		matches = pattern.FindStringSubmatch(bodyStr)
		if len(matches) < 1 {
			return nil, fmt.Errorf("kernel headers not found")
		}
	}

	// look for kernel headers common
	fullregexCommon := fmt.Sprintf(rmatch, kr.Major, kr.Minor, kr.Patch,
		extraVersionPartial, matchExtraGroupCommon, kr.Architecture.String())
	patternCommon := regexp.MustCompile(fullregexCommon)
	matchesCommon := patternCommon.FindStringSubmatch(bodyStr)
	if len(matchesCommon) < 1 {
		fullregexCommon = fmt.Sprintf(rmatchNew, matchExtraGroupCommon, kr.Major, kr.Minor, kr.Patch,
			extraVersionPartial, kr.Architecture.String())
		patternCommon = regexp.MustCompile(fullregexCommon)
		matchesCommon = patternCommon.FindStringSubmatch(bodyStr)
		if len(matchesCommon) < 1 {
			return nil, fmt.Errorf("kernel headers common not found")
		}
	}

	foundURLs := []string{fmt.Sprintf("%s%s", baseURL, matches[1])}
	foundURLs = append(foundURLs, fmt.Sprintf("%s%s", baseURL, matchesCommon[1]))

	return foundURLs, nil
}

func debianKbuildURLFromRelease(kr kernelrelease.KernelRelease) (string, error) {
	rmatch := `href="(linux-kbuild-%d\.%d.*%s\.deb)"`

	kbuildPattern := regexp.MustCompile(fmt.Sprintf(rmatch, kr.Major, kr.Minor, kr.Architecture.String()))
	baseURL := "http://mirrors.kernel.org/debian/pool/main/l/linux/"
	if kr.Major == 3 {
		baseURL = "http://mirrors.kernel.org/debian/pool/main/l/linux-tools/"
	}

	resp, err := http.Get(baseURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	match := kbuildPattern.FindStringSubmatch(string(body))

	if len(match) != 2 {
		return "", fmt.Errorf("kbuild not found")
	}

	return fmt.Sprintf("%s%s", baseURL, match[1]), nil
}
