# PatchyDOS Dune II v1.01 — Go source

Build with Go 1.23+:

```sh
go test ./...
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -H windowsgui" -o PatchyDOS-Dune2-x64.exe .
GOOS=windows GOARCH=386 CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -H windowsgui" -o PatchyDOS-Dune2-x86.exe .
```

For the runtime layout, put the executables beside an `assets/` folder that contains `Dune2_Kor_v1.01.json`, `Dune2-Battle-for-Arrakis-Korean-v1.01.cue`, `patches/Dune2-Battle-for-Arrakis-Korean-v1.01.bin.xdelta`, `xdelta3.exe`, and `xdelta3-x64.exe` (as supplied in the distribution ZIP).

CLI for non-Windows systems: `go run . apply /path/to/original/disc/folder` (requires a matching decoder executable in the assets folder).

v1.01 execution safeguards: 5-minute total decoder timeout, 90-second no-activity watchdog,
5-second elapsed/output-size logging, stderr/stdout diagnostics, Windows GUI cancellation,
and temporary-output cleanup after failure or cancellation.

Go's source scan checks for SHA256+MD5+size; no user file is modified. See README.ko.md and verification/.
