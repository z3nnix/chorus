# CLI Reference

This document describes the command-line interface for Chorus.

## Synopsis

```bash
chorus [target1] [target2] ... [targetN]
```

## Description

Chorus is a build automation tool that reads build instructions from a `chorus.build` file in the current directory and executes the specified targets.

## Usage

### Basic Usage

Run the default target (`all`):
```bash
chorus
```

Run specific targets:
```bash
chorus target1
```

Run multiple targets in sequence:
```bash
chorus target1 target2 target3
```

## Arguments

### Targets

Targets are the names defined in the `targets` section of your `chorus.build` file.

**Examples:**
```bash
# Build the 'all' target
chorus

# Build just the 'program' target
chorus program

# Clean and rebuild
chorus clean all

# Run tests
chorus test

# Install the built program
chorus install
```

If no target is specified, Chorus runs the `all` target by default.

### Multiple Targets

You can specify multiple targets in a single command:

```bash
chorus clean build test
```

Targets are executed in the order specified, left to right.

## Current Working Directory

Chorus must be run from a directory containing a `chorus.build` file. It does not search parent directories.

**Example directory structure:**
```
myproject/
├── chorus.build    ← Must run chorus from here
├── src/
│   ├── main.c
│   └── utils.c
└── include/
    └── utils.h
```

## Exit Codes

Chorus uses the following exit codes:

| Code | Meaning |
|------|---------|
| `0` | Success - all targets built successfully |
| `1` | Error - build failed, configuration error, or interrupt |

**Examples:**

```bash
chorus build
echo $?  # Prints 0 if successful

chorus invalid_target
echo $?  # Prints 1 (target not found)
```

## Output

Chorus provides colored output to indicate build progress:

### Output Indicators

| Indicator | Color | Meaning |
|-----------|-------|---------|
| `→ target` | Cyan | Building target (executing commands) |
| `→ target` | Yellow | Target up-to-date (skipping) |
| `✗ command` | Red | Command failed |
| `done 123ms` | Green | Build completed successfully |
| `interrupt` | Red | Build interrupted by user (Ctrl+C) |

### Example Output

```bash
$ chorus
→ main.o
→ utils.o
→ program
done 234ms
```

With up-to-date targets:
```bash
$ chorus
→ main.o
→ utils.o
→ program

done 12ms
```

After making changes:
```bash
$ touch src/main.c
$ chorus
→ main.o
→ utils.o
→ program

done 145ms
```

### Verbose Command Output

When Chorus executes commands, stdout and stderr from those commands are displayed directly:

```bash
$ chorus
→ main.o
src/main.c: In function 'main':
src/main.c:10:5: warning: unused variable 'x' [-Wunused-variable]
→ program

done 156ms
```

## Signal Handling

### Keyboard Interrupt (Ctrl+C)

Pressing Ctrl+C during a build will:
1. Stop the current command
2. Display `interrupt` message
3. Exit with code 1

**Example:**
```bash
$ chorus long_build
→ compiling...
^C
interrupt
```

The build is stopped immediately, and any partial build artifacts may remain.

### SIGTERM

Chorus also responds to SIGTERM signals, allowing graceful shutdown in automated environments.

## Environment

### Working Directory

Chorus executes all commands in the directory where `chorus.build` is located.

### Shell

All commands are executed using `sh -c`, which means:
- Standard shell features are available (pipes, redirects, etc.)
- Commands have access to environment variables
- Shell-specific features depend on your system's `/bin/sh`

**Example with shell features:**
```yaml
targets:
  archive:
    cmds:
      - "tar czf backup.tar.gz *.o && echo 'Archived'"
      - "ls -lh backup.tar.gz | awk '{print $5}'"
```

## Examples

### Common Build Workflows

**Clean build:**
```bash
chorus clean all
```

**Quick rebuild (just the main target):**
```bash
chorus program
```

**Build and test:**
```bash
chorus all test
```

**Build and install:**
```bash
chorus all install
```

**Force complete rebuild:**
```bash
rm -f *.o program
chorus
```

### Integration with Scripts

**Build script:**
```bash
#!/bin/bash
set -e

echo "Starting build..."
chorus clean all

echo "Running tests..."
chorus test

echo "Build successful!"
```

**Continuous Integration:**
```bash
#!/bin/bash
# CI build script

# Exit on any error
set -e

# Build project
chorus all || {
    echo "Build failed"
    exit 1
}

# Run tests
chorus test || {
    echo "Tests failed"
    exit 1
}

echo "CI build successful"
```

### Using with Make-style Patterns

Chorus can replace make in many scenarios:

**Before (Makefile):**
```make
all: program

program: main.o utils.o
	gcc -o $@ $^

main.o: main.c
	gcc -c main.c -o $@

clean:
	rm -f *.o program
```

**After (chorus.build):**
```yaml
targets:
  all:
    deps: [program]

  program:
    deps: [main.o, utils.o]
    cmds:
      - "gcc -o ${@} ${^}"

  main.o:
    deps: [main.c]
    cmds:
      - "gcc -c main.c -o ${@}"

  clean:
    phony: true
    cmds:
      - "rm -f *.o program"
```

**Usage is similar:**
```bash
# Instead of: make
chorus

# Instead of: make clean
chorus clean

# Instead of: make program
chorus program
```

## Limitations

### No Options or Flags

Currently, Chorus does not accept command-line options or flags. There are no equivalents to:
- `make -j4` (parallel builds)
- `make -n` (dry run)
- `make -B` (force rebuild)
- `make --version`

### No Help Command

There is no built-in `--help` or `-h` option. Refer to this documentation for usage information.

### No Parallel Execution

Chorus builds targets sequentially. Parallel builds are not supported.

## Troubleshooting

### "read config: open chorus.build: no such file or directory"

**Cause:** No `chorus.build` file in current directory

**Solution:** Ensure you're in the correct directory or create a `chorus.build` file

```bash
ls -la chorus.build  # Check if file exists
pwd                  # Verify current directory
```

### "target 'X' undefined"

**Cause:** Target does not exist in `chorus.build`

**Solution:** Check target name spelling and ensure it's defined

```bash
chorus all  # Try the default target instead
```

### Commands not executing as expected

**Cause:** Shell-specific syntax not available in `/bin/sh`

**Solution:** Use POSIX-compatible shell syntax, or invoke bash explicitly:

```yaml
cmds:
  - "bash -c 'your bash-specific command here'"
```

### Build appears to hang

**Cause:** A command is waiting for input or running indefinitely

**Solution:** Press Ctrl+C to interrupt and investigate the problematic command

## Tips and Tricks

### Timing Builds

The `done` message includes build time:
```bash
$ time chorus
→ program

done 145ms

real    0m0.152s
user    0m0.098s
sys     0m0.041s
```

### Conditional Builds

Use shell conditionals in commands:
```yaml
targets:
  conditional:
    cmds:
      - "test -f config.h || cp config.h.default config.h"
```

### Logging Build Output

Redirect output to a log file:
```bash
chorus all 2>&1 | tee build.log
```

### Checking if a Target Exists

Since Chorus doesn't have a dry-run mode, check your `chorus.build`:
```bash
grep "target_name:" chorus.build
```

Or use a YAML parser:
```bash
yq '.targets | keys' chorus.build
```
