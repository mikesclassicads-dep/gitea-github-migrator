# gitea-github-migrator
[![Build Status](https://drone.jonasfranz.software/api/badges/JonasFranzDEV/gitea-github-migrator/status.svg)](https://drone.jonasfranz.software/JonasFranzDEV/gitea-github-migrator)

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
- [ ] Pull Requests
- [ ] Statuses

## Installation

```bash
go get git.jonasfranz.software/JonasFranzDEV/gitea-github-migrator
cd $GOPATH/src/git.jonasfranz.software/JonasFranzDEV/gitea-github-migrator
dep ensure
go install
```

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
