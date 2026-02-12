# Getting Started with Chorus

This guide will help you install Chorus and create your first build file.

## Installation

### Unix-like Systems (Linux, BSD, macOS)

Download and install the latest release:

```bash
wget https://github.com/z3nnix/chorus/releases/download/1.0.2/chorus && \
sudo mv chorus /usr/bin/chorus && \
sudo chmod +x /usr/bin/chorus
```

### Verify Installation

Check that Chorus is installed correctly:

```bash
chorus --help
```

### Building from Source

If you prefer to build from source:

```bash
git clone https://github.com/z3nnix/chorus.git
cd chorus
go build -o chorus ./cmd/chorus
sudo mv chorus /usr/bin/chorus
```

## Your First Build

Let's create a simple C project to demonstrate Chorus.

### Project Structure

Create a new directory for your project:

```bash
mkdir hello-chorus
cd hello-chorus
```

### Create Source Files

Create a simple C program:

**main.c:**
```c
#include <stdio.h>

int main() {
    printf("Hello from Chorus!\n");
    return 0;
}
```

### Create a Build File

Create a file named `chorus.build`:

```yaml
variables:
  CC: "gcc"
  CFLAGS: "-Wall -O2"

targets:
  all:
    deps: [hello]

  hello:
    deps: [main.o]
    cmds:
      - "${CC} -o ${@} ${^}"

  main.o:
    cmds:
      - "${CC} ${CFLAGS} -c main.c -o ${@}"

  clean:
    phony: true
    cmds:
      - "rm -f *.o hello"
```

### Build Your Project

Run Chorus to build your project:

```bash
chorus
```

Or specify a target explicitly:

```bash
chorus all
```

### Run Your Program

```bash
./hello
```

You should see:
```
Hello from Chorus!
```

### Clean Up

Remove build artifacts:

```bash
chorus clean
```

## Understanding the Build Process

When you run `chorus`, here's what happens:

1. **Load Configuration**: Chorus reads `chorus.build` from the current directory
2. **Parse Dependencies**: It builds a dependency graph from your targets
3. **Check Timestamps**: For each target, it checks if rebuilding is necessary
4. **Execute Commands**: Commands are executed in dependency order
5. **Report Results**: Chorus displays colored output showing progress

### Output Indicators

- `→` (cyan): Target is being built
- `→` (yellow): Target is up-to-date, skipping
- `✗` (red): Command failed
- `done` (green): Build completed successfully

## Next Steps

- Read the [Configuration Reference](configuration.md) to learn about all available options
- Explore [Build File Syntax](syntax.md) for detailed syntax information
- Check out [Advanced Usage](advanced.md) for complex build scenarios
- See the [example directory](../example/) for a more complete example

## Common Issues

### "chorus.build not found"

Make sure you're in a directory that contains a `chorus.build` file.

### Permission Denied

If you get permission errors when installing:
```bash
sudo chmod +x /usr/bin/chorus
```

### Command Not Found

Ensure `/usr/bin` is in your PATH:
```bash
echo $PATH
```

If not, add it to your shell configuration (~/.bashrc, ~/.zshrc, etc.):
```bash
export PATH="/usr/bin:$PATH"
```
