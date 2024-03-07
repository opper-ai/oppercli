# Opper CLI

## Prerequisites

```shell
brew install golang
```

## Installation

```shell
make install
```

## Configuration

```shell
export OPPER_API_KEY=op-yourkeyhere
```

## Examples

```shell
opper joch/gpt3 tell me a short joke
```

```shell
diff -u test-*.py | opper joch/diff
```

```shell
opper joch/shell git revert to last commit
```

## Create function

```shell
opper -c joch/joker Respond to input with Linux or unix related jokes
```

## Delete function

```shell
opper -d joch/joker
```
