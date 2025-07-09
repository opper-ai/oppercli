# Opper CLI

Command line interface for Opper AI.

## Installation

### Using Homebrew

```shell
brew tap opper-ai/oppercli git@github.com:opper-ai/oppercli
brew install opper
```

### Manual Installation

To install oppercli, download the latest release for your platform on the [Github Releases](https://github.com/opper-ai/oppercli/releases) page. To easily use it, move it to your path.

Example on MacOS:

```shell
sudo curl -o /usr/local/bin/opper https://github.com/opper-ai/oppercli/releases/latest/download/opper-darwin-arm64
sudo chmod 755 /usr/local/bin/opper
```

## Usage

When first starting oppercli, it will prompt you for your API key.

Typing `opper` will show you the following help:

```
Usage:
  opper [command]

Available Commands:
  call        Call a function
  completion  Generate the autocompletion script for the specified shell
  config      Manage API keys and configuration
  functions   Manage functions
  help        Help about any command
  indexes     Manage indexes
  models      Manage models
  traces      Manage traces
  version     Print the version number

Flags:
      --debug        Enable debug output
  -h, --help         help for opper
      --key string   API key to use from config (default "default")

Use "opper [command] --help" for more information about a command.
```

Each command has subcommands that can be viewed using `opper [command] --help`. For example:

```
opper models --help

Manage models

Usage:
  opper models [command]

Examples:
  # List all models
  opper models list
  # Create a new model
  opper models create mymodel litellm-id api-key
  # Test a model
  opper models test mymodel

Available Commands:
  create      Create a new model
  delete      Delete a model
  get         Get model details
  list        List models
  test        Test a model with an interactive prompt

Flags:
  -h, --help   help for models

Global Flags:
      --debug        Enable debug output
      --key string   API key to use from config (default "default")
```

## Command line arguments and stdin

Many commands support receiving input either through command line arguments or standard input.

For example, when using the `call` command:

```shell
# Using command line arguments
opper call myfunction "respond in kind" "what is 2+2?"

# Using a specific model
opper call --model anthropic/claude-3-sonnet myfunction "respond in kind" "what is 2+2?"

# Using stdin
echo "what is 2+2?" | opper call myfunction "respond in kind"
echo '{"name":"Johnny", "age":41}' | opper call myfunction "only print age"
```

When using the `functions chat` command:

```shell
# Using command line arguments
opper functions chat myfunction "Hello there!"

# Using stdin
echo "Hello there!" | opper functions chat myfunction
```

## Adding a custom model

Execution of custom langauge models are done through [LiteLLM](https://docs.litellm.ai/docs/providers). In order for Opper to call your model, you need to provide configuraion appropriate for your model deployment.

Consider the following call:

```shell
opper models create my-model my-id api-key '{"api_base": "https://myoaiservice.azure.com", "api_version": "2024-06-01"}'
```

- `my-model` is the friendly name for this model in Opper, which users in your organization use when calling this model.
- `my-id` is the LiteLLM identifier for this model. Please see the [LiteLLM Providers](https://docs.litellm.ai/docs/providers) documentation for information on this.
- `api-key` is the API key required to connect to this service.
- `json extra` is a JSON object to pass model and deployment specific configuration as required by LiteLLM.

The following are examples for common cloud model deployments:

### Azure

In this example, we are using a GPT4 deployment in Azure. It has the following configuration:

```shell
opper models create example/my-gpt4 azure/my-gpt4-deployment my-api-key-here '{"api_base": "https://my-gpt4-endpoint.openai.azure.com/", "api_version": "2024-06-01"}'
```

- Endpoint: https://my-gpt4-endpoint.openai.azure.com/
- Deployment name: my-gpt4-deployment, which becomes azure/my-gpt4-deployment
- API key: my-api-key-here

## Show usage based on call tag

It is possible to get the usage grouped by a tag which you send as part of the call. In this example, we are passing in `customer_id` in a call:

```python
result, _ = await opper.call(
    name="respond",
    input="What is the capital of Sweden?",
    tags={"customer_id": "my-customer-id"},
)
```

Then we can query for usage per customer_id:

```
opper usage list --from-date=2025-05-15 --to-date=2025-05-16 --fields=total_tokens,cost --group-by=customer_id
Usage Events:

Time Bucket: 2025-05-15T00:00:00Z
Cost: 0.000005
Count: 1
customer_id: another-customer-id
total_tokens: 31

Time Bucket: 2025-05-15T00:00:00Z
Cost: 0.000016
Count: 3
customer_id: my-customer-id
total_tokens: 92

Time Bucket: 2025-05-15T00:00:00Z
Cost: 0.000046
Count: 1
customer_id: <nil>
total_tokens: 51
```

To have a more parsable list, add `--out=csv` to the list command.

### Event Type Filtering

The usage command now defaults to showing **generation events** (AI model calls) for better user experience. You can control this behavior:

```shell
# Shows generation events (default behavior)
opper usage list

# Shows ALL event types (generation, platform, span, embedding, etc.)
opper usage list --event-type=all

# Shows specific event types
opper usage list --event-type=platform
opper usage list --event-type=span
opper usage list --event-type=embedding
```

### Usage Summary

Get a comprehensive breakdown of costs and event counts by type:

```shell
# Summary of generation events (default)
opper usage list --summary

# Summary of all event types
opper usage list --summary --event-type=all

# Summary for specific date range
opper usage list --summary --from-date=2024-01-01 --to-date=2024-01-30
```

### Useful Fields for Generation Events

When analyzing generation events, these fields provide the most valuable insights:

```shell
# Token usage breakdown
opper usage list --fields=prompt_tokens,total_tokens,completion_tokens

# Grouped by model with token details
opper usage list --fields=prompt_tokens,completion_tokens --group-by=model

# Cost analysis over time
opper usage list --event-type=all --from-date=2025-07-05 --to-date=2025-07-10 --graph --graph-type=cost
```

## Building from source

```shell
brew install golang
make install
```

## Shell Completion

The CLI supports shell completion for bash, zsh, fish, and powershell. To enable it:

### Zsh
```shell
# Add this to your ~/.zshrc
source <(opper completion zsh)
```

### Bash
```shell
# Add this to your ~/.bashrc
source <(opper completion bash)
```

### Fish
```shell
opper completion fish | source
# To make it permanent
opper completion fish > ~/.config/fish/completions/opper.fish
```

### PowerShell
```powershell
opper completion powershell | Out-String | Invoke-Expression
# To make it permanent
opper completion powershell > opper.ps1
# Add the generated opper.ps1 file to your PowerShell profile
```

After enabling completion, you can use TAB to autocomplete commands, subcommands, and flags.