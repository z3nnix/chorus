package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"regexp"
	"strings"
	"syscall"
	"time"

	"github.com/fatih/color"
	"gopkg.in/yaml.v3"
)

// ANSI стили
var (
	cyan    = color.New(color.FgCyan).SprintFunc()
	green   = color.New(color.FgGreen).SprintFunc()
	red     = color.New(color.FgRed).SprintFunc()
	yellow  = color.New(color.FgYellow).SprintFunc()
	magenta = color.New(color.FgMagenta).SprintFunc()
	blue    = color.New(color.FgBlue).SprintFunc()
	bold    = color.New(color.Bold).SprintFunc()
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
		fmt.Printf("\n\n%s %s\n", red("✗"), bold("Build cancelled!"))
		os.Exit(1)
	}()

	printHeader()
	
	if len(os.Args) < 2 {
		processTarget(loadConfig(), "all")
	} else {
		config := loadConfig()
		for _, arg := range os.Args[1:] {
			processTarget(config, arg)
		}
	}

	fmt.Printf("\n%s %s %s\n\n",
		green("✓"),
		bold("Build completed successfully!"),
		blue("("+time.Since(startTime).Round(time.Millisecond).String()+")"),
	)
}

func printHeader() {
	header := `
    _______ ______  ___  __  ______
    / ___/ // / __ \/ _ \/ / / / __/
   / /__/ _  / /_/ / , _/ /_/ /\ \  
   \___/_//_/\____/_/|_|\____/___/  
                                    
`
	fmt.Println(magenta(header))
}

func loadConfig() *BuildConfig {
	data, err := os.ReadFile("chorus.build")
	if err != nil {
		exitWithError("Error reading build file:", err)
	}

	config := &BuildConfig{}
	if err := yaml.Unmarshal(data, config); err != nil {
		exitWithError("Error parsing build file:", err)
	}

	config.Variables["DATE"] = time.Now().Format("2006-01-02")
	return config
}

func processTarget(config *BuildConfig, target string) {
	t, exists := config.Targets[target]
	if !exists && target != "all" {
		exitWithError(fmt.Sprintf("Target '%s' not defined", target))
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
		printSuccess(target)
	} else {
		printSkipped(target)
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
	fmt.Printf("%s %s\n",
		cyan("●"),
		bold("Processing target:")+" "+magenta(target),
	)
	
	for _, cmd := range cmds {
		if strings.HasPrefix(cmd, "internal:") {
			handleInternalCommand(cmd)
			continue
		}

		cmd = expandVariables(config, cmd, target)
		start := time.Now()
		
		fmt.Printf("  %s %s\n",
			blue("⌛"),
			cyan(strings.TrimSpace(cmd)),
		)

		command := exec.Command("sh", "-c", cmd)
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		
		if err := command.Run(); err != nil {
			fmt.Printf("\r  %s %s %s\n",
				red("✗"),
				cyan(strings.TrimSpace(cmd)),
				red(fmt.Sprintf("[FAIL] (%s)", time.Since(start).Round(time.Millisecond))),
			)
			exitWithError("Command failed:", err)
		}
		
		fmt.Printf("\r  %s %s %s\n",
			green("✓"),
			cyan(strings.TrimSpace(cmd)),
			green(fmt.Sprintf("[OK] (%s)", time.Since(start).Round(time.Millisecond))),
		)
	}
}

func handleInternalCommand(cmd string) {
	start := time.Now()
	cmdParts := strings.SplitN(cmd, ":", 2)
	command := strings.TrimSpace(cmdParts[1])

	fmt.Printf("  %s %s\n",
		blue("⌛"),
		cyan(command),
	)

	var err error
	switch {
	case strings.HasPrefix(command, "load_nvm"):
		args := strings.Fields(command)
		if len(args) < 2 {
			err = fmt.Errorf("app file path required")
		} else {
			err = loadNVMHeader(args[1])
		}
	case strings.HasPrefix(command, "restore_nvm"):
		err = restoreNVMHeader()
	default:
		err = fmt.Errorf("unknown internal command")
	}

	if err != nil {
		fmt.Printf("\r  %s %s %s\n",
			red("✗"),
			cyan(command),
			red(fmt.Sprintf("[FAIL] (%s)", time.Since(start).Round(time.Millisecond))),
		)
		exitWithError("Internal command failed:", err)
	}

	fmt.Printf("\r  %s %s %s\n",
		green("✓"),
		cyan(command),
		green(fmt.Sprintf("[OK] (%s)", time.Since(start).Round(time.Millisecond))),
	)
}

func loadNVMHeader(appFile string) error {
	headerFile := "core/kernel/nvm/nvm.h"

	appContent, err := os.ReadFile(appFile)
	if err != nil {
		return fmt.Errorf("failed to read app file: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(appContent)), "\n")
	var processedLines []string
    
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
        
		escaped := strings.ReplaceAll(trimmed, `"`, `\"`)
		processedLines = append(processedLines, `"`+escaped+` \n"`)
	}

	result := strings.Join(processedLines, " ")
	if len(result) > 4 {
		result = result[:len(result)-4] + `"`
	}

	if err := copyFile(headerFile, headerFile+".bak"); err != nil {
		return fmt.Errorf("backup failed: %v", err)
	}

	headerContent, err := os.ReadFile(headerFile)
	if err != nil {
		return fmt.Errorf("failed to read header file: %v", err)
	}

	pattern := regexp.MustCompile(`static const char apps\[\] = ".*";`)
	newContent := pattern.ReplaceAllString(
		string(headerContent),
		fmt.Sprintf(`static const char apps[] = %s;`, result),
	)

	if err := os.WriteFile(headerFile, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("failed to write header file: %v", err)
	}

	return nil
}

func restoreNVMHeader() error {
	headerFile := "core/kernel/nvm/nvm.h"
	backupFile := headerFile + ".bak"

	if _, err := os.Stat(backupFile); os.IsNotExist(err) {
		return fmt.Errorf("backup file not found")
	}

	if err := copyFile(backupFile, headerFile); err != nil {
		return fmt.Errorf("restore failed: %v", err)
	}

	if err := os.Remove(backupFile); err != nil {
		return fmt.Errorf("failed to remove backup: %v", err)
	}

	return nil
}

func copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	return os.WriteFile(dst, input, 0644)
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

func printSuccess(target string) {
	fmt.Printf("%s %s %s\n\n",
		green("✔"),
		bold("Target"),
		magenta(target)+green(" completed successfully!"),
	)
}

func printSkipped(target string) {
	fmt.Printf("%s %s %s\n\n",
		yellow("ⓘ"),
		bold("Skipping"),
		cyan(target)+yellow(" (already up-to-date)"),
	)
}

func exitWithError(args ...interface{}) {
	fmt.Printf("\n%s %s\n\n",
		red("✗"),
		bold(red("BUILD FAILED!")),
	)
	log.Fatal(red(args...))
}
