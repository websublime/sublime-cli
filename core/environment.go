package core

import "os"

type Environment struct {
	Workspace      string
	WorkspaceRoot  string
	WorkspaceScope string
}

func NewEnvironment() *Environment {
	dir, _ := os.Getwd()

	return &Environment{
		WorkspaceRoot: dir,
	}
}

var environment = NewEnvironment()

func GetEnvironment() *Environment {
	return environment
}
