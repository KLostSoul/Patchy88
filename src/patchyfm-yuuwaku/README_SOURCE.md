# PatchyFM Yuuwaku v1.0 — source

Windows GUI patcher for FM TOWNS Yuuwaku, developed from the previous PatchFM Alltynex approach.

Build (Go 1.23+):

```sh
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -H windowsgui" -o PatchyFM-Yuuwaku-x64.exe .
GOOS=windows GOARCH=386 CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -H windowsgui" -o PatchyFM-Yuuwaku-x86.exe .
go test ./...
```

The executable must run beside `assets/`, which contains the supplied user xdelta patch, expected Korean CUE template, manifest, and separately licensed xdelta3 executables.

- `core.go`: metadata validation, recursive automatic source IMG/BIN discovery by hash, safely verified xdelta3 workflow, exact Korean BIN/CUE output.
- `main_windows.go`: native Win32 GUI with a Windows shell folder picker; no individual image file selection.
- `main_cli.go`: non-Windows CLI (`apply ORIGINAL_FOLDER`) for source debugging; requires a compatible native xdelta3 helper at the configured path.
- `core_test.go`: unit tests.

Generated binaries were cross-compiled and 8 folder-mode integration tests passed. The Windows GUI was not interactively run in the Linux build environment.

Apache-2.0 applies to this PatchyFM utility code only. It does not license the original game, Korean translation content, or supplied xdelta patch. xdelta3 is separately distributed under Apache-2.0.
