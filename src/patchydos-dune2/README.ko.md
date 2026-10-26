# PatchyDOS_Dune2_v1.01

MS-DOS판 **《Dune II: The Battle for Arrakis》 CD-ROM**용 한글패치 적용 도구입니다.
이 패처는 사용자가 준비한 원본 CD 이미지의 데이터 트랙에 제공받은
`Dune2-Battle-for-Arrakis-Korean-v1.01.bin.xdelta`를 적용하여 **한글판 BIN/CUE**를 생성합니다.

## Windows 사용법

1. `PatchyDOS_Dune2_v1.01.zip`을 폴더에 모두 압축 해제합니다. EXE만 이동하면 동작하지 않습니다.
2. Windows 64비트에서는 `PatchyDOS-Dune2-x64.exe`, 32비트에서는 `PatchyDOS-Dune2-x86.exe`를 실행합니다.
3. **폴더 선택...**으로 원본 CD 이미지가 있는 폴더를 선택합니다.
4. 해당 폴더와 하위 폴더에서 원본 **IMG/BIN** 파일을 자동 인식합니다. CUE/CCD/SUB를 개별 지정할 필요가 없습니다.
5. **한글패치 적용**을 누르면 원본 이미지가 들어 있는 폴더에 아래 파일을 생성합니다.

```text
Dune2-Battle-for-Arrakis-Korean-v1.01.bin
Dune2-Battle-for-Arrakis-Korean-v1.01.cue
```

에뮬레이터에서 사용할 때는 **새로 생성한 CUE**를 불러옵니다. 원본 CCD/SUB는 새로 생성된 이미지의 섹터 수와 다르므로 함께 사용하지 마십시오.

## 지원 원본

사용자 제공 `DUNE2CD(1).img`와 같은 내용의 원본 CD **MODE1/2352 데이터 트랙**만 지원합니다.
파일명은 달라도 되지만 전체 크기·MD5·SHA-256이 다음과 정확히 일치해야 합니다.

- 크기: `15,664,320 bytes` (6,660 sectors)
- MD5: `7ed518c967d786fc9a47ef7ab86aab97`
- SHA-256: `edcdce3de2e4d986421d99cc11be424277317fc3e8cd71cf78f8ed544bce81bf`

파일명이 동일한 다른 발매판, ISO 추출판(2,048바이트 섹터), 개조본은 지원하지 않습니다.

## 패치와 결과 무결성

패치: `Dune2-Battle-for-Arrakis-Korean-v1.01.bin.xdelta`
- 크기: `356,319 bytes`
- SHA-256: `41ef92e9a9cf062109371abce1e591fbe331e7b184b6614f6eedc1f74acef9db`

생성 BIN:
- 크기: `16,835,616 bytes` (7,158 sectors)
- MD5: `1cd47830e0925747622bb53e12590ae2`
- SHA-256: `5f3877d5d6be7b54b525ec91df19a46b01c4a6f5e83e3df6b683f05203053752`

생성 CUE:
- 크기: `103 bytes`
- SHA-256: `3d26796ff3eefd265b4f85f80a0ac6bf0ee95ae1f4f0c593fb93c8e009b3834b`

### 안정성 정책

- 원본 이미지(IMG/CCD/CUE/SUB)는 수정하지 않습니다.
- 폴더/하위 폴더를 검색하고 파일명과 관계없이 전체 해시로 자동 인식합니다.
- 정확한 원본이 여러 개 검출되면 대상이 모호하므로 적용을 중단합니다.
- 패치, CUE 템플릿, 디코더의 SHA-256을 검사합니다.
- 결과물은 임시 파일로 생성하고 전체 해시 검증 후에만 확정합니다.
- 기존 결과가 정상이라면 재사용하고, 다르면 **덮어쓰지 않고 중단**합니다.

## v1.01 실행 안정성 개선

- **실행 시간 제한:** xdelta3 실행 후 최대 5분. 제한을 넘으면 외부 프로세스를 종료하고 임시 파일을 정리합니다.
- **무응답 자동 감지:** 임시 BIN 파일 크기가 증가하지 않고 xdelta3의 stdout/stderr 출력도 90초간 없는 경우 프로세스를 중단합니다. 경과 시간은 매 5초 간격으로 로그에 기록됩니다.
- **진행 정보:** 화면에 경과 시간, 현재 임시 BIN 크기(MiB), 마지막 활동 이후 경과 시간을 표시합니다.
- **오류 정보:** xdelta3의 표준 오류와 표준 출력의 메시지를 로그에 표시하고, 실패·시간 초과·무응답·사용자 취소 사유를 구분합니다.
- **사용자 중단:** 작업 중에는 **작업 중단** 버튼으로 외부 디코더를 종료할 수 있습니다.
- **임시 파일 정리:** 정상 완료, 오류, 시간 초과 및 사용자 중단 이후 프로그램이 생성한 `.bin.tmp` 및 `.cue.tmp`를 정리합니다. 원본 게임 이미지는 손대지 않습니다.
- **참고:** 느린 외장 드라이브나 보안 소프트웨어의 검사로 활동 감지가 제한될 수 있습니다. 이 경우 짧은 경고 이후 자동 중단될 수 있으므로 원본 폴더의 저장 공간과 접근 권한을 확인하세요.

## 검증 범위

- 제공 원본과 xdelta에서 VCDIFF/LZMA 해독을 수행하고 **3개 윈도우 모두 Adler-32 일치**를 확인했습니다.
- 결과 이미지의 단일 MODE1/2352 트랙과 `CD001` ISO9660 PVD를 확인했습니다.
- Go 단위 테스트 및 Windows x64/x86 크로스 빌드에 성공했습니다.
- Linux에서 Python VCDIFF/LZMA 독립 디코더로 수정 xdelta를 실제 복원했으며 새 결과 BIN의 MD5 및 SHA-256을 확인했습니다. CLI 기반 종단 간 패치 처리, 재실행, 기존 결과 보호, 중복 원본 감지, CUE 변조 탐지도 검증했습니다. 외부 디코더 오류/무응답/시간 초과/취소 시나리오를 Go 테스트로 검증했습니다.
- **Windows GUI 직접 실행, Windows xdelta3 실행파일의 실제 동작, DOS에서 게임을 실행해 한글 표시를 확인하는 테스트는 수행하지 않았습니다.**

## 글꼴 라이선스 / 출처

동봉된 `FONT_ATTRIBUTIONS.md`에서 다음 글꼴의 출처와 라이선스를 확인할 수 있습니다.

- **달무리 글꼴** — RanolP 및 기여자. **Apache License 2.0**. [`LICENSES/DALMOORI_APACHE-2.0.txt`](LICENSES/DALMOORI_APACHE-2.0.txt)
- **GNU Unifont** — Unifont 기여자. **SIL OFL 1.1 또는 GPL v2 이상(폰트 임베딩 예외 포함)**. [`LICENSES/UNIFONT_OFL-1.1.txt`](LICENSES/UNIFONT_OFL-1.1.txt), [`LICENSES/UNIFONT_GPL-2.0_WITH_FONT_EXCEPTION.txt`](LICENSES/UNIFONT_GPL-2.0_WITH_FONT_EXCEPTION.txt)

패처 자체의 Apache-2.0 고지와 별도이며, 글꼴의 저작권/재배포 조건은 해당 글꼴에 적용됩니다. 게임 CD 원본, 완성 한글판 BIN, 글꼴 원본 파일은 동봉하지 않습니다.

## 배포 및 라이선스

원본 CD 이미지와 패치 완료 게임 BIN은 ZIP에 포함하지 않습니다. 패치에 필요한 xdelta 파일은 배포 ZIP에 포함합니다.
패처 소스는 Apache License 2.0으로 배포하며, 동봉 xdelta3는 [jmacd/xdelta v3.2.0](https://github.com/jmacd/xdelta/tree/v3.2.0) 계열의 Apache-2.0 프로그램입니다.
게임과 번역 패치의 권리는 각각의 권리자에게 있습니다.
