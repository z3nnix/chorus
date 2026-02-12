# Advanced Usage

This guide covers advanced patterns, techniques, and best practices for using Chorus effectively in complex projects.

## Multi-Language Projects

### C and C++ Mixed Project

```yaml
variables:
  CC: "gcc"
  CXX: "g++"
  CFLAGS: "-Wall -Werror -O2 -Iinclude"
  CXXFLAGS: "-Wall -Werror -O2 -std=c++17 -Iinclude"
  LDFLAGS: "-lstdc++ -lm"

targets:
  all:
    deps: [program]

  program:
    deps: [main.o, utils.o, algorithm.o]
    cmds:
      - "${CXX} -o ${@} ${^} ${LDFLAGS}"

  # C source file
  utils.o:
    deps: [src/utils.c]
    cmds:
      - "${CC} ${CFLAGS} -c ${<} -o ${@}"

  # C++ source files
  main.o:
    deps: [src/main.cpp]
    cmds:
      - "${CXX} ${CXXFLAGS} -c ${<} -o ${@}"

  algorithm.o:
    deps: [src/algorithm.cpp]
    cmds:
      - "${CXX} ${CXXFLAGS} -c ${<} -o ${@}"

  clean:
    phony: true
    cmds:
      - "rm -f *.o program"
```

### Project with Assembly

```yaml
variables:
  AS: "nasm"
  ASFLAGS: "-f elf64"
  CC: "gcc"
  CFLAGS: "-Wall -O2"

targets:
  all:
    deps: [program]

  program:
    deps: [main.o, asm_routines.o]
    cmds:
      - "${CC} -o ${@} ${^}"

  main.o:
    deps: [src/main.c]
    cmds:
      - "${CC} ${CFLAGS} -c ${<} -o ${@}"

  asm_routines.o:
    deps: [src/asm_routines.asm]
    cmds:
      - "${AS} ${ASFLAGS} ${<} -o ${@}"
```

## Build Variants

### Debug vs Release Builds

Create separate targets for different build configurations:

```yaml
variables:
  CC: "gcc"
  BASE_FLAGS: "-Wall -Iinclude"
  DEBUG_FLAGS: "${BASE_FLAGS} -g -DDEBUG -O0"
  RELEASE_FLAGS: "${BASE_FLAGS} -O3 -DNDEBUG"

targets:
  all:
    deps: [release]

  # Debug build
  debug:
    deps: [main_debug.o, utils_debug.o]
    cmds:
      - "${CC} -o program_debug ${^}"
      - "echo 'Debug build ready'"

  main_debug.o:
    deps: [src/main.c]
    cmds:
      - "${CC} ${DEBUG_FLAGS} -c ${<} -o ${@}"

  utils_debug.o:
    deps: [src/utils.c]
    cmds:
      - "${CC} ${DEBUG_FLAGS} -c ${<} -o ${@}"

  # Release build
  release:
    deps: [main_release.o, utils_release.o]
    cmds:
      - "${CC} -o program ${^}"
      - "strip program"
      - "echo 'Release build ready'"

  main_release.o:
    deps: [src/main.c]
    cmds:
      - "${CC} ${RELEASE_FLAGS} -c ${<} -o ${@}"

  utils_release.o:
    deps: [src/utils.c]
    cmds:
      - "${CC} ${RELEASE_FLAGS} -c ${<} -o ${@}"

  clean:
    phony: true
    cmds:
      - "rm -f *.o program program_debug"
```

Usage:
```bash
chorus debug    # Build debug version
chorus release  # Build release version
chorus          # Build release (default)
```

## Complex Dependency Chains

### Static and Dynamic Libraries

```yaml
variables:
  CC: "gcc"
  AR: "ar"
  CFLAGS: "-Wall -O2 -fPIC -Iinclude"

targets:
  all:
    deps: [program, libmath.a, libmath.so]

  # Static library
  libmath.a:
    deps: [add.o, multiply.o, divide.o]
    cmds:
      - "${AR} rcs ${@} ${^}"
      - "echo 'Created static library'"

  # Dynamic library
  libmath.so:
    deps: [add.o, multiply.o, divide.o]
    cmds:
      - "${CC} -shared -o ${@} ${^}"
      - "echo 'Created shared library'"

  # Program using static library
  program:
    deps: [main.o, libmath.a]
    cmds:
      - "${CC} -o ${@} main.o -L. -lmath"

  # Object files
  main.o:
    deps: [src/main.c]
    cmds:
      - "${CC} ${CFLAGS} -c ${<} -o ${@}"

  add.o:
    deps: [src/math/add.c]
    cmds:
      - "${CC} ${CFLAGS} -c ${<} -o ${@}"

  multiply.o:
    deps: [src/math/multiply.c]
    cmds:
      - "${CC} ${CFLAGS} -c ${<} -o ${@}"

  divide.o:
    deps: [src/math/divide.c]
    cmds:
      - "${CC} ${CFLAGS} -c ${<} -o ${@}"
```

### Nested Dependencies

```yaml
targets:
  all:
    deps: [app]

  app:
    deps: [app.o, libcore.a, libutils.a]
    cmds:
      - "gcc -o ${@} app.o -L. -lcore -lutils"

  libcore.a:
    deps: [core1.o, core2.o, libbase.a]
    cmds:
      - "ar rcs ${@} core1.o core2.o"

  libutils.a:
    deps: [utils1.o, utils2.o]
    cmds:
      - "ar rcs ${@} ${^}"

  libbase.a:
    deps: [base1.o, base2.o]
    cmds:
      - "ar rcs ${@} ${^}"

  # Object files...
  app.o:
    cmds:
      - "gcc -c src/app.c -o ${@}"

  core1.o:
    cmds:
      - "gcc -c src/core/core1.c -o ${@}"
  # ... etc
```

## Code Generation

### Generating Source Files

```yaml
variables:
  GENERATOR: "python3 scripts/generate.py"

targets:
  all:
    deps: [program]

  program:
    deps: [main.o, generated.o]
    cmds:
      - "gcc -o ${@} ${^}"

  # Regular source
  main.o:
    deps: [src/main.c, generated.h]
    cmds:
      - "gcc -c src/main.c -o ${@}"

  # Generated source
  generated.o:
    deps: [generated.c]
    cmds:
      - "gcc -c ${<} -o ${@}"

  generated.c:
    deps: [config.yaml, scripts/generate.py]
    cmds:
      - "${GENERATOR} config.yaml generated.c"

  generated.h:
    deps: [config.yaml, scripts/generate.py]
    cmds:
      - "${GENERATOR} config.yaml generated.h --header"

  clean:
    phony: true
    cmds:
      - "rm -f *.o program generated.c generated.h"
```

### Version Header Generation

```yaml
variables:
  VERSION: "1.2.3"
  GIT_HASH: "$(git rev-parse --short HEAD)"

targets:
  all:
    deps: [program]

  program:
    deps: [main.o, version.o]
    cmds:
      - "gcc -o ${@} ${^}"

  version.o:
    deps: [version.c]
    cmds:
      - "gcc -c ${<} -o ${@}"

  version.c:
    phony: true
    cmds:
      - |
        cat > version.c << 'EOF'
        const char* get_version(void) {
            return "${VERSION}";
        }
        const char* get_git_hash(void) {
            return "${GIT_HASH}";
        }
        EOF
```

## Testing Integration

### Unit Tests

```yaml
variables:
  CC: "gcc"
  CFLAGS: "-Wall -Iinclude -Itests"
  TEST_LIBS: "-lcriterion"

targets:
  all:
    deps: [program]

  test:
    phony: true
    deps: [run_tests]

  run_tests:
    phony: true
    deps: [test_runner]
    cmds:
      - "./test_runner"

  test_runner:
    deps: [test_main.o, test_utils.o, utils.o]
    cmds:
      - "${CC} -o ${@} ${^} ${TEST_LIBS}"

  test_main.o:
    deps: [tests/test_main.c]
    cmds:
      - "${CC} ${CFLAGS} -c ${<} -o ${@}"

  test_utils.o:
    deps: [tests/test_utils.c]
    cmds:
      - "${CC} ${CFLAGS} -c ${<} -o ${@}"

  # Production code
  utils.o:
    deps: [src/utils.c]
    cmds:
      - "${CC} ${CFLAGS} -c ${<} -o ${@}"

  program:
    deps: [main.o, utils.o]
    cmds:
      - "${CC} -o ${@} ${^}"
```

### Coverage Analysis

```yaml
variables:
  CC: "gcc"
  CFLAGS: "-Wall --coverage"

targets:
  coverage:
    phony: true
    deps: [run_tests_coverage, generate_coverage_report]

  run_tests_coverage:
    phony: true
    deps: [test_runner_cov]
    cmds:
      - "./test_runner_cov"

  generate_coverage_report:
    phony: true
    cmds:
      - "gcov *.c"
      - "lcov --capture --directory . --output-file coverage.info"
      - "genhtml coverage.info --output-directory coverage_html"
      - "echo 'Coverage report: coverage_html/index.html'"

  test_runner_cov:
    deps: [test_main_cov.o, utils_cov.o]
    cmds:
      - "${CC} ${CFLAGS} -o ${@} ${^} -lcriterion"

  test_main_cov.o:
    deps: [tests/test_main.c]
    cmds:
      - "${CC} ${CFLAGS} -c ${<} -o ${@}"

  utils_cov.o:
    deps: [src/utils.c]
    cmds:
      - "${CC} ${CFLAGS} -c ${<} -o ${@}"

  clean_coverage:
    phony: true
    cmds:
      - "rm -f *.gcda *.gcno *.gcov coverage.info"
      - "rm -rf coverage_html"
```

## Build Organization

### Subdirectory Structure

```yaml
variables:
  CC: "gcc"
  CFLAGS: "-Wall -O2 -Iinclude"
  BUILD_DIR: "build"
  OBJ_DIR: "${BUILD_DIR}/obj"

targets:
  all:
    deps: [setup_dirs, program]

  setup_dirs:
    phony: true
    cmds:
      - "mkdir -p ${BUILD_DIR}"
      - "mkdir -p ${OBJ_DIR}"

  program:
    deps: ["${OBJ_DIR}/main.o", "${OBJ_DIR}/utils.o"]
    cmds:
      - "${CC} -o ${BUILD_DIR}/program ${^}"

  "${OBJ_DIR}/main.o":
    deps: [src/main.c]
    cmds:
      - "${CC} ${CFLAGS} -c src/main.c -o ${@}"

  "${OBJ_DIR}/utils.o":
    deps: [src/utils.c]
    cmds:
      - "${CC} ${CFLAGS} -c src/utils.c -o ${@}"

  clean:
    phony: true
    cmds:
      - "rm -rf ${BUILD_DIR}"
```

### Multiple Executables

```yaml
variables:
  CC: "gcc"
  CFLAGS: "-Wall -O2"

targets:
  all:
    deps: [server, client, tool]

  # Server program
  server:
    deps: [server_main.o, network.o, common.o]
    cmds:
      - "${CC} -o ${@} ${^} -lpthread"

  server_main.o:
    deps: [src/server/main.c]
    cmds:
      - "${CC} ${CFLAGS} -c ${<} -o ${@}"

  # Client program
  client:
    deps: [client_main.o, network.o, common.o]
    cmds:
      - "${CC} -o ${@} ${^}"

  client_main.o:
    deps: [src/client/main.c]
    cmds:
      - "${CC} ${CFLAGS} -c ${<} -o ${@}"

  # Utility tool
  tool:
    deps: [tool_main.o, common.o]
    cmds:
      - "${CC} -o ${@} ${^}"

  tool_main.o:
    deps: [src/tool/main.c]
    cmds:
      - "${CC} ${CFLAGS} -c ${<} -o ${@}"

  # Shared objects
  network.o:
    deps: [src/shared/network.c]
    cmds:
      - "${CC} ${CFLAGS} -c ${<} -o ${@}"

  common.o:
    deps: [src/shared/common.c]
    cmds:
      - "${CC} ${CFLAGS} -c ${<} -o ${@}"

  clean:
    phony: true
    cmds:
      - "rm -f *.o server client tool"
```

## Platform-Specific Builds

### Detecting Platform

```yaml
variables:
  OS: "$(uname -s)"
  ARCH: "$(uname -m)"
  CC: "gcc"
  CFLAGS: "-Wall -O2"

targets:
  all:
    deps: [program]

  program:
    deps: [main.o, platform.o]
    cmds:
      - |
        case "${OS}" in
          Linux*)   ${CC} -o ${@} ${^} -lpthread -lrt ;;
          Darwin*)  ${CC} -o ${@} ${^} -lpthread ;;
          FreeBSD*) ${CC} -o ${@} ${^} -lpthread -lutil ;;
          *)        echo "Unsupported OS: ${OS}"; exit 1 ;;
        esac

  main.o:
    cmds:
      - "${CC} ${CFLAGS} -DOS_${OS} -c src/main.c -o ${@}"
```

## Documentation Generation

### Building with Doxygen

```yaml
targets:
  all:
    deps: [program, docs]

  docs:
    phony: true
    deps: [html_docs, pdf_docs]

  html_docs:
    deps: [Doxyfile, src/main.c, src/utils.c]
    cmds:
      - "doxygen Doxyfile"
      - "echo 'HTML docs: docs/html/index.html'"

  pdf_docs:
    deps: [html_docs]
    cmds:
      - "cd docs/latex && make"
      - "cp docs/latex/refman.pdf docs/manual.pdf"
      - "echo 'PDF docs: docs/manual.pdf'"

  clean_docs:
    phony: true
    cmds:
      - "rm -rf docs/html docs/latex docs/manual.pdf"
```

## Advanced Patterns

### Conditional Execution

```yaml
targets:
  conditional_build:
    phony: true
    cmds:
      - |
        if [ -f .use_clang ]; then
          export CC=clang
        else
          export CC=gcc
        fi
        $CC -o program src/*.c
```

### Dependency on External Files

```yaml
targets:
  config_dependent:
    deps: [.config, src/main.c]
    cmds:
      - "gcc -DCONFIG=$(cat .config) -c src/main.c -o main.o"
```

### Pre and Post Build Hooks

```yaml
targets:
  all:
    deps: [pre_build, program, post_build]

  pre_build:
    phony: true
    cmds:
      - "echo 'Starting build at ${DATE}'"
      - "scripts/check_dependencies.sh"

  program:
    deps: [main.o, utils.o]
    cmds:
      - "gcc -o ${@} ${^}"

  post_build:
    phony: true
    deps: [program]
    cmds:
      - "scripts/sign_binary.sh program"
      - "echo 'Build completed successfully'"
```

### Resource Embedding

```yaml
targets:
  program:
    deps: [main.o, resources.o]
    cmds:
      - "gcc -o ${@} ${^}"

  resources.o:
    deps: [resources.bin]
    cmds:
      - "objcopy -I binary -O elf64-x86-64 -B i386 ${<} ${@}"

  resources.bin:
    deps: [data/icon.png, data/config.json]
    cmds:
      - "tar czf ${@} -C data ."
```

## Performance Optimization

### Minimizing Rebuilds

Structure dependencies carefully to minimize unnecessary rebuilds:

```yaml
# Good: Separate interface from implementation
targets:
  main.o:
    deps: [src/main.c, include/utils.h]  # Only header
    cmds:
      - "gcc -c src/main.c -o ${@}"

  utils.o:
    deps: [src/utils.c, include/utils.h]
    cmds:
      - "gcc -c src/utils.c -o ${@}"

# Avoid: Including implementation as dependency
# This causes main.o to rebuild when utils.c changes
targets:
  main.o:
    deps: [src/main.c, src/utils.c]  # Bad!
    cmds:
      - "gcc -c src/main.c -o ${@}"
```

### Build Time Tracking

```yaml
targets:
  timed_build:
    phony: true
    cmds:
      - "echo 'Build started at $(date)'"
      - "START=$(date +%s); chorus all; END=$(date +%s); echo \"Build took $((END-START))s\""
```

## Best Practices Summary

1. **Keep it simple**: Don't over-complicate your build files
2. **Use variables**: Define repeated values once
3. **Mark phony targets**: Always mark non-file targets as phony
4. **Minimize dependencies**: Only list actual dependencies
5. **Organize logically**: Group related targets together
6. **Add comments**: Document complex build steps
7. **Test incrementally**: Ensure incremental builds work correctly
8. **Keep commands simple**: One logical operation per command
9. **Use meaningful names**: Target names should be clear
10. **Validate early**: Test your build file with small changes first
