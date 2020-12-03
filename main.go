package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/sensu-community/sensu-plugin-sdk/sensu"
	"github.com/sensu-community/sensu-plugin-sdk/templates"
	corev2 "github.com/sensu/sensu-go/api/core/v2"
)

// Config represents the handler plugin config.
type Config struct {
	sensu.PluginConfig
	TopicARN      string
	Message       string
	AssumeRoleARN string
}

var (
	plugin = Config{
		PluginConfig: sensu.PluginConfig{
			Name:     "sensu-aws-sns-handler",
			Short:    "AWS Simple Notification Service Handler",
			Keyspace: "sensu.io/plugins/sensu-aws-sns-handler/config",
		},
	}

	options = []*sensu.PluginConfigOption{
		{
			Path:      "topic-arn",
			Env:       "SNS_TOPIC_ARN",
			Argument:  "topic-arn",
			Shorthand: "t",
			Default:   "",
			Usage:     "The SNS Topic ARN",
			Value:     &plugin.TopicARN,
		},
		{
			Path:      "message-template",
			Env:       "SNS_MESSAGE_TEMPLATE",
			Argument:  "message-template",
			Shorthand: "m",
			Default:   "{{.Entity.Name}}/{{.Check.Name}}: {{.Check.State}}",
			Usage:     "The template for the message sent via SNS",
			Value:     &plugin.Message,
		},
		{
			Path:      "assume-role-arn",
			Env:       "SNS_ASSUME_ROLE_ARN",
			Argument:  "assume-role-arn",
			Shorthand: "a",
			Default:   "",
			Usage:     "The IAM role to assume upon succssful authentication",
			Value:     &plugin.AssumeRoleARN,
		},
	}
)

func main() {
	handler := sensu.NewGoHandler(&plugin.PluginConfig, options, checkArgs, executeHandler)
	handler.Execute()
}

func checkArgs(_ *corev2.Event) error {
	if len(plugin.TopicARN) == 0 {
		return fmt.Errorf("--topic-arn or SNS_TOPIC_ARN environment variable is required")
	}
	if len(plugin.AssumeRoleARN) > 0 {
		if !arn.IsARN(plugin.AssumeRoleARN) {
			return fmt.Errorf("aws-assume-role-arn %s is not a valid ARN", plugin.AssumeRoleARN)
		}
	}
	return nil
}

func executeHandler(event *corev2.Event) error {

	var svc *sns.SNS

	message, err := templates.EvalTemplate("SNSMessage", plugin.Message, event)
	if err != nil {
		return err
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	if arn.IsARN(plugin.AssumeRoleARN) {
		creds := stscreds.NewCredentials(sess, plugin.AssumeRoleARN)
		svc = sns.New(sess, &aws.Config{Credentials: creds})
	} else {
		svc = sns.New(sess)
	}

	// message should be a template with a specific default
	_, err = svc.Publish(&sns.PublishInput{
		Message:  &message,
		TopicArn: &plugin.TopicARN,
	})
	return err
}
