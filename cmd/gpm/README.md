# Good Package Manager (GPM)

The package manager Gotham deserves, but not the one it needs right now.

A super minimal git-based package manager that has no opinions about anything really.

Built because I wanted something simple to use for languages that don't have inbuilt package managers (or have shitty ones) that was slightly smoother than gitmodules, and thought it was interesting.

## Status

TBH I wouldn't use it if you weren't me, but, I'm not your parent.

- [X] Project Init
- [X] Module Add
- [X] Project Sync (and re-sync)
- [X] Semver module matching
- [X] Git hash based lock files
- [ ] Graceful handling of dirty modules
- [ ] Dry run mode w/ human and machine readable diffs
- [ ] Auto .gitignore management for modules
- [ ] PGP signature pinning and checking for git tags / commits

## Usage

Use `gpm --help` to list available commands, and `gpm OPTION --help` to list arguments for a given command.

1. Init a project with `gpm init` to create a `.gpm.yml` package file
2. [Re]sync a project (and modules) with `gpm sync`
3. Add dependencies with `gpm add`
4. Update dependencies (within semver ranges) with `gpm update`

