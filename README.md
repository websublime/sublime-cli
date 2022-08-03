# Sublime CLI

<p align="center">
  <img style="display: inline; margin: 0 6px" alt="GitHub issues" src="https://img.shields.io/github/issues/websublime/sublime-cli?style=flat-square">
  <img style="display: inline; margin: 0 6px" alt="GitHub pull requests" src="https://img.shields.io/github/issues-pr/websublime/sublime-cli?style=flat-square">
  <img style="display: inline; margin: 0 6px" alt="GitHub" src="https://img.shields.io/github/license/websublime/sublime-cli?style=flat-square">
  <img style="display: inline; margin: 0 6px" alt="PRS" src="https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat-square">
  <img style="display: inline; margin: 0 6px" alt="CI" src="https://github.com/websublime/sublime-cli/actions/workflows/release.yml/badge.svg?branch=main">
</p>

<p align="center">
  <img style="display: inline; margin: 0 6px" alt="OSS" src="https://forthebadge.com/images/badges/open-source.svg">
</p>

<p align="center">❄️ SB-CLI (WIP)</p>

Sublime CLI is a tool to create a frontend workspace, libs or packages to distribute as npm or global scripts to use on your page or micro frontend architecture. It is based in vite to build your dists and also it create artifacts to publish on the Websublime cloud. Current we only support github actions. The monorepo is powered by TurboRepo.

<p align="center">
  <img style="display: inline; margin: 0 6px" alt="Sublime Diagram" src="https://user-images.githubusercontent.com/495720/181646023-0828bee5-0ed9-4938-b558-b3b6f723d135.jpeg">
</p>

# Table of contents

- [Installation](#installation)
- [Usage](#usage)

# Usage

[(Back to top)](#table-of-contents)

```bash
CLI tool to manage monorepo packages.

Usage:
  sublime [flags]
  sublime [command]

Available Commands:
  action      Github action command
  completion  Generate the autocompletion script for the specified shell
  create      Create JS/TS packages
  help        Help about any command
  login       Login author on sublime cloud platform.
  register    Register author on sublime cloud platform.
  version     Print the version number of sublime
  workspace   Create a workspace.

Flags:
      --config string   Config file (default is .sublime.json).
  -h, --help            help for sublime
      --root string     Project working dir, default to current dir.

Use "sublime [command] --help" for more information about a command.
```

# Important

After instalation you will need to follow the next steps:

- Run command to register on the platform: ```sublime register```
- After registered please confirm your registration sent by email
- Now go to websublime.dev login and create your organization. Your organization should correspond to the github organization name.
- With your organization created please login thru the cli to create your local identity file: ```sublime login```
- Congrats! You are now able to start creating workspaces on your new organization.

## Create workspace

First let's start to create a workspace monorepo. The creation of the workspace will need some parameters to fullfill package.json needs.

```bash
> sublime workspace --organization websublime
```

The organization should be the github organization name, because artifacts will be release to github and you be able to install it via npm.
The CLI will prompt you with questions to be answer. All are mandatory.

After created your workspace will be ready to create packages inside of it.

## Create package/lib

Creating a library or package. Monorepo has two folders where you can create your packages they are: libs and packages. Packages on libs are designed to be common features to other packages use. You will see that by default one lib is present. This lib is a vite plugin that provide namespace resoltuion automatic between your packages. The CLI will prompt you with questions to be answer. All are mandatory

```bash
> sublime create
```

Packages are created from templates. The templates current supported are: lit, solid, vue and typescript. Maybe in the future will be incremented.

**Default template is: typescript**

Global parameters, can be used with any command before calling the command itself. There are two global parameters:

| Parameter | Description |
|---|---|
| --root | The root folder of your workspace |
| --config | The .sublime.json config file |

If you run the cli from inside your workspace folder this parameters are resolved automatic.

## Github action

Two predefined actions were created when you created an workspace. This actions will trigger based on:
- Branch name as: feat/...
- Tag creation

Also there one release that will put as npm package on github to be consume as that if you need.

The branch will create a snapshot artifact to be use for development needs. On tag will create the artifact for production.
Example of what runs on github action.

```bash
> sublime action --type branch --env "$NODE_ENV"
```

| Parameter | Description |
|---|---|
| --type | Kind is: branch or tag making the diference for prod or dev |
| --env | Environment in which you are right now (dev, prod) |

# Important

You can adjust your workflows if needeed but be aware that changing the predefined action where it runs sublime action command can break your deploysto the websublime cloud.

# Installation

[(Back to top)](#table-of-contents)

**Mandatory dependencies: NodeJS >= 16 and Yarn**

Download the suitable binary for your OS from the list [here](https://github.com/websublime/sublime-cli/releases) and install it. Make sure to make executable with chmod +x.

In your github repo you will need to setup the following secrets:

| Parameter | Description |
|---|---|
| GH_TOKEN | Github token |

This will be used on github actions to create npm deploys, artifacts uploads and releases.

For OSX you will need to allow it to be executed, because it is not signed as a trusted/signed user.

# Contributing

[(Back to top)](#table-of-contents)

Your contributions are always welcome! Please have a look at the [contribution guidelines](CONTRIBUTING.md) first. :tada:

Create branch, work on it and before submit run:
  - git add .
  - git commit -m "feat: title" -m "Description"
  - git add .
  - git commit --amend
  - git push origin feat/... -f

# License

[(Back to top)](#table-of-contents)


The MIT License (MIT) 2022 - [Websublime](https://github.com/websublime/). Please have a look at the [LICENSE.md](LICENSE.md) for more details.
