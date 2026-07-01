# Building as a Shared Library

In addition to the `chorus` CLI, the build engine can be compiled as a
C-shared library (`.so` on Linux, `.dylib` on macOS), so it can be embedded
into other programs (C, C++, or anything with a C FFI).

## Building the Library

```bash
git clone https://github.com/z3nnix/chorus.git
cd chorus
CGO_ENABLED=1 go build -buildmode=c-shared -o libchorus.so ./cmd/libchorus
```

This produces two files:
- `libchorus.so` - the shared library
- `libchorus.h` - a generated C header with the exported function declarations

Building the shared library is entirely optional: `cmd/chorus` still builds
as a regular standalone binary and does not require `cgo` or a C toolchain.

## C API

```c
int  ChorusBuild(const char* config_path, const char* targets);
char* ChorusLastError(void);
void  ChorusFreeString(char* s);
```

### `ChorusBuild`

Loads a `chorus.build` file from `config_path` (pass `""` to use `chorus.build`
in the current directory) and builds the given `targets`, a space-separated
list of target names (pass `""` to build the default `all` target).

Returns `0` on success, `-1` on failure. On failure, call `ChorusLastError` to
retrieve the error message.

### `ChorusLastError`

Returns the message from the last failed `ChorusBuild` call, or an empty
string if the last call succeeded. The returned string is heap-allocated and
must be released with `ChorusFreeString`.

### `ChorusFreeString`

Frees a string previously returned by `ChorusLastError`.

## Example

```c
#include <stdio.h>
#include "libchorus.h"

int main(void) {
    if (ChorusBuild("chorus.build", "all") != 0) {
        char* err = ChorusLastError();
        fprintf(stderr, "build failed: %s\n", err);
        ChorusFreeString(err);
        return 1;
    }
    return 0;
}
```

Compile and link against the library:

```bash
gcc example.c -o example -L. -lchorus
LD_LIBRARY_PATH=. ./example
```

## Notes

- The library does not install signal handlers (unlike the CLI), so `Ctrl+C`
  handling remains under the control of the host process.
- `ChorusBuild` is not safe to call concurrently on the same process without
  external synchronization, since build progress is written to `stdout`.
