package chorus

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func (c *Config) Build(targets ...string) error {
	if len(targets) == 0 {
		targets = []string{"all"}
	}

	for _, target := range targets {
		if err := c.processTarget(target); err != nil {
			return err
		}
	}
	return nil
}

func (c *Config) processTarget(target string) error {
	t, exists := c.Targets[target]
	if !exists && target != "all" {
		return fmt.Errorf("target '%s' undefined", target)
	}

	if t.executed {
		return nil
	}
	t.executed = true

	for _, dep := range t.Deps {
		if err := c.processTarget(dep); err != nil {
			return err
		}
	}

	if t.Phony || needsRebuild(target, t.Deps) {
		return c.executeCommands(t.Cmds, target)
	}

	fmt.Printf("%s %s\n", yellow("→"), target)
	return nil
}

func needsRebuild(target string, deps []string) bool {
	if target == "all" || strings.HasPrefix(target, "_") {
		return true
	}

	info, err := os.Stat(target)
	if os.IsNotExist(err) {
		return true
	}

	targetTime := info.ModTime()
	for _, dep := range deps {
		depInfo, err := os.Stat(dep)
		if err != nil {
			return true
		}
		if depInfo.ModTime().After(targetTime) {
			return true
		}
	}
	return false
}

func (c *Config) executeCommands(cmds []string, target string) error {
	fmt.Printf("%s %s\n", cyan("→"), target)

	for _, cmdStr := range cmds {
		cmdStr = c.expandVariables(cmdStr, target)

		command := exec.Command("sh", "-c", cmdStr)
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr

		if err := command.Run(); err != nil {
			fmt.Printf("%s %s\n", red("✗"), cmdStr)
			return fmt.Errorf("command failed: %w", err)
		}
	}
	return nil
}

func (c *Config) expandVariables(cmd string, target string) string {
	replacements := map[string]string{
		"@":  target,
		"@F": filepath.Base(target),
		"@D": filepath.Dir(target),
		"<":  c.firstDependency(target),
		"^":  strings.Join(c.Targets[target].Deps, " "),
	}

	for k, v := range c.Variables {
		cmd = strings.ReplaceAll(cmd, "${"+k+"}", v)
	}

	for k, v := range replacements {
		cmd = strings.ReplaceAll(cmd, "${"+k+"}", v)
	}

	return cmd
}

func (c *Config) firstDependency(target string) string {
	if deps := c.Targets[target].Deps; len(deps) > 0 {
		return deps[0]
	}
	return ""
}
