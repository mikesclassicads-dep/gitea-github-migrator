# gitea-github-migrator

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
`gh-token` is only required if you have more then 50 issues / repositories.

Migrate all repositories:
```bash
gitea-github-migrator migrate-all \
    --gh-user username \
    --gh-token GITHUB_TOKEN \
    --url http://gitea-url.tdl \
    --token GITEA_TOKEN \
    --owner 1
```

Migrate all repository without issues etc. (classic):
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
* This migrator does not work with Gitea instances utilizing a SQLite database.
* Comments / Issues will be added by token's user (information about date and author will be added)
* Current Date is utilized for creation date (information about date is added in a comment)