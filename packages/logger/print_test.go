package logger

import (
	"log"
	"strings"
	"testing"

	"github.com/gongt/sandbox-daemon/packages/myenv"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func refresh() *strings.Builder {
	var output strings.Builder
	log.SetOutput(&output)
	return &output
}

func random() string {
	uuid, _ := uuid.NewV7()
	return uuid.String()
}

func TestEnableDisable(t *testing.T) {
	myenv.RedirectDebugTesting(t)

	require.True(t, Disable("*"))

	var output *strings.Builder
	var message string
	output = refresh()

	message = random()
	REFLECT.Enable()
	DLog(string(REFLECT), "%s", message)
	assert.Contains(t, output.String(), message)

	message = random()
	output = refresh()
	REFLECT.Disable()
	DLog(string(REFLECT), "%s", message)
	assert.NotContains(t, output.String(), message)

	message = random()
	output = refresh()
	Enable("test_tag*")
	DLog("test_tag", "%s", message)
	Disable("test_tag*")
	assert.Contains(t, output.String(), message)

	message = random()
	output = refresh()
	Enable("*test_tag")
	DLog("test_tag", "%s", message)
	Disable("*test_tag")
	assert.Contains(t, output.String(), message)

	message = random()
	output = refresh()
	Enable("t*g")
	DLog("test_tag", "%s", message)
	Disable("t*g")
	assert.Contains(t, output.String(), message)
}
