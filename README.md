# gitea-github-migrator

[![Build Status](https://drone.jonasfranz.software/api/badges/JonasFranzDEV/gitea-github-migrator/status.svg)](https://drone.jonasfranz.software/JonasFranzDEV/gitea-github-migrator)
[![Latest Release](https://img.shields.io/badge/dynamic/json.svg?label=release&url=https%3A%2F%2Fgit.jonasfranz.software%2Fapi%2Fv1%2Frepos%2FJonasFranzDEV%2Fgitea-github-migrator%2Freleases&query=%24%5B0%5D.tag_name)](https://git.jonasfranz.software/JonasFranzDEV/gitea-github-migrator/releases)
[![Docker Pulls](https://img.shields.io/docker/pulls/ggmigrator/cli.svg)](https://hub.docker.com/r/ggmigrator/cli/)
[![Docker Pulls](https://img.shields.io/docker/pulls/ggmigrator/web.svg)](https://hub.docker.com/r/ggmigrator/web/)
[![Go Report Card](https://goreportcard.com/badge/git.jonasfranz.software/JonasFranzDEV/gitea-github-migrator)](https://goreportcard.com/report/git.jonasfranz.software/JonasFranzDEV/gitea-github-migrator)
[![GoDoc](https://godoc.org/git.jonasfranz.software/JonasFranzDEV/gitea-github-migrator?status.svg)](https://godoc.org/git.jonasfranz.software/JonasFranzDEV/gitea-github-migrator)
[![Coverage Status](https://coverage.jonasfranz.software/projects/1/badge.svg)](https://coverage.jonasfranz.software/projects/1)

A tool to migrate [GitHub](https://github.com) Repositories to [Gitea](https://gitea.io) including all issues, labels, milestones
and comments.


# IMPORTANT

**This repository got moved to the offical Gitea server and will be developed there. https://gitea.com/gitea/migrator/**

There will be some big changes on how the migrator works in order to improve the performance and the result of a migration.

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
make build
```
#### Web Support

Run `make web-build` instead of `make build` to include web support.

### From Binary
We provide binaries of master builds and all releases at our [minio storage server](https://storage.h.jonasfranz.software/minio/gitea-github-migrator/dist/).

Additionally we provide them for every release as release attachment under [releases](https://git.jonasfranz.software/JonasFranzDEV/gitea-github-migrator/releases).

You don't need any dependencies except the binary to run the migrator.

These binaries include web support by default.

### From Docker image

We provide a [cli docker image](https://hub.docker.com/r/ggmigrator/cli/):

For master builds:
```docker
docker run ggmigrator/cli:latest
```

For release builds:
```docker
docker run ggmigrator/cli:0.0.10
```


## Usage

### Command line

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

### Web interface

Since 0.1.0 gitea-github-migrator comes with an integrated web interface.
Follow these steps to get the web interface running:

1. Download or build a web-capable binary of the gitea-github-migrator. The builds on our storage server are build with web support included.
If you build from source, please follow [web support](#web-support).
2. Create `config.yml` file and change the properties according to your wishes. Please keep in mind that
you have to create a GitHub OAuth application to make the web interface work.
3. Run `./gitea-github-migrator web`
4. Visit `http://localhost:4000`

#### Docker

We're providing a docker image with web support. To start the web service please run:
```docker
docker run ggmigrator/web -p 4000:4000 -v data/:/data
```
Place your `config.yml` into `data/config.yml`.

#### Config
Example:
```yaml
# GitHub contains the OAuth2 application data obtainable from GitHub
GitHub:
  client_id: GITHUB_OAUTH_CLIENT_ID
  client_secret: GITHUB_OAUTH_SECRET
# Web contains the configuration for the integrated web server
Web:
  port: 4000
  host: 0.0.0.0
```

## Problems

- This migration tool does not work with Gitea instances using a SQLite database.
- Comments / Issues will be added in the name of the user to whom belongs the token (information about the original date and author will be added)
- The current date is used for creation date (information about the actual date is added in a comment)
