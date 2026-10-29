# PatchyFM_Yuuwaku_v1.02 — FM TOWNS 《誘惑 / Yuuwaku》

《誘惑》(Yuuwaku / Temptation) FM TOWNS 원본 CD의 데이터 트랙 IMG/BIN에
사용자 제공 한국어 v1.02 xdelta를 적용해 검증된 BIN/CUE를 생성하는 Windows 패처입니다.

## v1.02 — 릴리스 빌드의 메모리 사용과 배포 안전성 개선

- 한글화 빌드에서 **스택 메모리 오류 수정**을 반영한 새 xdelta를 사용합니다. v1.01과 비교하면 CD 이미지 내 `TMB.COM`의 내용이 변경되며, 그 외 273개 수록 파일은 동일합니다.
- `TMB.COM`이 변경됐다는 사실과 패치의 복원 무결성은 확인했으나, 실제 게임 실행 중 스택 메모리 문제가 완전히 해결됐는지는 **실기 또는 에뮬레이터 테스트가 필요**합니다.
- 새 결과를 `Temptation-Korean-v1.02.bin`과 `Temptation-Korean-v1.02.cue`로 생성하여 기존 v1.0/v1.01 출력과 충돌하지 않도록 했습니다.
- 원본·패치·결과의 해시와 크기를 검증하고, 잘못된 결과물을 덮어쓰지 않습니다.
- 패처의 입력·출력 검사에 스트리밍 해시를 사용하며, 외부 xdelta3 진단 데이터의 메모리 점유를 **출력별 최대 64KiB**로 제한합니다. xdelta3 실행 **5분 제한**, 출력 파이프 종료 지연 시 전체 결과 해시 검증을 통해서만 정상 처리하는 안전 절차를 추가했습니다.

## 사용법

1. `PatchyFM_Yuuwaku_v1.02.zip`의 내용을 **한 폴더에 모두** 압축 해제합니다. EXE만 이동하면 `assets/` 검증에 실패합니다.
2. Windows 64비트는 `PatchyFM-Yuuwaku-x64.exe`, 32비트는 `PatchyFM-Yuuwaku-x86.exe`를 실행합니다.
3. **폴더 선택...**으로 원본 IMG/BIN이 있는 폴더를 지정합니다. 하위 폴더도 검색하며 CUE/CCD/SUB를 지정할 필요는 없습니다.
4. **한글패치 적용**을 누르면 원본 이미지가 있는 폴더에 결과 BIN/CUE가 생성됩니다.
5. Tsugaru 등에서 새로 생성한 CUE를 CD-ROM으로 사용합니다. 별도의 게임 User Disk가 필요할 수 있으며 포함하지 않습니다.

```text
Temptation-Korean-v1.02.bin
Temptation-Korean-v1.02.cue
```

원본 IMG/BIN/CCD/CUE/SUB는 수정하지 않습니다. 같은 원본 IMG/BIN이 여러 개 있으면 중단합니다. 이미 정상적인 BIN/CUE가 존재하면 재사용하며, 같은 이름의 **내용이 다른 파일은 덮어쓰지 않습니다**. 임시 `.bin.tmp`·`.cue.tmp` 파일은 작업 성공/실패 후 정리됩니다. 오래된 결과를 보관하려면 기존 파일은 따로 보관하십시오.

## 전체 무결성 검증값

| 대상 | 크기 (바이트) | MD5 | SHA-256 |
|---|---:|---|---|
| 원본 `Yuuwaku.img` | 6,703,200 | `000b8a4e832d952fb29ea6b831ed8cd0` | `09095579184db0d9e1de7b17071550090782e134058b1d93aa432409dfc1557b` |
| v1.02 `Temptation-Korean-v1.02.xdelta` | 330,776 | `74f5f1d6980bb68b7cedf1456baa22db` | `7201fa46c7603e29a7d9deb3cececc00836eee424546ea1463aba69fa8c4808a` |
| v1.02 `Temptation-Korean-v1.02.bin` | 7,599,312 | `6a83336e6b2b85f63567712e3144895d` | `746c23498bbac452c314ca4f07962dc65cd6c0dfcb2dd90441ee7a2de904b5f5` |
| v1.02 `Temptation-Korean-v1.02.cue` | 89 | `b524f19d88ecb7db5afe6b8ea160a356` | `eadcea111c967e4df85bb8a91b341a707a2b9bdade5a06dcb659ba7571ab02c4` |

v1.02 xdelta는 사용자가 제공한 `Temptation-Korean-v1.0(2).xdelta`의 **바이트를 변경하지 않고 이름만 변경**해 배포합니다. xdelta와 생성 BIN의 해시는 이전 버전과 다르며, CD 이미지의 변경 바이트 수는 v1.01 대비 **1,651바이트**입니다. 이 수에는 CD 섹터 검사·패리티 데이터 등도 포함됩니다.

## 검증 및 배포 정책

- 새 xdelta를 **실제 원본 `Yuuwaku.img`에 독립 디코더로 적용**하여 VCDIFF 윈도우 Adler-32 체크섬 `0x45ed3812`와 위 출력 BIN의 MD5·SHA-256을 검증했습니다.
- Go 단위 테스트, Windows x64/x86 크로스 빌드, 출력 파일 보호 및 외부 디코더 실패/시간 제한/임시 파일 정리 테스트를 수행했습니다.
- Windows GUI와 동봉된 Windows xdelta3 실행파일의 실제 동작, FM TOWNS 게임 플레이 중 메모리 오류 수정 여부는 **아직 별도 검증이 필요**합니다.
- 원본 게임 CD, User Disk, 완성 BIN은 **배포 ZIP과 소스 ZIP에 포함하지 않습니다**. 배포물에는 오직 한글패치 xdelta, 실행파일, 검증 매니페스트, CUE 템플릿, 라이선스·고지 문서만 포함됩니다.
- 패처 소스는 Apache License 2.0이며, 동봉 xdelta3 v3.2.0도 Apache 2.0입니다. 게임과 한글 번역 콘텐츠의 권리는 각각의 권리자에게 있습니다. **FreeTOWNSOS는 사용하거나 포함하지 않습니다.**

## 배포 파일 구성

```text
PatchyFM_Yuuwaku_v1.02/
├── PatchyFM-Yuuwaku-x64.exe
├── PatchyFM-Yuuwaku-x86.exe
├── README.ko.md
├── LICENSE
├── THIRD_PARTY_LICENSE.txt
├── NOTICE_MODIFICATIONS.txt
├── UPSTREAM_XDELTA.txt
├── SHA256SUMS.txt
└── assets/
    ├── Yuuwaku_Kor_v1.02.json
    ├── Temptation-Korean-v1.02.cue
    ├── xdelta3.exe
    ├── xdelta3-x64.exe
    └── patches/
        └── Temptation-Korean-v1.02.xdelta
```

## 저장소 소스와 빌드

- [PatchyFM Yuuwaku Go 소스](../src/patchyfm-yuuwaku/)
- 주요 엔진: `core.go`, Windows GUI: `main_windows.go`, CLI: `main_cli.go`
- 테스트: `core_test.go`, `release_safety_test.go`, `verification/verify_vcdiff_lzma.py`

```sh
cd src/patchyfm-yuuwaku
go test ./...
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -H windowsgui" -o PatchyFM-Yuuwaku-x64.exe .
GOOS=windows GOARCH=386 CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -H windowsgui" -o PatchyFM-Yuuwaku-x86.exe .
```

실행파일 옆에 배포 ZIP의 `assets/` 폴더를 함께 배치해야 합니다. 소스 트리에는 원본 게임, 전체 결과 BIN, xdelta 및 xdelta3 실행파일을 넣지 않습니다.

- 배포 ZIP SHA-256: `45ab5cf39088660d4c849b4028f6051032b2828252f6a2c520ff8bb3170d837b`
- 소스 ZIP SHA-256: `39db288ba5e6fc90b4f5b6ec38e3b84cbece346957caa0722cb9bd1c6466996f`
