package utils

import "os/exec"

func GetLastCommit(path string) (string, error) {
	gitCmd := exec.Command("git", "-C", path, "rev-parse", "head")
	output, err := gitCmd.Output()

	return string(output), err
}

func GetBeforeLastCommit(path string) (string, error) {
	gitCmd := exec.Command("git", "-C", path, "rev-parse", "head^")
	output, err := gitCmd.Output()

	return string(output), err
}

func GetShortCommit(path string, hash string) (string, error) {
	gitCmd := exec.Command("git", "-C", path, "rev-parse", "--short", hash)
	output, err := gitCmd.Output()

	return string(output), err
}

func GetBeforeLastDiff(path string, searchFor string) (string, error) {
	gitCmd := exec.Command("git", "-C", path, "diff", "--name-only", "HEAD^", "HEAD", "--", searchFor)
	output, err := gitCmd.Output()

	return string(output), err
}

func GetCommitsCount(path string) (string, error) {
	gitCmd := exec.Command("git", "-C", path, "rev-list", "--objects", "--all", "--count")
	output, err := gitCmd.Output()

	return string(output), err
}
