# Patchy88

Patchy88은 PC-8801 게임의 한글화 패치를 안전하게 적용하기 위한 도구 모음입니다. 《몽환전사 바리스》·《몽환전사 바리스 II》용 IPS 패처와 《Mirrors》용 xdelta 패처를 제공하며, 같은 검증 원칙을 바탕으로 만든 FM TOWNS용 변형 **PatchyFM**도 이 저장소에서 관리합니다.

**PatchyDOS**는 MS-DOS판 《Dune II: The Battle for Arrakis》 CD 한글패치 전용 도구입니다. 폴더에서 원본 IMG/BIN을 자동 인식하고 xdelta를 적용하여 검증된 BIN/CUE를 생성합니다. v1.01 최신 수정본은 디코더 무응답 감지, 5분 실행 제한, 작업 중단, WaitDelay 오류 보정 및 실행 시각별 자동 진단 로그 저장을 지원합니다. 로그는 실행파일 옆 `logs/` 폴더에 생성되며 쓰기 불가 시 `%APPDATA%`의 사용자 설정 폴더를 이용합니다.

FM TOWNS 패처는 《Alltynex》용 PatchyFM과 《誘惑 / Yuuwaku》용 **PatchyFM_Yuuwaku_v1.01**을 제공합니다. Alltynex는 원본 ZIP에서 한글판 ISO를 생성하고, Yuuwaku v1.01은 지정 폴더와 하위 폴더의 원본 IMG/BIN을 해시로 자동 인식하고 수정된 xdelta를 적용하여 한글판 BIN/CUE를 생성합니다. 두 패처 모두 결과 전체의 크기·MD5·SHA-256을 검증합니다.

《Mirrors》 v1.01은 일본판과 영문판의 원본이 서로 다르지만, 두 판본 모두 **동일한 최종 한글판 CCD·IMG·SUB**를 생성합니다. 원본과 결과는 모두 전체 MD5·SHA-256과 크기로 검증합니다.

46OkuMen의 **Pachy98 / romtools**가 사용한 패치 배포 개념에서 출발했지만, PC-88의 D88/ROM과 FM TOWNS의 ZIP→ISO 변환에 맞게 구조를 바꿨습니다. PC-88 IPS 판에서는 NDC나 xdelta3를 사용하지 않고 D88/ROM에 IPS를 직접 적용하면서 IPS가 실제로 건드리는 원본 영역을 검증합니다. Mirrors·PatchyFM·PatchyDOS는 공식 xdelta3 v3.2.0을 사용합니다.

> **중요:** 원본 게임 D88, 원본 KANJI ROM, 원본 게임 ZIP, 패치 완료 전체 게임 이미지/ISO는 저장소와 배포물에 포함하지 않습니다.

## 게임별 문서

게임별 사용법, 지원 범위, 버전, 검증 결과와 해시는 README에서 분리했습니다.

- [몽환전사 바리스용 Patchy88](docs/VALIS1.md)
- [몽환전사 바리스 II용 Patchy88](docs/VALIS2.md)
- [Mirrors용 Patchy88 — Mirrors_Kor1.01](docs/MIRRORS.md)
- [PatchyFM — Alltynex Kor v1.0](docs/PATCHYFM_ALLTYNEX.md)
- [PatchyFM_Yuuwaku_v1.01 — 誘惑 한글패처](docs/PATCHYFM_YUUWAKU.md)
- [PatchyDOS_Dune2_v1.01 — Dune II 한글패처](docs/PATCHYDOS_DUNE2.md)
- [Mirrors용 Patchy88 소스 및 빌드](src/mirrors-go/README.ko.md)
- [PatchyFM Alltynex 소스 및 빌드](src/patchyfm-alltynex/README.ko.md)
- [PatchyFM Yuuwaku 소스 및 빌드](src/patchyfm-yuuwaku/README_SOURCE.md)
- [PatchyDOS Dune II 소스 및 빌드](src/patchydos-dune2/README_SOURCE.md)

## PC-88 IPS 판의 핵심 원리

Patchy88은 전체 파일 SHA-256 하나만으로 호환성을 판정하지 않습니다.

D88은 디스크/섹터 컨테이너이며 패치와 무관한 저장 데이터나 이미지 직렬화 차이 때문에 **전체 D88 SHA-256이 달라도 실제 패치 대상 바이트는 같을 수 있습니다.**

Patchy88은 IPS가 실제로 건드리는 원본 영역의 적용 전/후 SHA-256을 별도로 기록하여 상태를 판정합니다.

```text
D88 / ROM
  │
  ├─ IPS가 건드리지 않는 영역 → 호환성 판정에서 제외
  │
  └─ IPS가 건드리는 영역
       ├─ 원본 상태 해시 일치 → ORIGINAL
       ├─ 패치 결과 해시 일치 → ALREADY_PATCHED
       ├─ 원본/패치 상태 혼재 → PARTIAL
       └─ 어느 쪽도 아님       → INCOMPATIBLE
```

표준 IPS에는 출력 오프셋과 새 데이터/RLE 데이터만 들어 있고 **원본(preimage) 바이트는 포함되지 않습니다.** 따라서 Patchy88은 지원 원본에서 패치 대상 영역의 `before` / `after` SHA-256을 별도로 기록합니다.

자세한 내용은 [검증 구조 문서](docs/VALIDATION.md)를 참조하십시오.

## PC-88 IPS 판의 안전 적용 절차

```text
입력 선택
→ IPS 및 검증 정보 확인
→ 패치 대상 영역 사전검증
→ 상태 판정
→ 원본 백업
→ 임시 결과 생성
→ IPS 적용
→ 패치 대상 영역 사후검증
→ 검증 성공 시에만 (K) 결과 확정
```

상태 판정:

- `ORIGINAL` — 정상 원본, 패치 가능
- `ALREADY_PATCHED` — 이미 정상 패치됨, 재적용 안 함
- `PARTIAL` — 부분 패치 상태, 적용 거부
- `INCOMPATIBLE` — 지원하지 않는 데이터, 적용 거부

`PARTIAL`, `INCOMPATIBLE` 상태에서는 원본 파일을 변경하지 않습니다.

## PC-88 IPS 판의 파일명 검사

Patchy88은 원본 파일명을 호환성 검사 조건으로 사용하지 않습니다.

예를 들어 다음처럼 파일명이 달라도 패치 대상 원본 데이터가 지원 기준과 일치하면 정상 식별할 수 있습니다.

```text
one.d88
two.d88
font.rom
```

## PC-88 IPS 판의 출력과 백업

패치 결과는 기본적으로 확장자 앞에 `(K)`를 붙여 생성합니다.

```text
game.d88
→ game.d88.bak
→ game(K).d88
```

백업 폴더를 따로 만들지 않고 원본과 같은 폴더에 `.bak`을 둡니다. 기존 백업이 있으면 덮어쓰지 않습니다.

```text
game.d88.bak
game.d88.1.bak
game.d88.2.bak
```

## D88 처리 범위

현재 PC-88 IPS 판은 IPS의 **D88 파일 오프셋을 직접 사용**합니다.

- NDC 미사용
- 파일시스템 추출 미사용
- xdelta3 미사용
- 논리 섹터 정규화 미사용
- 현재 지원 원본과 동일/호환되는 D88 직렬화에서 패치 대상 영역이 일치해야 함

논리적으로 같은 섹터 내용이라도 D88 내부 오프셋 배치가 달라지면 안전하게 `INCOMPATIBLE`로 거부할 수 있습니다.

## xdelta 판의 검증

### Mirrors

- 일본판/영문판 원본 CCD·IMG·SUB를 판별
- 필요한 xdelta를 적용
- 두 판본 모두 동일한 최종 한글판 CCD·IMG·SUB를 생성
- 결과 전체의 크기·MD5·SHA-256 검증
- 검증 실패 시 결과를 확정하지 않음

### PatchyFM Alltynex

```text
사용자 원본 ZIP 선택
→ 원본 ZIP 전체 크기/MD5/SHA-256 검증
→ 동봉 xdelta 및 xdelta3 SHA-256 검증
→ 임시 ISO 생성
→ ISO 전체 크기/MD5/SHA-256 검증
→ 검증 성공 시에만 Alltynex (Kor v1.0).iso 확정
```

원본 ZIP은 수정하지 않습니다. 같은 이름의 기존 ISO가 올바른 결과 해시라면 그대로 인정하고, 해시가 다르면 덮어쓰지 않습니다.

PatchyFM Alltynex v1.0은 실제 원본 ZIP에서 ISO 생성, 결과 검증, 게임 실행까지 확인했습니다.

### PatchyFM Yuuwaku v1.01

```text
원본 IMG/BIN이 있는 폴더 선택
→ 하위 폴더 포함 원본 자동 검색
→ 원본 크기/MD5/SHA-256 검증
→ xdelta 및 xdelta3 검증
→ 임시 BIN 생성/결과 검증
→ 한글판 BIN/CUE 확정
```

CUE·CCD·SUB의 수동 지정은 필요하지 않으며 원본은 수정하지 않습니다. 원본과 동일한 이미지가 여러 개 검색되면 잘못된 폴더에 결과가 생성되는 것을 방지하기 위해 중단합니다. 자세한 내용은 [Yuuwaku 사용법](docs/PATCHYFM_YUUWAKU.md)을 참조하십시오.

### PatchyDOS Dune II

```text
CD 이미지가 있는 폴더 선택
→ 하위 폴더에서 지원 원본 IMG/BIN 자동 검색
→ 원본 크기·MD5·SHA-256 검사
→ 동봉 xdelta 및 xdelta3 무결성 검사
→ 5분 전체 제한 / 90초 무활동 감시하며 임시 BIN 생성
→ 최종 BIN/CUE 전체 해시 검사
→ 검증 후 결과 확정 (실패/취소 시 임시 파일 정리)
```

기존 원본 게임 이미지는 수정하지 않습니다. 최신 v1.01 수정 배포 ZIP은 `PatchyDOS_Dune2_v1.01_WaitDelayFix_Log.zip`이며, 종료 지연 오류 수정과 .log 파일 자동 저장을 포함합니다. [Dune II v1.01 사용법 및 검증값](docs/PATCHYDOS_DUNE2.md)을 참조하십시오.

## 저장소 구조

```text
Patchy88/
├─ README.md
├─ LICENSE
├─ NOTICE.md
├─ docs/
│  ├─ VALIS1.md
│  ├─ VALIS2.md
│  ├─ MIRRORS.md
│  ├─ PATCHYFM_ALLTYNEX.md
│  ├─ PATCHYFM_YUUWAKU.md
│  ├─ PATCHYDOS_DUNE2.md
│  ├─ VALIDATION.md
│  ├─ Valis1_PC88_Validation.md
│  ├─ Valis2_PC88_KOR_Hash_List.md
│  └─ RELEASE_ARTIFACTS.md
├─ patches/
│  └─ README.md
└─ src/
   ├─ valis1-python/
   ├─ valis2-go/
   ├─ mirrors-go/
   ├─ patchyfm-alltynex/
   ├─ patchyfm-yuuwaku/
   └─ patchydos-dune2/
```

## 관련 문서

- [몽환전사 바리스](docs/VALIS1.md)
- [몽환전사 바리스 II](docs/VALIS2.md)
- [Mirrors용 Patchy88 사용법](docs/MIRRORS.md)
- [PatchyFM Alltynex 사용법](docs/PATCHYFM_ALLTYNEX.md)
- [PatchyFM Yuuwaku 사용법](docs/PATCHYFM_YUUWAKU.md)
- [PatchyDOS Dune II 사용법](docs/PATCHYDOS_DUNE2.md)
- [Valis IPS 검증 구조](docs/VALIDATION.md)
- [IPS 파일 식별값](patches/README.md)
- [배포 산출물 기록](docs/RELEASE_ARTIFACTS.md)

## 라이선스 / 파생판 고지

Patchy88은 Pachy98/romtools의 패치 배포 개념에서 출발한 파생 프로그램입니다.

- Original project: 46OkuMen / romtools / Pachy98
- Original repository: https://github.com/46OkuMen/romtools
- Original license: Apache License 2.0

Valis 판에서는 PC-8801 D88/ROM 직접 IPS 적용과 패치 영역 검증을 사용합니다. Mirrors와 PatchyFM에서는 xdelta3와 원본/결과 전체 해시 검증을 사용합니다.

PatchyFM Alltynex가 생성하는 ISO의 부팅 환경에는 CaptainYS(Soji Yamakawa)의 **FreeTOWNSOS / TSUGARU OS**가 사용됩니다. FreeTOWNSOS의 원 라이선스는 PatchyFM 배포물과 소스 문서에 함께 포함합니다.

자세한 변경 고지는 [NOTICE.md](NOTICE.md)를 참조하십시오.

## 저작권 관련

이 저장소와 배포물은 **원본 게임 디스크 이미지, 원본 KANJI ROM, 원본 게임 ZIP, 패치 완료 전체 게임 이미지/ISO**를 제공하지 않습니다.

사용자는 적법하게 준비한 원본을 사용해야 합니다.
