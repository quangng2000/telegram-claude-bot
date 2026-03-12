package claude

import (
	"time"
)

type Runner struct {
	Model   string
	Timeout time.Duration
}

func NewRunner(model string, timeoutMin int) *Runner {
	return &Runner{
		Model:   model,
		Timeout: time.Duration(timeoutMin) * time.Minute,
	}
}

func (r *Runner) Run(workDir, prompt string) (string, error) {
	args := []string{
		"-p",
		"--model", r.Model,
		"--permission-mode", "bypassPermissions",
		"--continue",
		prompt,
	}
	return r.exec(args, workDir, r.Timeout)
}

func (r *Runner) RunOneShot(workDir, prompt string) (string, error) {
	args := []string{
		"-p",
		"--model", r.Model,
		"--no-session-persistence",
		prompt,
	}
	return r.exec(args, workDir, 60*time.Second)
}

func (r *Runner) exec(args []string, workDir string, timeout time.Duration) (string, error) {
	raw, err := execWithPTY(args, workDir, timeout)
	if err != nil {
		return "", err
	}
	clean := StripAnsi(raw)
	if clean == "" {
		return "No response from Claude.", nil
	}
	return clean, nil
}
