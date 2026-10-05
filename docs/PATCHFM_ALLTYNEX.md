# PatchFM — Alltynex Kor v1.0

PatchFM은 Patchy88의 검증 원칙을 FM TOWNS용으로 단순화한 xdelta 패처입니다.

이 버전은 FM TOWNS판 《Alltynex》의 원본 ZIP을 직접 입력으로 받아 한글판 ISO를 생성합니다.

## 배포 원칙

PatchFM 배포물에는 다음 파일을 포함하지 않습니다.

- 원본 `alltynex_fmtowns.zip`
- 완성 `Alltynex (Kor v1.0).iso`

사용자는 적법하게 준비한 원본 ZIP을 직접 선택해야 합니다.

배포물에는 다음 요소만 포함합니다.

```text
PatchFM_Alltynex_v1.0/
├─ PatchFM-Alltynex-x64.exe
├─ PatchFM-Alltynex-x86.exe
├─ README.ko.md
├─ LICENSE
├─ THIRD_PARTY_LICENSE.txt
├─ NOTICE_MODIFICATIONS.txt
├─ FREETOWNSOS_LICENSE.md
├─ UPSTREAM_XDELTA.txt
├─ SHA256SUMS.txt
└─ assets/
   ├─ Alltynex_Kor_v1.0.json
   ├─ xdelta3-x64.exe
   ├─ xdelta3.exe
   └─ patches/
      └─ Alltynex_Kor_v1.0.xdelta
```

## 지원 원본

PatchFM은 파일명만으로 원본을 인정하지 않습니다. 선택한 ZIP 전체를 크기·MD5·SHA-256으로 검사합니다.

| 항목 | 값 |
|---|---|
| 크기 | `574,208 bytes` |
| MD5 | `bca57cbe9a2f23c922118f40fffae43a` |
| SHA-256 | `19b662f038e50e99535f60dbe6e07a77e269f574ed9250d46fc0655bb2aa6b3a` |

검증값이 다르면 xdelta를 적용하지 않습니다.

## xdelta

동봉 패치:

`Alltynex_Kor_v1.0.xdelta`

SHA-256:

`d7ef05131f9f2f2809e73e7e77adf77f4ff19f89a2c59bcf6e604ec8f4cfc99b`

PatchFM 시작 시 xdelta 패치와 현재 실행 아키텍처에 맞는 xdelta3 실행파일의 SHA-256도 검증합니다.

## 패치 절차

```text
원본 ZIP 선택
→ 원본 ZIP 전체 크기/MD5/SHA-256 검증
→ xdelta 및 xdelta3 무결성 확인
→ 임시 ISO에 xdelta 적용
→ 임시 ISO 전체 크기/MD5/SHA-256 검증
→ 검증 성공 시에만 최종 ISO 이름으로 확정
```

원본 ZIP은 수정하지 않습니다.

## 결과 ISO

| 항목 | 값 |
|---|---|
| 파일명 | `Alltynex (Kor v1.0).iso` |
| 크기 | `5,222,400 bytes` |
| MD5 | `d5a438b7c7aac9e5b08baea6cb361f70` |
| SHA-256 | `114a9ffa82e32bcc8cacb8b0475810422056560260eae5697ccf32b17bc91a54` |

같은 폴더에 같은 이름의 ISO가 이미 있을 경우:

- 결과 해시와 정확히 일치하면 정상 패치 결과로 인정
- 해시가 다르면 덮어쓰지 않고 중단

## 실제 검증 상태

PatchFM Alltynex v1.0은 사용자 환경에서 다음 과정을 실제로 완료했습니다.

1. 원본 ZIP 선택
2. 원본 해시 검증
3. xdelta 적용
4. 한글판 ISO 생성
5. 결과 ISO 전체 해시 검증
6. 생성 ISO 실행
7. 한글 패치가 적용된 게임 실행 확인

따라서 단순 빌드 테스트가 아니라 실제 입력부터 실행까지의 종단 간 검증이 완료된 배포판입니다.

## v1.0 파일 선택 종료 문제 수정

초기 v1.0 Windows 빌드에서는 `ZIP 선택` 버튼을 누르면 프로그램이 종료되는 문제가 있었습니다.

Win32 `OPENFILENAMEW`의 필터 문자열은 NUL로 구분되어야 하는데, Go의 `syscall.StringToUTF16`에 내부 NUL이 포함된 문자열을 직접 넘기면서 panic이 발생한 것이 원인이었습니다.

UTF-16 필터 항목을 각각 생성하여 올바르게 결합하도록 수정했습니다.

버전 표기는 사용자의 요구에 따라 **v1.0을 그대로 유지**했습니다.

## FreeTOWNSOS / TSUGARU OS

한글판 ISO의 부팅 환경에는 CaptainYS(Soji Yamakawa)의 FreeTOWNSOS / TSUGARU OS를 사용합니다.

- Upstream: https://github.com/captainys/FreeTOWNSOS
- Copyright 2021 Soji Yamakawa
- 원본 라이선스: `src/patchfm-alltynex/FREETOWNSOS_LICENSE.md`

원 프로젝트의 `LICENSE.md` 전문을 수정하지 않고 배포물에 포함합니다. FreeTOWNSOS 라이선스에 기재된 ORICON과 Free386 관련 고지도 유지합니다.

## xdelta3

공식 jmacd/xdelta v3.2.0 계열을 사용합니다.

- Upstream: https://github.com/jmacd/xdelta/tree/v3.2.0
- License: Apache License 2.0
- 상세 고지: `src/patchfm-alltynex/UPSTREAM_XDELTA.md`

## 소스

소스 위치:

`src/patchfm-alltynex/`

주요 파일:

- `core.go` — 원본/결과 해시 검증 및 xdelta 적용
- `main_windows.go` — Windows GUI
- `main_cli.go` — 비 Windows CLI 빌드
- `core_test.go` — 코어 단위 테스트
- `Alltynex_Kor_v1.0.json` — 검증 매니페스트
