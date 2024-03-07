# Opper CLI

## Installation

```shell
brew install golang
make install
export OPPER_API_KEY=op-yourkeyhere
```

## Functions

### Create

```shell
opper -c joch/joker Respond to input with Linux or unix related jokes
```

### Delete

```shell
opper -d joch/joker
```

### List

List all functions:

```shell
opper -l
```

Filter the list by a search term:

```shell
opper -l joch
```

## Examples

```shell
opper -c joch/gpt4

opper joch/gpt4 tell me a short joke
```

```shell
opper -c joch/diff You are provided with a diff. Generate a summary of the changes in bullet form.

git diff | opper joch/diff
```

```shell
opper -c joch/shell You are a bash shell assistant. Help the user with creating commands or provide help. Be concise. When responding with a command, just respond with the command. Do not add markdown formatting.

opper joch/shell git revert to last commit
```
