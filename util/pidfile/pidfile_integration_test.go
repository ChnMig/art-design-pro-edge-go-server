package pidfile

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"
)

func TestPidFileLifecycle_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("跳过集成测试（-short）")
	}

	tmpDir := t.TempDir()
	pidPath := filepath.Join(tmpDir, "api-server.pid")

	exe, err := os.Executable()
	if err != nil {
		t.Fatalf("获取测试二进制路径失败: %v", err)
	}

	cmd := exec.Command(exe, "-test.run", "^TestPidFileHelperProcess$", "-test.v")
	cmd.Env = append(os.Environ(),
		"PIDFILE_HELPER_PROCESS=1",
		"PIDFILE_PATH="+pidPath,
	)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Start(); err != nil {
		t.Fatalf("启动 helper 进程失败: %v", err)
	}
	t.Cleanup(func() {
		if cmd.Process != nil && cmd.ProcessState == nil {
			_ = cmd.Process.Kill()
			_ = cmd.Wait()
		}
	})

	wantPID := strconv.Itoa(cmd.Process.Pid)
	deadline := time.Now().Add(2 * time.Second)
	for {
		data, err := os.ReadFile(pidPath)
		if err == nil {
			gotPID := strings.TrimSpace(string(data))
			if gotPID == wantPID {
				break
			}
		}

		if time.Now().After(deadline) {
			t.Fatalf("pid 文件未在预期时间内写入或内容不正确，want=%s，当前输出：\n%s", wantPID, out.String())
		}
		time.Sleep(10 * time.Millisecond)
	}

	if err := cmd.Process.Signal(os.Interrupt); err != nil {
		t.Fatalf("发送中断信号失败: %v", err)
	}

	waitCh := make(chan error, 1)
	go func() { waitCh <- cmd.Wait() }()

	select {
	case err := <-waitCh:
		if err != nil {
			t.Fatalf("helper 进程退出失败: %v\n输出：\n%s", err, out.String())
		}
	case <-time.After(5 * time.Second):
		_ = cmd.Process.Kill()
		_ = <-waitCh
		t.Fatalf("helper 进程未在预期时间内退出\n输出：\n%s", out.String())
	}

	if _, err := os.Stat(pidPath); !os.IsNotExist(err) {
		t.Fatalf("进程退出后 pid 文件应被删除，Stat() err=%v", err)
	}
}

func TestPidFileHelperProcess(t *testing.T) {
	if os.Getenv("PIDFILE_HELPER_PROCESS") != "1" {
		return
	}

	path := os.Getenv("PIDFILE_PATH")
	if err := Write(path, os.Getpid()); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "写入 pid 文件失败: %v\n", err)
		os.Exit(2)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	if err := Remove(path); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "删除 pid 文件失败: %v\n", err)
		os.Exit(3)
	}
	os.Exit(0)
}
