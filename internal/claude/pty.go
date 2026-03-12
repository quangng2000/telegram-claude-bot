package claude

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/creack/pty"
)

func execWithPTY(args []string, workDir string, timeout time.Duration) (string, error) {
	cmd := exec.Command("claude", args...)
	cmd.Dir = workDir
	cmd.Env = cleanEnv()

	ptmx, err := pty.Start(cmd)
	if err != nil {
		return "", fmt.Errorf("failed to start claude: %w", err)
	}
	defer ptmx.Close()

	pty.Setsize(ptmx, &pty.Winsize{Rows: 50, Cols: 200})
	return readWithTimeout(ptmx, cmd, timeout)
}

func readWithTimeout(ptmx *os.File, cmd *exec.Cmd, timeout time.Duration) (string, error) {
	done := make(chan struct{})
	var output strings.Builder

	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := ptmx.Read(buf)
			if n > 0 {
				output.Write(buf[:n])
			}
			if err != nil {
				break
			}
		}
		close(done)
	}()

	select {
	case <-done:
		cmd.Wait()
		return output.String(), nil
	case <-time.After(timeout):
		cmd.Process.Kill()
		return "", fmt.Errorf("timed out after %v", timeout)
	}
}

func cleanEnv() []string {
	env := os.Environ()
	clean := make([]string, 0, len(env))
	for _, e := range env {
		if !strings.HasPrefix(e, "CLAUDECODE=") {
			clean = append(clean, e)
		}
	}
	return clean
}
