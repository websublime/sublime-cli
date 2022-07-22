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

	CommandRoot       string = "sublime"
	CommandFlagRoot   string = "root"
	CommandFlagConfig string = "config"

	MessageCommandConfigUsage string = "Config file (default is .sublime.json)"
	MessageCommandRootUsage   string = "Project working dir, default to current dir"
	MessageCommandRootShort   string = "CLI tool to manage monorepo packages"

	MessageErrorCommandExecution string = "Unable to execute command."
	MessageErrorCurrentDirectory string = "Unable to get current directory"
	MessageErrorHomeDirectory    string = "Unable to get user home directory"
)
