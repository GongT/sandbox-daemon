package instance

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestProcessGroup(t *testing.T) {
	g := NewProcessGroup()

	tmpDir := filepath.Join(os.TempDir(), "test")
	t.Logf("创建临时目录: %s", tmpDir)
	err := os.MkdirAll(tmpDir, 0755)
	require.NoError(t, err)

	g.SetDir("/")
	g.SetOverlayRoot(tmpDir)

	err = os.WriteFile(filepath.Join(tmpDir, "hello.txt"), []byte("Hello, World!"), 0644)
	if err != nil {
		panic(err)
	}

	init := g.CreateProcess([]string{"sleep", "infinity"})
	init_buff := test_output(init)

	init.Start()
	defer init.Join()
	defer init.Kill(os.Kill)

	go func() {
		stat, err := init.Join()
		t.Logf("init输出 %s", init_buff.String())
		require.Error(t, err, "asdasd")
		require.Equal(t, -1, stat.ExitCode())
	}()

	time.Sleep(1 * time.Second) // 等待init启动完成

	// test_fast_execute(t, g, "pwd")
	// test_fast_execute(t, g, "ls", "-lhA", ".")
	// test_fast_execute(t, g, "mount")
	test_fast_execute(t, g, "bash", "-c", "echo wow > /xxx.txt")

	// test
	test := g.CreateProcess([]string{"cat", "/hello.txt"})
	buff := test_output(test)

	test.Start()
	test.Join()
	require.Equal(t, "Hello, World!", buff.String())
	require.FileExists(t, filepath.Join(tmpDir, "xxx.txt"))

	require.Equal(t, 1, len(g.child_processes.instances))

	os.RemoveAll(tmpDir)
}

func test_output(instance *ProcessInstance) *bytes.Buffer {
	var outBuffer bytes.Buffer
	instance.SetBeforeStartHook(func(c *exec.Cmd) {
		c.Stdout = &outBuffer
		c.Stderr = &outBuffer
	})

	return &outBuffer
}

func test_fast_execute(t *testing.T, g *ProcessGroup, cmdline ...string) {
	// test
	test := g.CreateProcess(cmdline)
	buff := test_output(test)

	test.Start()
	test.Join()

	t.Logf("------------- %v -------------\n%s--------------------------------", cmdline, buff.String())
}
