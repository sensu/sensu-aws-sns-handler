[![Sensu Bonsai Asset](https://img.shields.io/badge/Bonsai-Download%20Me-brightgreen.svg?colorB=89C967&logo=sensu)](https://bonsai.sensu.io/assets/nixwiz/sensu-aws-sns-handler)
![Go Test](https://github.com/nixwiz/sensu-aws-sns-handler/workflows/Go%20Test/badge.svg)
![goreleaser](https://github.com/nixwiz/sensu-aws-sns-handler/workflows/goreleaser/badge.svg)

# Sensu AWS SNS Handler

## Table of Contents
- [Overview](#overview)
- [Usage examples](#usage-examples)
  - [Help output](#help-output)
  - [Templates](#templates)
- [Configuration](#configuration)
  - [Asset registration](#asset-registration)
  - [Handler definition](#handler-definition)
  - [Annotations](#annotations)
  - [AWS Credentials](#aws-credentials)
- [Installation from source](#installation-from-source)
- [Additional notes](#additional-notes)
- [Contributing](#contributing)

## Overview

The sensu-aws-sns-handler is a [Sensu Handler][4] that allows you to send alerts via AWS
Simple Notification Service (SNS).

## Usage examples

### Help output

```
AWS Simple Notification Service Handler

Usage:
  sensu-aws-sns-handler [flags]
  sensu-aws-sns-handler [command]

Available Commands:
  help        Help about any command
  version     Print the version number of this plugin

Flags:
  -t, --topic-arn string          The SNS Topic ARN
  -m, --message-template string   The template for the message sent via SNS (default "{{.Entity.Name}}/{{.Check.Name}}: {{.Check.State}}")
  -a, --assume-role-arn string    The IAM role to assume upon succssful authentication
  -h, --help                      help for sensu-aws-sns-handler
```

### Templates

This handler provides options for using templates to populate the values
provided by the event in the message sent via SNS. More information on
template syntax and format can be found in [the documentation][8]


## Configuration

### Asset registration

[Sensu Assets][5] are the best way to make use of this plugin. If you're not using an asset, please
consider doing so! If you're using sensuctl 5.13 with Sensu Backend 5.13 or later, you can use the
following command to add the asset:

```
sensuctl asset add nixwiz/sensu-aws-sns-handler
```

If you're using an earlier version of sensuctl, you can find the asset on the [Bonsai Asset Index][6].

### Handler definition

```yml
---
type: Handler
api_version: core/v2
metadata:
  name: sensu-aws-sns-handler
  namespace: default
spec:
  command: sensu-aws-sns-handler
  type: pipe
  runtime_assets:
  - nixiwz/sensu-aws-sns-handler
```

### Annotations

All arguments for this handler are tunable on a per entity or check basis based on annotations.  The
annotations keyspace for this handler is `sensu.io/plugins/sensu-aws-sns-handler/config`.

**NOTE**: Due to [check token substituion][9], supplying a template value such
as for `message-template` as a check annotation requires that you place the
desired template as a [golang string literal][10] (enlcosed in backticks)
within another template definition.  This does not apply to entity annotations.


#### Examples

To change the SNS topic ARN for a particular entity, in that agent's agent.yml file
add the following:

```yml
annotations:
  sensu.io/plugins/sensu-aws-sns-handler/config/topic-arn: "arn:aws:sns:us-west-2:0123456789012:sensu-alerts"
```

To change the message template for a particular check, for that checks's metadata add the following:

```yml
type: CheckConfig
api_version: core/v2
metadata:
  annotations:
    fatigue_check/occurrences: "3"
    sensu.io/plugins/sensu-aws-sns-handler/config/message-template: "{{`{{.Entity.Name}}/{{.Check.Name}}: {{.Check.State}}, {{.Check.Occurrences}}`}}"
[...]
```

###  AWS Credentials

This plugin makes use of the AWS SDK for Go.  The SDK uses the [default credential provider chain][2]
to find AWS credentials.  The SDK uses the first provider in the chain that returns credentials
without an error. The default provider chain looks for credentials in the following order:

1. Environment variables (AWS_SECRET_ACCESS_KEY, AWS_ACCESS_KEY_ID, and AWS_REGION).

2. Shared credentials file (typically ~/.aws/credentials).

3. If your application is running on an Amazon EC2 instance, IAM role for Amazon EC2.

4. If your application uses an ECS task definition or RunTask API operation, IAM role for tasks.

The SDK detects and uses the built-in providers automatically, without requiring manual configurations.
For example, if you use IAM roles for Amazon EC2 instances, your applications automatically use the
instance’s credentials. You don’t need to manually configure credentials in your application.

Source: [Configuring the AWS SDK for Go][3]

This plugin also supports assuming a new role upon authentication using the `--assume-role-arn`
option.

If you go the route of using environment variables, it is highly suggested you use them via the
[Env secrets provider][7].

## Installation from source

The preferred way of installing and deploying this plugin is to use it as an Asset. If you would
like to compile and install the plugin from source or contribute to it, download the latest version
or create an executable from this source.

From the local path of the sensu-aws-sns-handler repository:

```
go build
```

## Contributing

For more information about contributing to this plugin, see [Contributing][1].

[1]: https://github.com/sensu/sensu-go/blob/master/CONTRIBUTING.md
[2]: https://docs.aws.amazon.com/sdk-for-go/api/aws/defaults/#CredChain
[3]: https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html
[4]: https://docs.sensu.io/sensu-go/latest/reference/handlers/
[5]: https://docs.sensu.io/sensu-go/latest/reference/assets/
[6]: https://bonsai.sensu.io/assets/nixwiz/sensu-aws-sns-handler
[7]: https://docs.sensu.io/sensu-go/latest/guides/secrets-management/#use-env-for-secrets-management
[8]: https://docs.sensu.io/sensu-go/latest/observability-pipeline/observe-process/handler-templates/
[9]: https://docs.sensu.io/sensu-go/latest/observability-pipeline/observe-schedule/checks/#check-token-substitution
[10]: https://golang.org/ref/spec#String_literals
