/*
 * Copyright 2020 Infolabs Inc & Associates
 *
 *    Licensed under the Apache License, Version 2.0 (the "License");
 *    you may not use this file except in compliance with the License.
 *    You may obtain a copy of the License at
 *
 *        http://www.apache.org/licenses/LICENSE-2.0
 *
 *    Unless required by applicable law or agreed to in writing, software
 *    distributed under the License is distributed on an "AS IS" BASIS,
 *    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *    See the License for the specific language governing permissions and
 *    limitations under the License.
 */

package mpesa

import (
	"fmt"
	"os"
	"strconv"
)

const (
	defMajorVer = 2
	defMinorVer = 5
	defPatchVer = 8
	defDevStage = 3
	defDevStageIter = 2

	envMajorVer = "MAJOR_VERSION"
	envMinorVer = "MINOR_VERSION"
	envPatchVer = "PATCH_VERSION"
	envDevStage = "DEV_STAGE"
	envDevStageIter = "STAGE_ITERATION"
)

//type devStage int

const (
	Alpha = iota +1
	Beta
	ReleaseCandidate
	Release
	PostReleaseFix
)

var devStageName = map[int]string{
	Alpha: "Alpha",
	Beta: "Beta",
	ReleaseCandidate:"Release Candidate",
	Release: "Release",
	PostReleaseFix: "Post Release Fix",

}

func stage(stage int)string  {
	return devStageName[stage]
}



type version struct {
	Major int
	Minor int
	Patch int
	DevStage int
	DevStageIter int
}


//SemVer returns the semantic version of the software
//Semantic version is universal and widely used way of software versioning.
//It contains 3 components, in the format of MAJOR.MINOR.PATCH where:
//1. MAJOR: When you do incompatible changes, means it may contains changes which
//are may not be compatible with older version, or major feature or functional changes to it.
//2. MINOR: When you do compatible changes, means it add features backwards-compatible manner,
//3. PATCH: when you make backwards-compatible bug fixes.
//4. QUALIFIER: If you’re doing release very fluently and you need one for pre-release or
//additional qualifier you can use qualifier at last.
func (v version) semVer() string {

	var qualifier string
	switch v.DevStage {
	case Alpha:
		qualifier = "a"

	case Beta:
		qualifier = "b"

	case ReleaseCandidate:
		qualifier = "rc"

	default:
		qualifier = ""
	}


	//e.g 1.2.0 and 1.2.5
	if qualifier == ""{

		if v.DevStage==Release{
			return fmt.Sprintf("%d.%d.0",v.Major, v.Minor)
		}

		return fmt.Sprintf("%d.%d.%d",v.Major, v.Minor,v.Patch)
	}


	//e.g 1.4.8-a.2
	return fmt.Sprintf("%d.%d.%d-%s.%d",v.Major, v.Minor,v.Patch,qualifier,v.DevStageIter)

}

func loadVersionEnv(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		env, _ := strconv.Atoi(v)

		return env

	}

	return fallback
}

func loadConfig() version {

	return version{
		Major:        loadVersionEnv(envMajorVer,defMajorVer),
		Minor:        loadVersionEnv(envMinorVer,defMinorVer),
		Patch:        loadVersionEnv(envPatchVer,defPatchVer),
		DevStage:     loadVersionEnv(envDevStage,defDevStage),
		DevStageIter: loadVersionEnv(envDevStageIter,defDevStageIter),
	}

}

func Version() string {
	v := loadConfig()
	return v.semVer()
}
