# Legacy Opper Go CLI — retired

This Go CLI and its Homebrew tap are retired as of October 6, 2026. They no longer receive releases, maintenance, or security fixes. Use the maintained [Opper CLI](https://github.com/opper-ai/cli), published as [`@opperai/cli`](https://www.npmjs.com/package/@opperai/cli).

## Migrate from Homebrew

Both CLIs provide an `opper` command. Remove the legacy installation and tap before installing the replacement:

```shell
brew uninstall opper-ai/oppercli/opper
brew untap opper-ai/oppercli
npm install -g @opperai/cli
opper --version
opper login
```

The maintained CLI requires Node.js 20.12 or later. If you installed the Go CLI manually or with `go install`, use `type -a opper` to locate the old binary and remove it from PATH before installing the replacement.

Keep `~/.oppercli`: the maintained CLI can import its API keys and base URLs on the first non-forced `opper login` when no new configuration exists. It leaves the legacy file in place. Existing `~/.opper/config.json` is not overwritten; use `opper config add` or sign in again if you already configured the new CLI. See the [current CLI documentation](https://docs.opper.ai/developer-tools/cli) for supported commands and flags.

## Retirement status

- The Homebrew formula is disabled and points users to the npm package.
- Automated release publishing and formula updates have been removed.
- Source, tags, and existing release assets are retained for historical reference. Previously installed binaries are not remotely disabled.

The repository is intended to remain archived after these changes are merged. The [historical README](README-legacy.md) documents the old CLI; its installation instructions and command examples are no longer supported.
