# Patchy88

Patchy88은 PC-8801용 IPS 패치를 안전하게 적용하기 위한 패처입니다.

46OkuMen의 **Pachy98 / romtools**가 사용한 패치 배포 개념에서 출발했지만, PC-88의 D88/ROM을 대상으로 하기 위해 구조를 바꿨습니다. NDC를 이용한 파일 추출/재삽입이나 xdelta3를 사용하지 않고, **D88/ROM에 IPS를 직접 적용하면서 IPS가 실제로 건드리는 원본 영역을 검증**합니다.

> **중요:** 원본 게임 D88, 원본 KANJI ROM, 패치 완료 전체 게임 이미지는 저장소에 포함하지 않습니다.

## 게임별 문서

게임별 사용법, 지원 범위, 버전, 검증 결과와 해시는 README에서 분리했습니다.

- [몽환전사 바리스용 Patchy88](docs/VALIS1.md)
- [몽환전사 바리스 II용 Patchy88](docs/VALIS2.md)

## 핵심 원리

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

## 안전 적용 절차

기본 처리 순서는 다음과 같습니다.

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

## 파일명은 검사 조건이 아님

Patchy88은 원본 파일명을 호환성 검사 조건으로 사용하지 않습니다.

예를 들어 다음처럼 파일명이 달라도 패치 대상 원본 데이터가 지원 기준과 일치하면 정상 식별할 수 있습니다.

```text
one.d88
two.d88
font.rom
```

## 출력과 백업

패치 결과는 기본적으로 확장자 앞에 `(K)`를 붙여 생성합니다.

```text
game.d88
→ game.d88.bak
→ game(K).d88
```

백업 폴더를 따로 만들지 않고 원본과 같은 폴더에 `.bak`을 둡니다.

기존 백업이 있으면 덮어쓰지 않습니다.

```text
game.d88.bak
game.d88.1.bak
game.d88.2.bak
```

## D88 처리 범위

현재 Patchy88은 IPS의 **D88 파일 오프셋을 직접 사용**합니다.

- NDC 미사용
- 파일시스템 추출 미사용
- xdelta3 미사용
- 논리 섹터 정규화 미사용
- 현재 지원 원본과 동일/호환되는 D88 직렬화에서 패치 대상 영역이 일치해야 함

논리적으로 같은 섹터 내용이라도 D88 내부 오프셋 배치가 달라지면 안전하게 `INCOMPATIBLE`로 거부할 수 있습니다.

## 저장소 구조

```text
Patchy88/
├─ README.md
├─ LICENSE
├─ NOTICE.md
├─ docs/
│  ├─ VALIS1.md
│  ├─ VALIS2.md
│  ├─ VALIDATION.md
│  ├─ Valis1_PC88_Validation.md
│  ├─ Valis2_PC88_KOR_Hash_List.md
│  └─ RELEASE_ARTIFACTS.md
├─ patches/
│  └─ README.md
└─ src/
   ├─ valis1-python/
   └─ valis2-go/
```

## 관련 문서

- [몽환전사 바리스](docs/VALIS1.md)
- [몽환전사 바리스 II](docs/VALIS2.md)
- [Patchy88 검증 구조](docs/VALIDATION.md)
- [IPS 파일 식별값](patches/README.md)
- [배포 산출물 기록](docs/RELEASE_ARTIFACTS.md)

## 라이선스 / 파생판 고지

Patchy88은 Pachy98/romtools의 패치 배포 개념에서 출발한 PC-8801용 파생 프로그램입니다.

- Original project: 46OkuMen / romtools / Pachy98
- Original repository: https://github.com/46OkuMen/romtools
- Original license: Apache License 2.0

NDC/xdelta 기반 처리 대신 PC-88 D88/ROM 직접 IPS 적용, 패치 대상 영역 검증, 내용 기반 자동 식별, 트랜잭션식 적용 등을 추가했습니다.

자세한 변경 고지는 `NOTICE.md`를 참조하십시오.

## 저작권 관련

이 저장소와 배포물은 **원본 게임 디스크 이미지, 원본 KANJI ROM, 패치 완료 전체 게임 이미지**를 제공하지 않습니다.

사용자는 적법하게 준비한 원본을 사용해야 합니다.
