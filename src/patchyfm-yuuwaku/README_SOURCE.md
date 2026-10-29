# PatchyFM Yuuwaku v1.02 — build and verification

Standalone Windows patcher for FM TOWNS Yuuwaku, written in Go 1.23+.

## Build

```sh
go test ./...
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -H windowsgui" -o PatchyFM-Yuuwaku-x64.exe .
GOOS=windows GOARCH=386 CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -H windowsgui" -o PatchyFM-Yuuwaku-x86.exe .
```

`assets/` from the binary release is required beside the executables. It contains the supplied xdelta, the template CUE, integrity manifest and the externally licensed xdelta3 helpers. These payload binaries are deliberately **not** duplicated inside the source ZIP or source Git tree.

## Memory and deployment safety

- `hashFile()` computes MD5 and SHA-256 by streaming with `io.Copy`; whole CD images are **not** held in patcher memory.
- xdelta stdout/stderr are capped to the most recent 64 KiB each to avoid unbounded diagnostic buffering. Important previous output might be omitted if excessive diagnostics are emitted.
- External decoding has a 5-minute wall-clock deadline, 3-second `WaitDelay`, and differentiated timeout / nonzero-exit diagnostics.
- The `WaitDelay` workaround does not blindly accept an exit status: only a complete expected BIN matching size, MD5 and SHA-256 can be accepted.
- Input, patch, decoder helper and CUE template hashes are checked. Existing non-matching outputs are refused without overwrite; newly created temporary files are cleaned on failure.
- The on-disk output is named v1.02, intentionally separate from v1.01.

## Tests and scope

`release_safety_test.go` covers bounded diagnostics, safety timeouts, failure cleanup, output verification, overwrite protection and idempotence. `core_test.go` covers source discovery and hash checking. `verification/verify_vcdiff_lzma.py` independently replays the actual user xdelta, **using Python memory buffers**; this development-only script is not used by the shipping patcher.

The patch comes from user-supplied `Temptation-Korean-v1.0(2).xdelta`, SHA-256 `7201fa46c7603e29a7d9deb3cececc00836eee424546ea1463aba69fa8c4808a` (330,776 bytes). The expected BIN is 7,599,312 bytes, SHA-256 `746c23498bbac452c314ca4f07962dc65cd6c0dfcb2dd90441ee7a2de904b5f5`.

The emulator/game runtime memory fault is not verified as fixed by these patch-integrity tests. Windows GUI and packaged xdelta3 helper runtime were not executed in the Linux cross-build environment.

Apache-2.0 applies to the PatchyFM utility code; xdelta3 v3.2.0 is separately distributed under Apache-2.0. Original game data and translation patch rights are separate.
