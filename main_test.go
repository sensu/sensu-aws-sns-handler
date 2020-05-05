package main

import (
	"testing"

	corev2 "github.com/sensu/sensu-go/api/core/v2"
	"github.com/stretchr/testify/assert"
)

func TestCheckArgs(t *testing.T) {
	assert := assert.New(t)
	event := corev2.FixtureEvent("entity1", "check1")
	plugin.TopicARN = "arn:aws:sns:us-west-2:012345678901:sensu-alerts"
        assert.NoError(checkArgs(event))
}
