# Chorus
> Its an alternative build system that writes on GO. Btw, for the NovariaOS

# Setup
Unix-like systems(Linux-based, BSD-like, macos maybe.. and etc.) only
```sh
curl https://raw.githubusercontent.com/z3nnix/chorus/refs/heads/main/install.sh | sh
```

# Building as a shared library
Chorus can optionally be built as a C-shared library (`.so`) for embedding
into other programs, without giving up the CLI. See
[docs/library.md](docs/library.md).

# Project, that used Chorus
- [NovariaOS](https://github.com/novariaos) <br>
- [Perano lang](https://github.com/noxzion/perano-lang) <br>
_No more.._

# Example of chorus.build file
[From NovariaOS](https://github.com/z3nnix/NovariaOS/blob/main/chorus.build)