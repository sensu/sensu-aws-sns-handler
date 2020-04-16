package main

import (
	"bytes"
	"fmt"
	"io"
	"text/template"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/sensu-community/sensu-plugin-sdk/sensu"
	"github.com/sensu/sensu-go/types"
)

// Config represents the handler plugin config.
type Config struct {
	sensu.PluginConfig
	TopicARN string
	Message  string
}

// used to handle getting text/template or html/template
type templater interface {
        Execute(wr io.Writer, data interface{}) error
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
		&sensu.PluginConfigOption{
			Path:      "topic-arn",
			Env:       "SNS_TOPIC_ARN",
			Argument:  "topic-arn",
			Shorthand: "t",
			Default:   "",
			Usage:     "The SNS Topic ARN",
			Value:     &plugin.TopicARN,
		},
		&sensu.PluginConfigOption{
			Path:      "message-template",
			Env:       "SNS_MESSAGE_TEMPLATE",
			Argument:  "message-template",
			Shorthand: "m",
			Default:   "{{.Entity.Name}}/{{.Check.Name}}: {{.Check.State}}",
			Usage:     "The template for the message sent via SNS",
			Value:     &plugin.Message,
		},
	}
)

func main() {
	handler := sensu.NewGoHandler(&plugin.PluginConfig, options, checkArgs, executeHandler)
	handler.Execute()
}

func checkArgs(_ *types.Event) error {
	if len(plugin.TopicARN) == 0 {
		return fmt.Errorf("--topic-arn or SNS_TOPIC_ARN environment variable is required")
	}
	return nil
}

func executeHandler(event *types.Event) error {

	message, err := resolveTemplate(plugin.Message, event)
	if err != nil {
		return err
	}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := sns.New(sess)

	// message should be a template with a specific default
	_, err = svc.Publish(&sns.PublishInput{
		Message:  &message,
		TopicArn: &plugin.TopicARN,
	})
	return err
}

func resolveTemplate(templateValue string, event *types.Event) (string, error) {
        var (
                resolved bytes.Buffer
                tmpl     templater
                err      error
        )

        tmpl, err = template.New("test").Parse(templateValue)

        if err != nil {
                return "", err
        }

        err = tmpl.Execute(&resolved, *event)
        if err != nil {
                return "", err
        }

        return resolved.String(), nil
}

