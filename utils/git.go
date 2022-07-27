/*
Copyright © 2022 Websublime.dev organization@websublime.dev

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package utils

import (
	"os/exec"
	"strings"
)

func InitGit(path string) (string, error) {
	gitCmd := exec.Command("git", "init")
	gitCmd.Dir = path
	output, err := gitCmd.Output()

	return strings.Replace(string(output), "\n", "", -1), err
}

func GetLastCommit(path string) (string, error) {
	gitCmd := exec.Command("git", "rev-parse", "head")
	gitCmd.Dir = path
	output, err := gitCmd.Output()

	return strings.Replace(string(output), "\n", "", -1), err
}

func GetBeforeLastCommit(path string) (string, error) {
	gitCmd := exec.Command("git", "rev-parse", "head^")
	gitCmd.Dir = path
	output, err := gitCmd.Output()

	return strings.Replace(string(output), "\n", "", -1), err
}

func GetShortCommit(path string, hash string) (string, error) {
	gitCmd := exec.Command("git", "rev-parse", "--short", hash)
	gitCmd.Dir = path
	output, err := gitCmd.Output()

	return strings.Replace(string(output), "\n", "", -1), err
}

func GetBeforeAndLastDiff(path string) (string, error) {
	gitCmd := exec.Command("git", "--no-pager", "diff", "--name-only", "HEAD^", "HEAD")
	gitCmd.Dir = path
	output, err := gitCmd.CombinedOutput()

	return string(output), err
}

func GetBeforeDiff(path string) (string, error) {
	gitCmd := exec.Command("git", "--no-pager", "diff", "--name-only", "HEAD^")
	gitCmd.Dir = path
	output, err := gitCmd.Output()

	return string(output), err
}

func GetCommitsCount(path string) (string, error) {
	gitCmd := exec.Command("git", "rev-list", "--objects", "--all", "--count")
	gitCmd.Dir = path
	output, err := gitCmd.Output()

	return strings.Replace(string(output), "\n", "", -1), err
}

func GetDiffBetweenTags(path string) (string, error) {
	gitCmd := exec.Command("git", "rev-list", "--tags", "--max-count=1")
	gitCmd.Dir = path
	lastTagByte, er := gitCmd.Output()

	if er != nil {
		return "", er
	}

	gitCmd = exec.Command("git", "rev-list", "--tags", "--skip=1", "--max-count=1")
	gitCmd.Dir = path
	tagBeforeByte, erro := gitCmd.Output()

	if erro != nil {
		return "", er
	}

	lastTag := strings.Replace(string(lastTagByte), "\n", "", -1)
	tagBefore := strings.Replace(string(tagBeforeByte), "\n", "", -1)

	if tagBefore == "" {
		tagBefore = "main"
	}

	if lastTag == "" {
		lastTag = "HEAD"
	}

	gitCmd = exec.Command("git", "--no-pager", "diff", "--name-only", tagBefore, lastTag)
	gitCmd.Dir = path
	output, err := gitCmd.Output()

	return string(output), err
}
