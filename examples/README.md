# Examples

## OSH (Opper Shell)

Get suggestions for shell commands, edit and then run them.

1. Copy `osh` to your `bin` directory
2. Create opper function:
```shell
opper -c shell You are a bash shell assistant. Help the user with creating commands or provide help. Be concise. When responding with a command, just respond with the command. Do not add markdown formatting.
```
3. Run `osh [what you want to do]`

Example: Create a git branch
```shell
$ osh git create and check out branch my-branch

git checkout -b my-branch
```

You can edit the command, and when you press enter, it will be executed.
