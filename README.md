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

When first starting oppercli, it will prompt you for your API key. You can alternatively provide an envrionment variable with the key instead:

```shell
export OPPER_API_KEY=op-yourkeyhere
```

Typing `opper` will show you the following help:

```
Usage:
  opper <command> <subcommand> [arguments]

Commands:
  functions:
    list [filter]              List functions, optionally filtering by name
    create <name> [instructions] Create a new function
    delete <name>              Delete a function
    get <name>                 Get function details
    chat <name> [message]      Chat with a function
      Input: echo "message" | opper functions chat <name>
             opper functions chat <name> <message...>

  models:
    list [filter]              List custom language models
    create <name> <litellm-id> <key> [extra] Create a new model
    delete <name>              Delete a model
    get <name>                 Get model details

  indexes:
    list [filter]              List indexes
    create <name>              Create a new index
    delete <name>              Delete an index
    get <name>                 Get index details
    query <name> <query>       Query an index
    add <name> <key> <content> Add content to an index
    upload <name> <file>       Upload and index a file

  call <name> <instructions>   Call a function
    Input: echo "input" | opper call <name> <instructions>
           opper call <name> <instructions> <input...>

  help                         Show this help message
  ```

## Command line arguments and stdin

The prompt can be passed on the command line, or as standard input. If you want to pass standard input, and combine it with a prompt, add a `-` on the command line before writing the prompt.

```shell
opper gpt4 tell me a short joke
echo "tell me a short joke" | opper gpt4
echo '{"name":"Johnny", "age":41}' | opper gpt4 - only print age
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
opper models create example/my-gpt4 azure/gpt4-production my-api-key-here '{"api_base": "https://my-gpt4-deployment.openai.azure.com/", "api_version": "2024-06-01"}'
```

- Endpoint: https://my-gpt4-deployment.openai.azure.com/
- Deployment name: gpt4-production, which becomes azure/gpt4-production
- API key: my-api-key-here

## Building from source

```shell
brew install golang
make install
```
