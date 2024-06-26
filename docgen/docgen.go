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

package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/diginfra/driverkit/cmd"
	"github.com/spf13/cobra/doc"
)

const outputDir = "docs"
const websiteTemplate = `---
title: %s
weight: %d
---

`

var (
	targetWebsite    bool
	websitePrepender = func(num int) func(filename string) string {
		total := num
		return func(filename string) string {
			num = num - 1
			title := strings.TrimPrefix(strings.TrimSuffix(strings.ReplaceAll(filename, "_", " "), ".md"), fmt.Sprintf("%s/", outputDir))
			return fmt.Sprintf(websiteTemplate, title, total-num)
		}
	}
	websiteLinker = func(filename string) string {
		if filename == "driverkit.md" {
			return "_index.md"
		}
		return filename
	}
)

// docgen
func main() {
	// Get mode
	flag.BoolVar(&targetWebsite, "website", targetWebsite, "")
	flag.Parse()

	// Get root command
	configOpts, err := cmd.NewConfigOptions()
	if err != nil {
		// configOpts will never be nil here
		if configOpts != nil {
			configOpts.Printer.Logger.Fatal("error setting driverkit config options defaults",
				configOpts.Printer.Logger.Args("err", err.Error()))
		} else {
			os.Exit(1)
		}
	}
	rootOpts, err := cmd.NewRootOptions()
	if err != nil {
		configOpts.Printer.Logger.Fatal("error setting driverkit root options defaults",
			configOpts.Printer.Logger.Args("err", err.Error()))
	}
	driverkit := cmd.NewRootCmd(configOpts, rootOpts)
	root := driverkit.Command()
	num := len(root.Commands()) + 1

	// Setup prepender hook
	prepender := func(num int) func(filename string) string {
		return func(filename string) string {
			return ""
		}
	}
	if targetWebsite {
		prepender = websitePrepender
	}

	// Setup links hook
	linker := func(filename string) string {
		return filename
	}
	if targetWebsite {
		linker = websiteLinker
	}

	// Generate markdown docs
	err = doc.GenMarkdownTreeCustom(root, outputDir, prepender(num), linker)
	if err != nil {
		configOpts.Printer.Logger.Fatal("markdown generation", configOpts.Printer.Logger.Args("err", err.Error()))
	}

	if targetWebsite {
		err = os.Rename(path.Join(outputDir, "driverkit.md"), path.Join(outputDir, "_index.md"))
		if err != nil {
			configOpts.Printer.Logger.Fatal("renaming main docs page", configOpts.Printer.Logger.Args("err", err.Error()))
		}
	}

	if err = stripSensitive(); err != nil {
		configOpts.Printer.Logger.Fatal("error replacing sensitive data", configOpts.Printer.Logger.Args("err", err.Error()))
	}
}

func stripSensitive() error {
	f, err := os.Open(outputDir)
	if err != nil {
		return err
	}
	files, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return err
	}

	for _, file := range files {
		filePath := path.Join(outputDir, file.Name())
		file, err := os.ReadFile(filePath)
		if err != nil {
			return err
		}

		envMark := []byte{36} // $
		for _, s := range cmd.Sensitive {
			target := []byte(os.Getenv(s))
			file = bytes.ReplaceAll(file, target, append(envMark, []byte(s)...))
		}
		if err = os.WriteFile(filePath, file, 0666); err != nil {
			return err
		}
	}

	return nil
}
