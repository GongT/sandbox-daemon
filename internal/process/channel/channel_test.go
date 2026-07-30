package channel

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChildCloseInput(t *testing.T) {
	cmd := exec.Command("bash", "-c", "exec 0<&-; sleep 1; echo 'hello world'")

	sd := AttachInputSender(cmd)
	defer sd.Destroy()

	w := sd.GetWriter()
	go func() {
		w.WriteLine("test line")
		w.WriteLine("test line")
		w.WriteLine("test line")
		w.WriteLine("test line")
	}()

	err := cmd.Run()
	if err != nil {
		t.Fatal(err)
	}
}

func TestReadOutput(t *testing.T) {
	cmd := exec.Command("bash", "--noprofile", "--norc")

	cmd.Env = []string{
		"PS1=",
	}
	for _, line := range os.Environ() {
		if strings.HasPrefix(line, "PATH=") {
			cmd.Env = append(cmd.Env, line)
		}
	}

	stdin := AttachInputSender(cmd)
	stdout := AttachOutputHandler(cmd)
	stderr := AttachErrorHandler(cmd)
	defer stdin.Destroy()
	defer stdout.Destroy()
	defer stderr.Destroy()

	require.NoError(t, cmd.Start())

	go func() {
		result := cmd.Wait()
		assert.Error(t, result)
	}()
	go func() {
		x := stdout.GetReader("debug")
		y := stderr.GetReader("debug")
		for {
			select {
			case line, ok := <-x.GetChannel():
				if ok {
					fmt.Println("stdout:", line)
				}
			case line, ok := <-y.GetChannel():
				if ok {
					fmt.Println("stderr:", line)
				}
			case <-x.Wait():
				break
			case <-y.Wait():
				break
			}
		}
	}()

	w := stdin.GetWriter()
	r1 := stdout.GetReader("r1")
	r2 := stdout.GetReader("r2")
	e := stderr.GetReader("e")

	w.WriteLine("echo 'hello world'")

	str1, ok1 := r1.Readline()
	assert.Equal(t, "hello world", str1)
	assert.True(t, ok1)

	str2, ok2 := r2.Readline()
	assert.Equal(t, "hello world", str2)
	assert.True(t, ok2)

	w.WriteLine("echo 'error line' >&2")
	str3, ok3 := e.Readline()
	assert.Equal(t, "error line", str3)
	assert.True(t, ok3)

	w.WriteLine("exit 123")
}
