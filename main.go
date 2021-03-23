package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/arn"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/ec2metadata"
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
	Subject       string
	AssumeRoleARN string
	UseEC2Region  bool
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
			Path:      "subject-template",
			Env:       "SNS_SUBJECT_TEMPLATE",
			Argument:  "subject-template",
			Shorthand: "s",
			Default:   "{{.Check.State}} - {{.Entity.Name}}/{{.Check.Name}}",
			Usage:     "The template for the subject sent via SNS",
			Value:     &plugin.Subject,
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
		{
			Path:      "use-ec2-region",
			Env:       "",
			Argument:  "use-ec2-region",
			Shorthand: "u",
			Default:   false,
			Usage:     "Query the EC2 metadata for the region to use for SNS",
			Value:     &plugin.UseEC2Region,
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

	subject, err := templates.EvalTemplate("SNSSubject", plugin.Subject, event)
	if err != nil {
		return err
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	awsConfig := &aws.Config{}

	if plugin.UseEC2Region {
		ec2md := ec2metadata.New(sess)
		region, err := ec2md.Region()
		if err != nil {
			return fmt.Errorf("Cannot determine region from EC2 metadata: %v", err)
		}
		awsConfig.Region = aws.String(region)
	}

	if arn.IsARN(plugin.AssumeRoleARN) {
		awsConfig.Credentials = stscreds.NewCredentials(sess, plugin.AssumeRoleARN)
	}
	svc = sns.New(sess, awsConfig)

	// message should be a template with a specific default
	publishOut, err := svc.Publish(&sns.PublishInput{
		Subject:  &subject,
		Message:  &message,
		TopicArn: &plugin.TopicARN,
	})
	if err != nil {
		return fmt.Errorf("Failed to publish message to SNS: %v", err)
	}

	// FUTURE: send to activities hub
	fmt.Printf("SNS message published, message id %s\n", *publishOut.MessageId)
	return nil

}
