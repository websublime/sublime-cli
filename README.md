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

<p align="center">❄️ SB-CLI</p>

Sublime CLI is a tool to create a frontend workspace, libs or packages to distribute as npm or global scripts to use on your page or micro frontend architecture. It is based in vite to build your dists and also it create artifacts to publish on Supabase storage. Current we only support lit and github actions.

# Table of contents

- [Usage](#usage)
- [Installation](#installation)

# Usage

[(Back to top)](#table-of-contents)

```bash
CLI tool to manage projects

Usage:
  sublime [flags]
  sublime [command]

Available Commands:
  action      Creates artifacts on github actions
  create      Create libs or packages
  help        Help about any command
  version     Print the version number of sublime
  workspace   Create a workspace project

Flags:
      --config string   config file (default is .sublime.json)
  -h, --help            help for sublime-cli
      --root string     Project working dir, default to current dir

Use "sublime [command] --help" for more information about a command.
```

## Create workspace

First let's start to create a workspace monorepo. The creation of the workspace will need some parameters to fullfill package.json needs.

```bash
> sublime workspace --name ws-libs-ui --scope @ws --repo sublime/ws-libs-ui --username miguelramos --email miguel@websublime.dev
```

| Parameter | Description |
|---|---|
| --name | This will be the folder name for creating workspace |
| --scope | This is the scope(organization) prefix namespace |
| --repo | This is mandatory short github repo name |
| --username | This will be used on package.json definitions |
| --email | Same porpose as username |

After the creation of your workspace just get into the workspace folder and run yarn.

Now inside your workspace let's create a library or a package.

- Library, could be something that you want to share thru your packages
- Packages, independents lit components for your micro frontends

## Create package/lib

Creating a library or package is on the same command, only parameter changes

```bash
> sublime create --name utils --type lib --template lit
```

| Parameter | Description |
|---|---|
| --name | This will be the folder name for creating lib or pkg, also will be prepend with scope for packing name |
| --type | Defines the type, supported are: lib or pkg |
| --template | Template to use for your lib/pkg (current only lit but will have solid, vue, react, typescript only) |

Global parameters, can be used with any command before calling the command itself. There are two global parameters:

| Parameter | Description |
|---|---|
| --root | The root folder of your workspace |
| --config | The .sublime.json config file |

If you run the cli from inside your workspace folder this para meters are resolved.

## Perform github action

Two predefined actions were created when you created an workspace. This actions will trigger based on:
- Branch name as: feat/...
- Tag creation

Also there one release that will put as npm package on github to be consume as that if you need.

The branch will create a snapshot artifact to be use for development needs. On tag will create the artifact for production.
Example of what runs on github action.

```bash
> sublime action --kind branch --bucket "$BUCKET" --url "$STORAGE_URL" --key "$STORAGE_KEY" --env "$NODE_ENV"
```

| Parameter | Description |
|---|---|
| --kind | Kind is: branch or tag making the diference for prod or dev |
| --bucket | Storage bucket name to put assets from dist |
| --url | Storage base url |
| --key | Storage api key |
| --env | Environment in which you are right now (dev, prod) |

For now only Supabase storage is supported.

# Installation

[(Back to top)](#table-of-contents)

**Mandatory dependencies: NodeJS >= 16 and Yarn**

Download the suitable binary for your OS from the list [here](https://github.com/websublime/sublime-cli/releases) and install it. Make sure to make executable with chmod +x.

In your github repo you will need to setup the following secrets:

| Parameter | Description |
|---|---|
| GH_TOKEN | Github token |
| BUCKET | Storage bucket name |
| STORAGE_URL | Base Storage url |
| STORAGE_KEY | Storage api secret key |

This will be used on github actions to create npm deploys, artifacts uploads and releases.

For OSX you will need to allow it to be executed, because it is not signed as a trusted/signed user.

Also current all artifacts are upload to supabase storage, so you need to create two buckets there:

- One will be used to deploy all dist assets on it
- Second create a bucket with the name "manifests", where this will have info about wich package or lib you release from github actions. This manifest file can be used to serve the main scripts for your system.

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
