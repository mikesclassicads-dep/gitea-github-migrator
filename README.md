# gitea-github-migrator
[![Build Status](https://drone.jonasfranz.software/api/badges/JonasFranzDEV/gitea-github-migrator/status.svg)](https://drone.jonasfranz.software/JonasFranzDEV/gitea-github-migrator)
[![Latest Release](https://img.shields.io/badge/dynamic/json.svg?label=release&url=https%3A%2F%2Fgit.jonasfranz.software%2Fapi%2Fv1%2Frepos%2FJonasFranzDEV%2Fgitea-github-migrator%2Freleases&query=%24%5B0%5D.tag_name)](https://git.jonasfranz.software/JonasFranzDEV/gitea-github-migrator/releases)
[![Docker Pulls](https://img.shields.io/docker/pulls/jonasfranz/gitea-github-migrator.svg)](https://hub.docker.com/r/jonasfranz/gitea-github-migrator/)

A tool to migrate [GitHub](https://github.com) Repositories to [Gitea](https://gitea.io) including all issues, labels, milestones
and comments.

## Features
Migrates:

- [x] Repositories
- [x] Issues
- [x] Labels
- [x] Milestones
- [x] Comments
- [ ] Users
- [x] Pull Requests (as issue)
- [ ] Statuses

## Installation

### From source

```bash
go get git.jonasfranz.software/JonasFranzDEV/gitea-github-migrator
cd $GOPATH/src/git.jonasfranz.software/JonasFranzDEV/gitea-github-migrator
dep ensure
go install
```

### From Binary
We provide binaries of master builds and all releases at our [minio storage server](https://storage.h.jonasfranz.software/minio/gitea-github-migrator/dist/).

Additionally we provide them for every release as release attachment under [releases](https://git.jonasfranz.software/JonasFranzDEV/gitea-github-migrator/releases).

You don't need any dependencies except the binary to run the migrator.
## Usage

Migrate one repository:
```bash
gitea-github-migrator migrate \
    --gh-repo owner/reponame \
    --gh-token GITHUB_TOKEN \
    --url http://gitea-url.tdl \
    --token GITEA_TOKEN \
    --owner 1
```
`gh-token` is only required if you have more than 50 issues / repositories.

Migrate all repositories:
```bash
gitea-github-migrator migrate-all \
    --gh-user username \
    --gh-token GITHUB_TOKEN \
    --url http://gitea-url.tdl \
    --token GITEA_TOKEN \
    --owner 1
```

Migrate all repositories without issues etc. (classic):
```bash
gitea-github-migrator migrate-all \
    --gh-user username \
    --gh-token GITHUB_TOKEN \
    --url http://gitea-url.tdl \
    --token GITEA_TOKEN \
    --owner 1
    --only-repo
```

## Problems
* This migration tool does not work with Gitea instances using a SQLite database.
* Comments / Issues will be added in the name of the user to whom belongs the token (information about the original date and author will be added)
* The current date is used for creation date (information about the actual date is added in a comment)
