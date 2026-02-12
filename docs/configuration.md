# Configuration Reference

This document provides a complete reference for `chorus.build` configuration files.

## File Format

Chorus uses YAML format for build configuration. The file must be named `chorus.build` and placed in the root of your project.

## Top-Level Structure

```yaml
variables:
  # Variable definitions

targets:
  # Target definitions
```

## Variables Section

Variables define reusable values that can be referenced in commands.

### Syntax

```yaml
variables:
  VARIABLE_NAME: "value"
  ANOTHER_VAR: "another value"
```

### User-Defined Variables

You can define any variables you need:

```yaml
variables:
  CC: "gcc"
  CXX: "g++"
  CFLAGS: "-Wall -Werror -O2"
  LDFLAGS: "-lm -lpthread"
  OUTPUT_DIR: "build"
  VERSION: "1.0.0"
```

### Built-in Variables

Chorus automatically provides these variables:

| Variable | Description | Example |
|----------|-------------|---------|
| `DATE` | Current date in YYYY-MM-DD format | `2024-01-15` |

### Variable Expansion

Variables are expanded in commands using `${VARIABLE_NAME}` syntax:

```yaml
targets:
  main.o:
    cmds:
      - "${CC} ${CFLAGS} -c main.c -o ${@}"
```

## Targets Section

Targets define the build units and their dependencies.

### Basic Target Syntax

```yaml
targets:
  target_name:
    deps: [dependency1, dependency2]
    cmds:
      - "command 1"
      - "command 2"
    phony: false
```

### Target Properties

#### `deps` (optional)

List of dependencies that must be built before this target.

```yaml
targets:
  program:
    deps: [main.o, utils.o, config.o]
    cmds:
      - "gcc -o program ${^}"
```

**Type:** Array of strings  
**Default:** Empty array

#### `cmds` (optional)

List of shell commands to execute when building this target.

```yaml
targets:
  clean:
    cmds:
      - "rm -rf *.o"
      - "rm -rf build/"
```

**Type:** Array of strings  
**Default:** Empty array

**Note:** Each command is executed in a separate shell (`sh -c`).

#### `phony` (optional)

Marks a target as phony (not a real file). Phony targets are always executed.

```yaml
targets:
  clean:
    phony: true
    cmds:
      - "rm -rf *.o program"
```

**Type:** Boolean  
**Default:** `false`

### Automatic Variables

Chorus provides automatic variables that can be used in commands:

| Variable | Description | Example |
|----------|-------------|---------|
| `${@}` | Current target name | `main.o` |
| `${<}` | First dependency | `main.c` |
| `${^}` | All dependencies (space-separated) | `main.o utils.o config.o` |

#### Example Usage

```yaml
targets:
  program:
    deps: [main.o, utils.o]
    cmds:
      - "gcc -o ${@} ${^}"
      # Expands to: gcc -o program main.o utils.o

  main.o:
    deps: [main.c]
    cmds:
      - "gcc -c ${<} -o ${@}"
      # Expands to: gcc -c main.c -o main.o
```

### Special Targets

#### `all` Target

By convention, `all` is the default target executed when no target is specified:

```yaml
targets:
  all:
    deps: [program, docs, tests]
```

When you run just `chorus`, it executes the `all` target.

#### Underscore-Prefixed Targets

Targets starting with underscore are treated as always out-of-date:

```yaml
targets:
  _rebuild:
    deps: [clean, all]
```

These are useful for forcing rebuilds or running tasks.

## Complete Example

Here's a comprehensive example demonstrating all features:

```yaml
variables:
  CC: "gcc"
  CXX: "g++"
  CFLAGS: "-Wall -Werror -O2 -Iinclude"
  CXXFLAGS: "-Wall -Werror -O2 -std=c++17 -Iinclude"
  LDFLAGS: "-lm -lpthread"
  BUILD_DIR: "build"
  VERSION: "1.2.3"

targets:
  # Default target
  all:
    deps: [program, library]

  # Main program
  program:
    deps: [main.o, utils.o, libmath.a]
    cmds:
      - "${CC} -o ${@} main.o utils.o -L. -lmath ${LDFLAGS}"

  # Object files
  main.o:
    deps: [src/main.c]
    cmds:
      - "${CC} ${CFLAGS} -c ${<} -o ${@}"

  utils.o:
    deps: [src/utils.c]
    cmds:
      - "${CC} ${CFLAGS} -c ${<} -o ${@}"

  math.o:
    deps: [src/math.c]
    cmds:
      - "${CC} ${CFLAGS} -c ${<} -o ${@}"

  # Static library
  libmath.a:
    deps: [math.o]
    cmds:
      - "ar rcs ${@} ${^}"

  # Phony targets
  clean:
    phony: true
    cmds:
      - "rm -f *.o *.a program"
      - "rm -rf ${BUILD_DIR}"

  install:
    phony: true
    deps: [all]
    cmds:
      - "mkdir -p /usr/local/bin"
      - "cp program /usr/local/bin/"

  test:
    phony: true
    deps: [program]
    cmds:
      - "./program --test"

  version:
    phony: true
    cmds:
      - "echo 'Version: ${VERSION}'"
      - "echo 'Build Date: ${DATE}'"
```

## Best Practices

### Variable Naming

- Use UPPERCASE for variables: `CC`, `CFLAGS`, `BUILD_DIR`
- Use descriptive names: `OUTPUT_DIR` instead of `OUT`
- Group related variables together

### Target Organization

- Put the default `all` target first
- Group related targets together
- Use comments to separate sections

### Phony Targets

Always mark non-file targets as phony:
- `clean`, `install`, `test`, `run`
- `help`, `version`, `docs`

### Command Best Practices

- One logical operation per command
- Use variables for repeated values
- Check for errors explicitly when needed:
  ```yaml
  cmds:
    - "test -f input.txt || { echo 'Missing input.txt'; exit 1; }"
  ```

## Validation

Chorus validates your configuration on load:

- YAML syntax must be correct
- Target names must be unique
- Circular dependencies are not allowed
- Referenced dependencies must exist

Common validation errors:

```
parse config: yaml: line 5: found character that cannot start any token
```
Solution: Check YAML syntax, proper indentation

```
target 'program' undefined
```
Solution: Ensure the target is defined in the `targets` section
