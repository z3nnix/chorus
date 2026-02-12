# Chorus Documentation

Welcome to the comprehensive documentation for **Chorus**, a minimalistic build system written in Go, designed for simplicity and efficiency.

## Table of Contents

- [Getting Started](getting-started.md) - Installation, quick start, and first build
- [Configuration Reference](configuration.md) - Complete reference for chorus.build files
- [Build File Syntax](syntax.md) - Detailed syntax guide for writing build files
- [CLI Reference](cli.md) - Command-line interface documentation
- [Advanced Usage](advanced.md) - Advanced patterns and best practices
- [FAQ](faq.md) - Frequently asked questions

## What is Chorus?

Chorus is an alternative build system that provides a simple, declarative way to define build targets and their dependencies. It uses YAML-based configuration files (`chorus.build`) to describe your build process.

### Key Features

- **Simple Syntax**: YAML-based configuration that's easy to read and write
- **Dependency Management**: Automatic dependency tracking and incremental builds
- **Variable Expansion**: Support for user-defined variables and automatic variables
- **Phony Targets**: Support for tasks that don't produce output files
- **Cross-Platform**: Designed for Unix-like systems (Linux, BSD, macOS)
- **Fast**: Written in Go for performance and reliability
- **Minimal**: No unnecessary features, just what you need to build

### Quick Example

```yaml
variables:
  CC: "gcc"
  CFLAGS: "-c -Wall"

targets:
  all:
    deps: [program]

  program:
    deps: [main.o]
    cmds:
      - "${CC} -o ${@} ${^}"

  main.o:
    cmds:
      - "${CC} ${CFLAGS} main.c -o ${@}"
```

## Projects Using Chorus

- [NovariaOS](https://github.com/novariaos) - An operating system project
- [Perano lang](https://github.com/noxzion/perano-lang) - A programming language

## Community and Support

- **GitHub**: [z3nnix/chorus](https://github.com/z3nnix/chorus)
- **Issues**: Report bugs or request features on GitHub Issues
- **License**: Check the LICENSE file in the repository

## Contributing

Contributions are welcome! Please check the GitHub repository for guidelines.
