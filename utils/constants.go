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

type ErrorType string

type PackageType string

const (
	Library PackageType = "lib"
	Package PackageType = "pkg"
)

const (
	ErrorUnknown          ErrorType = "E_UNKNOWN"
	ErrorOpenFile         ErrorType = "EOPEN_FILE"
	ErrorReadFile         ErrorType = "EREAD_FILE"
	ErrorCmdExecution     ErrorType = "ECMD_EXECUTION"
	ErrorMissingDirectory ErrorType = "EMISSING_DIRECTORY"
	ErrorPromptInvalid    ErrorType = "EPROMPT_INVALID"

	CommandRoot       string = "sublime"
	CommandFlagRoot   string = "root"
	CommandFlagConfig string = "config"

	CommandRegister string = "register"

	MessageCommandConfigUsage   string = "Config file (default is .sublime.json)."
	MessageCommandRootUsage     string = "Project working dir, default to current dir."
	MessageCommandRootShort     string = "CLI tool to manage monorepo packages."
	MessageCommandRegisterShort string = "Register author on sublime cloud platform."
	MessageCommandRegisterLong  string = `
	As an author you will register in websublime.dev platform to be able
	to create an organization like in github and sync all your packages you create
	in a monorepo style. This packages are intended to be JS UI/Libs whatever and
	they will be available as single scripts or npm packages.
	`
	MessageCommandRegisterNamePrompt          string = "Please provide your name:"
	MessageErrorCommandRegisterNamePrompt     string = "Name provided is not valid."
	MessageCommandRegisterUsernamePrompt      string = "Please provide your github username:"
	MessageErrorCommandRegisterUsernamePrompt string = "Username provided is not valid."
	MessageCommandRegisterEmailPrompt         string = "Please provide your email:"
	MessageErrorCommandRegisterEmailPrompt    string = "Email provided is not valid."
	MessageCommandRegisterPasswordPrompt      string = "Please provide a password:"
	MessageErrorCommandRegisterPasswordPrompt string = "Password provided is not valid."

	MessageErrorCommandExecution string = "Unable to execute command."
	MessageErrorCurrentDirectory string = "Unable to get current directory."
	MessageErrorHomeDirectory    string = "Unable to get user home directory."
)
