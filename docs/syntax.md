# Build File Syntax

This guide provides detailed information about the syntax and semantics of `chorus.build` files.

## File Structure

A `chorus.build` file is written in YAML and consists of two main sections:

```yaml
variables:
  # Variable declarations

targets:
  # Target definitions
```

## YAML Basics

Chorus uses standard YAML syntax. Here are the key concepts:

### Indentation

YAML uses indentation (spaces, not tabs) to represent structure:

```yaml
targets:
  target1:
    deps: [dep1, dep2]
    cmds:
      - "command 1"
      - "command 2"
```

**Important:** Always use spaces for indentation, not tabs. Use 2 spaces per indentation level.

### Strings

Strings can be quoted or unquoted:

```yaml
variables:
  UNQUOTED: value
  SINGLE_QUOTED: 'value with spaces'
  DOUBLE_QUOTED: "value with ${variable} expansion"
```

Use quotes when:
- Value contains special characters
- Value starts with special YAML characters (-, [, {, etc.)
- You want to preserve exact spacing

### Lists

Lists can be written in two styles:

**Flow style (inline):**
```yaml
deps: [file1.o, file2.o, file3.o]
```

**Block style:**
```yaml
deps:
  - file1.o
  - file2.o
  - file3.o
```

Both are equivalent. Use flow style for short lists, block style for longer ones.

### Comments

Comments start with `#`:

```yaml
# This is a comment
variables:
  CC: "gcc"  # Inline comment
```

## Variables

### Declaration

Variables are declared in the `variables` section:

```yaml
variables:
  VARIABLE_NAME: "value"
```

### Naming Rules

- Variable names are case-sensitive
- Can contain letters, numbers, and underscores
- Convention: Use UPPERCASE for user-defined variables

**Good:**
```yaml
variables:
  CC: "gcc"
  CFLAGS: "-Wall"
  BUILD_DIR: "build"
  OUTPUT_FILE: "program"
```

**Avoid:**
```yaml
variables:
  cc: "gcc"           # Lowercase (not conventional)
  c-flags: "-Wall"    # Contains hyphen
  my var: "value"     # Contains space
```

### Variable References

Reference variables using `${VARIABLE_NAME}`:

```yaml
targets:
  main.o:
    cmds:
      - "${CC} ${CFLAGS} -c main.c"
```

**Nested expansion:**
Variables can reference other variables:

```yaml
variables:
  BASE_FLAGS: "-Wall -Werror"
  CFLAGS: "${BASE_FLAGS} -O2"
  DEBUG_FLAGS: "${BASE_FLAGS} -g"
```

### Automatic Variables

Chorus provides these automatic variables in target commands:

| Variable | Meaning | When to Use |
|----------|---------|-------------|
| `${@}` | Current target name | Output file name |
| `${<}` | First dependency | Single input file |
| `${^}` | All dependencies (space-separated) | Multiple input files |

**Example:**
```yaml
targets:
  program:
    deps: [main.o, util.o, math.o]
    cmds:
      - "gcc -o ${@} ${^}"
      # Expands to: gcc -o program main.o util.o math.o
```

## Targets

### Target Names

Target names can be:
- Filenames: `main.o`, `program`, `libmath.a`
- Task names: `all`, `clean`, `install`, `test`
- Any string without spaces

**Special behaviors:**
- Targets starting with `_` are always considered out-of-date
- Target `all` is the default when no target is specified

### Dependencies

Dependencies are specified using the `deps` field:

```yaml
targets:
  program:
    deps: [main.o, utils.o]
```

**Dependency resolution order:**
1. Chorus builds dependencies before the target
2. If a dependency is out-of-date, it's rebuilt first
3. If any dependency changes, the target is rebuilt

**Circular dependencies:**
Circular dependencies are not allowed and will cause an error:

```yaml
targets:
  a:
    deps: [b]
  b:
    deps: [a]  # Error: circular dependency
```

### Commands

Commands are shell commands executed to build the target:

```yaml
targets:
  main.o:
    cmds:
      - "gcc -c main.c -o main.o"
      - "echo 'Compiled main.o'"
```

**Important details:**
- Each command runs in a separate shell (`sh -c`)
- Commands execute in order
- If a command fails (non-zero exit), the build stops
- Commands have access to environment variables

**Multiple commands:**
```yaml
targets:
  install:
    cmds:
      - "mkdir -p /usr/local/bin"
      - "cp program /usr/local/bin/"
      - "chmod +x /usr/local/bin/program"
```

**Long commands:**
Use YAML multi-line strings for readability:

```yaml
targets:
  complex:
    cmds:
      - |
        gcc -c main.c \
          -Wall -Werror \
          -O2 -g \
          -Iinclude \
          -o main.o
```

### Phony Targets

Phony targets don't represent files:

```yaml
targets:
  clean:
    phony: true
    cmds:
      - "rm -rf *.o program"
```

**Characteristics:**
- Always executed (never skipped as up-to-date)
- Don't check file timestamps
- Useful for: `clean`, `install`, `test`, `run`, `help`

**Example phony targets:**
```yaml
targets:
  all:
    deps: [build, test]

  build:
    deps: [program]

  test:
    phony: true
    deps: [program]
    cmds:
      - "./program --test"

  clean:
    phony: true
    cmds:
      - "rm -rf *.o program"

  install:
    phony: true
    deps: [program]
    cmds:
      - "cp program /usr/local/bin/"
```

## Incremental Builds

Chorus automatically determines when to rebuild targets based on timestamps.

### Rebuild Rules

A target is rebuilt if:
1. The target file doesn't exist
2. Any dependency is newer than the target
3. The target is phony
4. The target name starts with `_`
5. The target is `all`

### Timestamp Comparison

```yaml
targets:
  program:
    deps: [main.o, utils.o]
    cmds:
      - "gcc -o ${@} ${^}"
```

Chorus checks:
- Does `program` exist?
- Is `main.o` newer than `program`?
- Is `utils.o` newer than `program`?

If any answer is yes, `program` is rebuilt.

### Force Rebuild

To force a rebuild:
1. Use a phony target
2. Delete the target file
3. Touch a dependency file
4. Use underscore-prefixed target name

## Complete Syntax Example

Here's a complete example showing all syntax features:

```yaml
# Build configuration for MyProject
# Compiler: GCC 11+

variables:
  # Compiler settings
  CC: "gcc"
  CXX: "g++"
  
  # Flags
  COMMON_FLAGS: "-Wall -Werror"
  CFLAGS: "${COMMON_FLAGS} -std=c11"
  CXXFLAGS: "${COMMON_FLAGS} -std=c++17"
  LDFLAGS: "-lm -lpthread"
  
  # Directories
  SRC_DIR: "src"
  BUILD_DIR: "build"
  INC_DIR: "include"
  
  # Other
  VERSION: "2.1.0"
  APP_NAME: "myapp"

targets:
  # Default target
  all:
    deps: [program, tests]

  # Main program
  program:
    deps: [main.o, utils.o, math.o]
    cmds:
      - "${CC} -o ${APP_NAME} ${^} ${LDFLAGS}"
      - "echo 'Built ${APP_NAME} v${VERSION}'"

  # Object files
  main.o:
    deps: ["${SRC_DIR}/main.c"]
    cmds:
      - "${CC} ${CFLAGS} -I${INC_DIR} -c ${<} -o ${@}"

  utils.o:
    deps: ["${SRC_DIR}/utils.c"]
    cmds:
      - "${CC} ${CFLAGS} -I${INC_DIR} -c ${<} -o ${@}"

  math.o:
    deps: ["${SRC_DIR}/math.c"]
    cmds:
      - "${CC} ${CFLAGS} -I${INC_DIR} -c ${<} -o ${@}"

  # Test binary
  tests:
    phony: true
    deps: [test_runner]
    cmds:
      - "./${<}"

  test_runner:
    deps: [test_main.o, utils.o, math.o]
    cmds:
      - "${CC} -o ${@} ${^} ${LDFLAGS}"

  test_main.o:
    deps: ["${SRC_DIR}/test_main.c"]
    cmds:
      - "${CC} ${CFLAGS} -I${INC_DIR} -c ${<} -o ${@}"

  # Utility targets
  clean:
    phony: true
    cmds:
      - "rm -f *.o"
      - "rm -f ${APP_NAME} test_runner"
      - "echo 'Cleaned build artifacts'"

  install:
    phony: true
    deps: [program]
    cmds:
      - "mkdir -p /usr/local/bin"
      - "cp ${APP_NAME} /usr/local/bin/"
      - "echo 'Installed to /usr/local/bin/${APP_NAME}'"

  uninstall:
    phony: true
    cmds:
      - "rm -f /usr/local/bin/${APP_NAME}"
      - "echo 'Uninstalled ${APP_NAME}'"

  # Force rebuild
  _rebuild:
    deps: [clean, all]

  # Version info
  version:
    phony: true
    cmds:
      - "echo 'Version: ${VERSION}'"
      - "echo 'Build Date: ${DATE}'"
```

## Syntax Validation

Common syntax errors and how to fix them:

### Error: "yaml: line X: mapping values are not allowed in this context"

**Cause:** Incorrect indentation or missing colon

**Fix:**
```yaml
# Wrong
targets:
target1:
  cmds:

# Correct
targets:
  target1:
    cmds:
```

### Error: "target 'X' undefined"

**Cause:** Dependency references non-existent target

**Fix:** Ensure all dependencies are defined:
```yaml
targets:
  program:
    deps: [main.o]  # main.o must be defined

  main.o:  # Define it here
    cmds:
      - "gcc -c main.c -o main.o"
```

### Error: "found character that cannot start any token"

**Cause:** Special characters not quoted

**Fix:**
```yaml
# Wrong
cmds:
  - echo ${VAR}

# Correct
cmds:
  - "echo ${VAR}"
```

## Style Guide

### Recommended style:

```yaml
variables:
  CC: "gcc"
  CFLAGS: "-Wall -O2"

targets:
  all:
    deps: [program]

  program:
    deps: [main.o, utils.o]
    cmds:
      - "${CC} -o ${@} ${^}"

  main.o:
    cmds:
      - "${CC} ${CFLAGS} -c main.c -o ${@}"

  clean:
    phony: true
    cmds:
      - "rm -f *.o program"
```

**Guidelines:**
- 2-space indentation
- Blank line between targets
- Quote string values
- Use flow style for short dependency lists
- Add comments for complex sections
- Group related targets together
