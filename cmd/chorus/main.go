package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/fatih/color"
	"gopkg.in/yaml.v3"
)

var (
	cyan   = color.New(color.FgCyan).SprintFunc()
	green  = color.New(color.FgGreen).SprintFunc()
	red    = color.New(color.FgRed).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
)

type BuildConfig struct {
	Variables map[string]string `yaml:"variables"`
	Targets   map[string]struct {
		Deps     []string `yaml:"deps"`
		Cmds     []string `yaml:"cmds"`
		Phony    bool     `yaml:"phony"`
		Executed bool
	} `yaml:"targets"`
}

func main() {
	color.NoColor = false
	startTime := time.Now()
	
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Printf("\n%s\n", red("interrupt"))
		os.Exit(1)
	}()
	
	config := loadConfig()
	
	if len(os.Args) < 2 {
		processTarget(config, "all")
	} else {
		for _, arg := range os.Args[1:] {
			processTarget(config, arg)
		}
	}

	fmt.Printf("\n%s %s\n", green("done"), time.Since(startTime).Round(time.Millisecond))
}

func loadConfig() *BuildConfig {
	data, err := os.ReadFile("chorus.build")
	if err != nil {
		exitWithError("read config:", err)
	}

	config := &BuildConfig{}
	if err := yaml.Unmarshal(data, config); err != nil {
		exitWithError("parse config:", err)
	}

	config.Variables["DATE"] = time.Now().Format("2006-01-02")
	return config
}

func processTarget(config *BuildConfig, target string) {
	t, exists := config.Targets[target]
	if !exists && target != "all" {
		exitWithError(fmt.Sprintf("target '%s' undefined", target))
	}

	if t.Executed {
		return
	}
	t.Executed = true

	for _, dep := range t.Deps {
		processTarget(config, dep)
	}

	if t.Phony || needsRebuild(target, t.Deps) {
		executeCommands(config, t.Cmds, target)
	} else {
		fmt.Printf("%s %s\n", yellow("→"), target)
	}
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

func executeCommands(config *BuildConfig, cmds []string, target string) {
	fmt.Printf("%s %s\n", cyan("→"), target)
	
	for _, cmd := range cmds {
		cmd = expandVariables(config, cmd, target)
		
		command := exec.Command("sh", "-c", cmd)
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		
		if err := command.Run(); err != nil {
			fmt.Printf("%s %s\n", red("✗"), cmd)
			exitWithError("command failed:", err)
		}
	}
}

func expandVariables(config *BuildConfig, cmd string, target string) string {
	replacements := map[string]string{
		"@": target,
		"<": firstDependency(config, target),
		"^": strings.Join(config.Targets[target].Deps, " "),
	}

	for k, v := range config.Variables {
		cmd = strings.ReplaceAll(cmd, "${"+k+"}", v)
	}

	for k, v := range replacements {
		cmd = strings.ReplaceAll(cmd, "${"+k+"}", v)
	}

	return cmd
}

func firstDependency(config *BuildConfig, target string) string {
	if deps := config.Targets[target].Deps; len(deps) > 0 {
		return deps[0]
	}
	return ""
}

func exitWithError(args ...interface{}) {
	log.Fatal(red(args...))
}