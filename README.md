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