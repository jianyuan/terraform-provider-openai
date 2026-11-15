# Terraform Provider OpenAI

Unofficial Terraform provider for OpenAI.

> [!WARNING]
> **This provider is deprecated.** OpenAI now publishes an official Terraform provider. Please use it instead:
>
> - Terraform Registry: [openai/openai](https://registry.terraform.io/providers/openai/openai/latest/docs)
> - Source code: [openai/terraform-provider-openai](https://github.com/openai/terraform-provider-openai)
> - Guide: [Terraform guide on the OpenAI developer docs](https://developers.openai.com/api/docs/guides/terraform)

## Thank you

Thank you to everyone who has used this provider since 2024.

I built this provider and maintain the [Sentry provider](https://github.com/jianyuan/terraform-provider-sentry) (29.5M+ downloads). If you need a Terraform provider for your own API, or help building AI agents and LLM features, I take on consulting work through [jian.fyi](https://jian.fyi/?utm_source=github&utm_medium=referral&utm_campaign=terraform-provider-openai).

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.21

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Using the provider

Please refer to the [provider documentation on the Terraform registry](https://registry.terraform.io/providers/jianyuan/openai/latest/docs).

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

_Note:_ Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```
