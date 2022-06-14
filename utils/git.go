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
	"fmt"
	"os/exec"
	"strings"
)

func InitGit(path string) (string, error) {
	gitCmd := exec.Command("git", "-C", path, "init")
	output, err := gitCmd.Output()

	return strings.Replace(string(output), "\n", "", -1), err
}

func GetLastCommit(path string) (string, error) {
	gitCmd := exec.Command("git", "-C", path, "rev-parse", "head")
	output, err := gitCmd.Output()

	return strings.Replace(string(output), "\n", "", -1), err
}

func GetBeforeLastCommit(path string) (string, error) {
	gitCmd := exec.Command("git", "-C", path, "rev-parse", "head^")
	output, err := gitCmd.Output()

	return strings.Replace(string(output), "\n", "", -1), err
}

func GetShortCommit(path string, hash string) (string, error) {
	gitCmd := exec.Command("git", "-C", path, "rev-parse", "--short", hash)
	output, err := gitCmd.Output()

	return strings.Replace(string(output), "\n", "", -1), err
}

func GetBeforeAndLastDiff(path string, searchFor string) (string, error) {
	gitCmd := exec.Command("git", "-C", path, "diff", "--name-only", "HEAD^", "HEAD", "--", fmt.Sprintf("'*%s*'", searchFor))
	output, err := gitCmd.Output()

	return string(output), err
}

func GetBeforeDiff(path string, searchFor string) (string, error) {
	gitCmd := exec.Command("git", "-C", path, "diff", "--name-only", "HEAD^", "--", fmt.Sprintf("'*%s*'", searchFor))
	output, err := gitCmd.Output()

	return string(output), err
}

func GetCommitsCount(path string) (string, error) {
	gitCmd := exec.Command("git", "-C", path, "rev-list", "--objects", "--all", "--count")
	output, err := gitCmd.Output()

	return strings.Replace(string(output), "\n", "", -1), err
}
