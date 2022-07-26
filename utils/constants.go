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
	ErrorMissingFile      ErrorType = "EMISSING_FILE"
	ErrorCmdExecution     ErrorType = "ECMD_EXECUTION"
	ErrorMissingDirectory ErrorType = "EMISSING_DIRECTORY"
	ErrorCreateDirectory  ErrorType = "ECREATE_DIRECTORY"
	ErrorPromptInvalid    ErrorType = "EPROMPT_INVALID"
	ErrorInvalidAuthor    ErrorType = "EAUTHOR_INVALID"
	ErrorInvalidTemplate  ErrorType = "ETEMPLATE_INVALID"

	CommandRoot       string = "sublime"
	CommandFlagRoot   string = "root"
	CommandFlagConfig string = "config"

	CommandRegister string = "register"
	CommandLogin    string = "login"

	MessageCommandConfigUsage string = "Config file (default is .sublime.json)."
	MessageCommandRootUsage   string = "Project working dir, default to current dir."
	MessageCommandRootShort   string = "CLI tool to manage monorepo packages."

	MessageErrorAuthorFileMissing string = "Author file not found. Please register first then login to cloud service."
	MessageErrorParseFile         string = "Unable to parse file."
	MessageErrorIndentFile        string = "Unable to indent file."
	MessageErrorWriteFile         string = "Unable to write file"

	// Register command
	MessageCommandRegisterShort string = "Register author on sublime cloud platform."
	MessageCommandRegisterLong  string = `
	As an author you will register in websublime.dev platform to be able
	to create an organization like in github and sync all your packages you create
	in a monorepo style. This packages are intended to be JS UI/Libs whatever and
	they will be available as single scripts or npm packages.
	`
	MessageCommandRegisterNamePrompt     string = "Please provide your name:"
	MessageCommandRegisterUsernamePrompt string = "Please provide your github username:"
	MessageCommandRegisterEmailPrompt    string = "Please provide your email:"
	MessageCommandRegisterPasswordPrompt string = "Please provide a password:"
	MessageCommandRegisterProgressInit   string = "Start registration process."
	MessageCommandRegisterProgressAuthor string = "Author registered. Init local config."
	MessageCommandRegisterLocalAuthor    string = "Author local metadata saved."
	MessageCommandRegisterProgressDone   string = "Almost done!"
	MessageCommandRegisterNextStep       string = "Please check your email and confirm your registration. After confirmed please visit our platform and create your organization to be able to create packages."

	MessageErrorCommandExecution string = "Unable to execute command."
	MessageErrorCurrentDirectory string = "Unable to get current directory."
	MessageErrorHomeDirectory    string = "Unable to get user home directory."

	MessageErrorCommandRegisterNamePrompt     string = "Name provided is not valid."
	MessageErrorCommandRegisterUsernamePrompt string = "Username provided is not valid."
	MessageErrorCommandRegisterEmailPrompt    string = "Email provided is not valid."
	MessageErrorCommandRegisterPasswordPrompt string = "Password provided is not valid."
	MessageErrorCommandRegisterHomeDir        string = "Error creating data author directory."
	MessageErrorCommandRegisterReadTemplate   string = "Unable to read template file."
	MessageErrorCommandRegisterWriteTemplate  string = "Unable to write template file."

	// Login command
	MessageCommandLoginShort string = "Login author on sublime cloud platform."
	MessageCommandLoginLong  string = `
	As an author you will login in websublime.dev platform to be able
	to create an packages and released them to the platform.
	`

	MessageCommandLoginEmailPrompt    string = "Login with your email:"
	MessageCommandLoginPasswordPrompt string = "Login with your password:"
	MessageCommandLoginProgressInit   string = "Attempt to log you in on the cloud platform"
	MessageCommandLoginAuthor         string = "Author loggedin. Init update author data."
	MessageCommandLoginSuccess        string = "Author data update and loggedin."

	MessageErrorCommandLoginEmailPrompt    string = "Email is not valid."
	MessageErrorCommandLoginPasswordPrompt string = "Password is not valid."
)
