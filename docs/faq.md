# Frequently Asked Questions (FAQ)

## General Questions

### What is Chorus?

Chorus is a minimalistic build automation tool written in Go. It uses YAML-based configuration files to define build targets and their dependencies, providing a simpler alternative to traditional build systems like Make.

### Why use Chorus instead of Make?

Chorus offers several advantages:
- **Simpler syntax**: YAML is more readable than Makefile syntax
- **No tab issues**: YAML uses spaces, avoiding Make's tab/space confusion
- **Modern design**: Built with modern development practices in mind
- **Clear error messages**: Better error reporting than Make
- **Easy to learn**: Minimal concepts to understand

### Is Chorus production-ready?

Chorus is actively used in projects like NovariaOS and Perano lang. However, it's a minimalistic tool with a focused feature set. Evaluate whether its features meet your project's needs.

### What platforms does Chorus support?

Chorus is designed for Unix-like systems:
- Linux (all distributions)
- BSD variants (FreeBSD, OpenBSD, NetBSD)
- macOS
- Other POSIX-compliant systems

Windows is not officially supported.

## Installation and Setup

### How do I install Chorus?

Download and install the binary:
```bash
wget https://github.com/z3nnix/chorus/releases/download/1.0.2/chorus && \
sudo mv chorus /usr/bin/chorus && \
sudo chmod +x /usr/bin/chorus
```

See [Getting Started](getting-started.md) for details.

### Can I install Chorus without sudo?

Yes, install to a user directory:
```bash
wget https://github.com/z3nnix/chorus/releases/download/1.0.2/chorus
chmod +x chorus
mkdir -p ~/.local/bin
mv chorus ~/.local/bin/
```

Then ensure `~/.local/bin` is in your PATH:
```bash
export PATH="$HOME/.local/bin:$PATH"
```

### How do I build Chorus from source?

```bash
git clone https://github.com/z3nnix/chorus.git
cd chorus
go build -o chorus ./cmd/chorus
```

Requirements: Go 1.16 or later

### How do I update Chorus?

Download the latest release and replace the binary:
```bash
wget https://github.com/z3nnix/chorus/releases/latest/download/chorus -O /tmp/chorus
sudo mv /tmp/chorus /usr/bin/chorus
sudo chmod +x /usr/bin/chorus
```

## Configuration

### Where should I put my chorus.build file?

In the root directory of your project, alongside your source files:
```
myproject/
├── chorus.build    ← Here
├── src/
├── include/
└── README.md
```

### Can I use a different filename?

No, Chorus only looks for `chorus.build` in the current directory.

### Can Chorus search parent directories for chorus.build?

No, you must run Chorus from the directory containing `chorus.build`.

### How do I comment my chorus.build file?

Use `#` for comments:
```yaml
# This is a comment
variables:
  CC: "gcc"  # Inline comment
```

### Can I include other build files?

No, Chorus doesn't support including or importing other files. All configuration must be in a single `chorus.build` file.

### How do I handle different build configurations?

Create separate targets for each configuration:
```yaml
targets:
  debug:
    deps: [main_debug.o]
    cmds:
      - "gcc -g -o program_debug main_debug.o"

  release:
    deps: [main_release.o]
    cmds:
      - "gcc -O3 -o program main_release.o"
```

## Variables

### How do I define a variable?

In the `variables` section:
```yaml
variables:
  CC: "gcc"
  CFLAGS: "-Wall -O2"
```

### Can variables reference other variables?

Yes:
```yaml
variables:
  BASE_FLAGS: "-Wall -Werror"
  CFLAGS: "${BASE_FLAGS} -O2"
  DEBUG_FLAGS: "${BASE_FLAGS} -g"
```

### Can I use environment variables?

Yes, using shell command substitution:
```yaml
variables:
  HOME_DIR: "$(echo $HOME)"
  USER: "$(whoami)"
```

### What automatic variables are available?

| Variable | Meaning |
|----------|---------|
| `${@}` | Current target name |
| `${<}` | First dependency |
| `${^}` | All dependencies |
| `${DATE}` | Current date (YYYY-MM-DD) |

### Can I pass variables from the command line?

No, Chorus doesn't support command-line variable overrides. Use environment variables with command substitution instead:
```bash
CC=clang chorus
```

Then in your chorus.build:
```yaml
variables:
  CC: "$(echo ${CC:-gcc})"  # Use CC env var, default to gcc
```

## Targets

### What is a phony target?

A phony target doesn't represent a real file. It's used for tasks like cleaning, testing, or installing:
```yaml
targets:
  clean:
    phony: true
    cmds:
      - "rm -f *.o program"
```

### When should I mark a target as phony?

Mark targets as phony when they:
- Don't create a file (clean, test, install, run)
- Should always execute regardless of timestamps
- Represent actions rather than files

### How does Chorus decide when to rebuild?

A target is rebuilt if:
1. The target file doesn't exist
2. Any dependency is newer than the target
3. The target is marked phony
4. The target name starts with underscore
5. The target is named "all"

### Can I have circular dependencies?

No, circular dependencies are not allowed and will cause an error:
```yaml
targets:
  a:
    deps: [b]
  b:
    deps: [a]  # Error!
```

### How do I create a default target?

By convention, create an `all` target. It's executed when you run `chorus` without arguments:
```yaml
targets:
  all:
    deps: [program]
```

### Can I have multiple commands per target?

Yes, list them in the `cmds` array:
```yaml
targets:
  install:
    cmds:
      - "mkdir -p /usr/local/bin"
      - "cp program /usr/local/bin/"
      - "chmod +x /usr/local/bin/program"
```

### What happens if a command fails?

The build stops immediately, and Chorus exits with code 1. Subsequent commands are not executed.

## Usage

### How do I build my project?

Run `chorus` in the directory with `chorus.build`:
```bash
chorus          # Builds 'all' target
chorus program  # Builds 'program' target
chorus clean all # Builds 'clean' then 'all'
```

### Can I build multiple targets at once?

Yes, list them as arguments:
```bash
chorus clean all test
```

They execute in order, left to right.

### How do I do a clean build?

```bash
chorus clean all
```

Or manually:
```bash
rm -f *.o
chorus
```

### Can I do parallel builds?

No, Chorus currently doesn't support parallel builds. All targets build sequentially.

### How do I see what Chorus is doing?

Chorus displays colored output showing each target being built:
- Cyan `→` means executing
- Yellow `→` means up-to-date
- Commands' stdout/stderr are shown directly

### Can I do a dry run?

No, Chorus doesn't have a dry-run mode. To see what would build, check file timestamps manually or use a test directory.

### How do I force a rebuild?

Options:
1. Delete the target file: `rm program`
2. Touch a dependency: `touch src/main.c`
3. Use a phony target for the build step
4. Clean and rebuild: `chorus clean all`

## Troubleshooting

### "chorus.build: no such file or directory"

**Cause:** No `chorus.build` file in current directory.

**Solution:** Ensure you're in the correct directory or create a `chorus.build` file.

### "target 'X' undefined"

**Cause:** Referenced target doesn't exist.

**Solution:** Check spelling and ensure the target is defined in `chorus.build`.

### "yaml: line X: mapping values are not allowed"

**Cause:** YAML syntax error, usually indentation.

**Solution:** Check indentation uses spaces (not tabs) and colons are correctly placed.

### "command failed"

**Cause:** A shell command exited with non-zero status.

**Solution:** 
- Check the command's error output
- Test the command manually
- Verify file paths and variable expansions

### Build seems stuck

**Cause:** A command is waiting for input or running a long operation.

**Solution:** Press Ctrl+C to interrupt and investigate the command.

### Variables not expanding

**Cause:** Missing `${}` syntax or quoting issues.

**Solution:** Use `"${VAR}"` with quotes:
```yaml
# Wrong
cmds:
  - ${CC} -o program main.o

# Correct
cmds:
  - "${CC} -o program main.o"
```

### Automatic variables not working

**Cause:** Using automatic variables outside of commands.

**Solution:** Automatic variables (`${@}`, `${<}`, `${^}`) only work in target commands:
```yaml
targets:
  program:
    deps: [main.o]
    cmds:
      - "gcc -o ${@} ${^}"  # Correct
```

## Comparison with Other Tools

### Chorus vs Make

| Feature | Chorus | Make |
|---------|--------|------|
| Syntax | YAML | Makefile |
| Learning curve | Easy | Moderate |
| Parallel builds | No | Yes (-j) |
| Pattern rules | No | Yes |
| Conditional logic | Via shell | Built-in |
| Whitespace | Spaces | Tabs required |

### Chorus vs CMake

CMake is a meta-build system (generates Makefiles). Chorus is a direct build tool like Make. CMake is more complex but more powerful for large, cross-platform projects.

### Chorus vs Ninja

Ninja is focused on speed and parallel execution. Chorus prioritizes simplicity and readability over performance.

### When should I use Chorus?

Use Chorus when:
- You want simple, readable build files
- Your project is straightforward (C/C++, small to medium size)
- You're on Unix-like systems
- You don't need parallel builds
- You prefer YAML over Make syntax

Consider alternatives when:
- You need parallel builds for speed
- You need advanced features (pattern rules, complex conditionals)
- You need cross-platform support (Windows)
- Your project is very large or complex

## Advanced Topics

### Can I use Chorus for non-C projects?

Yes! Chorus works with any command-line tools:
- Rust: `rustc`, `cargo`
- Go: `go build`
- Java: `javac`
- Python: compilation, packaging
- Documentation: LaTeX, Markdown to HTML
- Any shell commands

### How do I integrate Chorus with CI/CD?

Run Chorus in your CI script:
```bash
#!/bin/bash
set -e

# Build
chorus all || exit 1

# Test
chorus test || exit 1

# Package
chorus package || exit 1
```

### Can I generate chorus.build files programmatically?

Yes, since it's YAML, you can generate it with any tool that outputs YAML:
```python
import yaml

config = {
    'variables': {'CC': 'gcc'},
    'targets': {
        'all': {'deps': ['program']},
        'program': {'cmds': ['gcc -o program main.c']}
    }
}

with open('chorus.build', 'w') as f:
    yaml.dump(config, f)
```

### How do I debug my chorus.build file?

1. Start simple and add complexity incrementally
2. Test each target individually
3. Add echo commands to see variable values:
   ```yaml
   cmds:
     - "echo 'CC=${CC}, CFLAGS=${CFLAGS}'"
     - "${CC} ${CFLAGS} -c main.c"
   ```
4. Check YAML syntax with a validator
5. Review [Syntax Guide](syntax.md) for common issues

## Getting Help

### Where can I find more examples?

- [Example directory](../example/) in the repository
- [NovariaOS build file](https://github.com/z3nnix/NovariaOS/blob/main/chorus.build)
- [Advanced Usage Guide](advanced.md)

### How do I report bugs?

Open an issue on GitHub: [z3nnix/chorus/issues](https://github.com/z3nnix/chorus/issues)

Include:
- Chorus version
- Operating system
- Your `chorus.build` file
- Error messages
- Steps to reproduce

### How can I contribute?

Contributions are welcome! Check the GitHub repository for:
- Open issues
- Feature requests
- Documentation improvements
- Bug fixes

### Is there a community?

Check the GitHub repository for discussions and issues. The project welcomes community involvement.
