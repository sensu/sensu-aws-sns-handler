package main

import (
	"testing"

	corev2 "github.com/sensu/sensu-go/api/core/v2"
	"github.com/stretchr/testify/assert"
)

func TestResolveTemplate(t *testing.T) {
	assert := assert.New(t)
	event := corev2.FixtureEvent("entity1", "check1")
	event.Check.State = "passing"
	plugin.Message = "{{.Entity.Name}}/{{.Check.Name}}: {{.Check.State}}"
	expectedMessage := "entity1/check1: passing"
	msg, _ := resolveTemplate(plugin.Message, event)
	assert.Equal(expectedMessage, msg)
}

