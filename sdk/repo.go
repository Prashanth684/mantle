// Copyright 2015 CoreOS, Inc.
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

package sdk

import (
	"log"
	"os"
	"path/filepath"
)

const (
	// In the SDK chroot the repo is always at this location
	chrootRepoRoot = "/mnt/host/source"

	// Assorted paths under the repo root
	defaultCacheDir = ".cache"
	defaultBuildDir = "src/build"
)

func isDir(dir string) bool {
	stat, err := os.Stat(dir)
	return err == nil && stat.IsDir()
}

func envDir(env string) string {
	dir := os.Getenv(env)
	if dir == "" {
		return ""
	}
	if !filepath.IsAbs(dir) {
		log.Fatalf("%s is not an absolute path: %q", env, dir)
	}
	return dir
}

func RepoRoot() string {
	if dir := envDir("REPO_ROOT"); dir != "" {
		return dir
	}

	if isDir(chrootRepoRoot) {
		return chrootRepoRoot
	}

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Invalid working directory: %v", err)
	}

	for dir := wd; ; dir = filepath.Dir(dir) {
		if isDir(filepath.Join(dir, ".repo")) {
			return dir
		} else if filepath.IsAbs(dir) {
			break
		}
	}

	return wd
}

func RepoCache() string {
	return filepath.Join(RepoRoot(), defaultCacheDir)
}

func BuildRoot() string {
	if dir := envDir("BUILD_ROOT"); dir != "" {
		return dir
	}
	return filepath.Join(RepoRoot(), defaultBuildDir)
}
