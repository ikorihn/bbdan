# bbdan

`bbdan` is an unofficial command line tool wrapping Bitbucket Cloud REST API 2.0.

## Motivation

Bitbucket Cloud REST API 2.0 does not provide the way to manage permission of a repository.
(see <https://community.atlassian.com/t5/Bitbucket-questions/Setting-repo-permissions-using-BB-Cloud-API-2-0/qaq-p/1038792>)

With so many repositories, I would be happy if I could at least delete or update permissions through command line, so I developed this.

## Features

- List, delete permissions for a repository.
- Copy permissions for a repository to another repository.

## Install

### Build from source

```shell
$ go install
```

## Configure

Create `config.toml` at your `$XDG_CONFIG_HOME/bbdan` (`~/.config/bbdan` if unset) directory and configure like below.

`$XDG_CONFIG_HOME/bbdan/config.toml`

```toml
username = "your-bitbucket-user-name"
password = "your-app-passwords"
```

### App password

You should generate [app password](https://support.atlassian.com/bitbucket-cloud/docs/app-passwords/).
Make sure it has permissions `repository:admin` .

## Usage

To see all available commands, use `bbdan -h` .

### `permission list`

List permissions for a repository.

```shell
$ bbdan permission list workspace repository
List permissions for workspace/repository
==== RESULT ====
type, id, name, permission
group, developer, developer, read
user, {aaaaaaaa-8888-1111-abcd-12345abc}, user-1, write
```

### `permission copy`

Copy permissions of a repository to another repository.

```shell
$ bbdan permission copy workspace my-repository other-repository
Copy permissions from workspace/my-repository to workspace/other-repository
? Choose operations:  [Use arrows to move, space to select, <right> to all, <left> to none, type to filter]
> [ ]  Add: user user-1 (WRITE)
  [ ]  Add: user user-2 (ADMIN)
  [ ]  Update: group developer WRITE => ADMIN
```

With `--batch` (`-b`), force overwrite permissions of target repository

```shell
$ bbdan permission copy -b workspace my-repository other-repository
Copy permissions from workspace/my-repository to workspace/other-repository
```

### `permission remove`

Select and remove permission of a repository.

```shell
$ bbdan permission remove workspace repository
Remove selected permissions from workspace/repository
? Choose operations:  [Use arrows to move, space to select, <right> to all, <left> to none, type to filter]
> [ ]  Remove: group developer (WRITE)
  [ ]  Remove: user user-1 (WRITE)
  [ ]  Remove: user user-2 (ADMIN)
```
