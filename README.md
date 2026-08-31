# Patchy88

Patchy88은 PC-8801용 IPS 패치를 안전하게 적용하기 위한 패처입니다.

46OkuMen의 **Pachy98 / romtools**가 사용한 패치 배포 개념에서 출발했지만, PC-88의 D88/ROM을 대상으로 하기 위해 구조를 바꿨습니다. NDC를 이용한 파일 추출/재삽입이나 xdelta3를 사용하지 않고, **D88/ROM에 IPS를 직접 적용하면서 IPS가 실제로 건드리는 원본 영역을 검증**합니다.

현재 저장소에는 《몽환전사 바리스》와 《몽환전사 바리스 II》용 Patchy88의 소스와 검증 자료를 정리합니다.

> **중요:** 원본 게임 D88, 원본 KANJI ROM, 패치 완료 전체 게임 이미지는 저장소에 포함하지 않습니다.

## 현재 버전

### 몽환전사 바리스

- Windows x64/x86 독립 실행판: **Patchy88 Valis1 v1.0.7**
- Python판: **Patchy88 Valis1 Python v1.0.7**
- 대상: Disk A, KANJI1 ROM
- 사용자가 D88/ROM을 직접 선택
- 결과는 원본과 같은 폴더에 `(K)` 접미사로 생성

```text
VALIS.D88
→ VALIS.D88.bak
→ VALIS(K).D88

KANJI1.ROM
→ KANJI1.ROM.bak
→ KANJI1(K).ROM
```

### 몽환전사 바리스 II

- Windows x64/x86 독립 실행판: **Patchy88 Valis2 v1.0.0**
- 대상: Disk A~G + KANJI1 ROM
- 사용자는 원본 8개가 들어 있는 폴더 하나만 지정
- 파일명 대신 패치 대상 원본 데이터로 A~G/KANJI1 자동 식별
- 8개 전체를 사전검증한 뒤에만 패치 시작
- 8개 임시 결과를 모두 사후검증한 뒤 실제 결과 확정

## 핵심 원리

Patchy88은 **전체 파일 SHA-256 하나만으로 원본을 판정하지 않습니다.**

D88은 디스크/섹터 컨테이너이므로 패치와 무관한 저장 데이터나 직렬화 차이 때문에 전체 SHA가 달라도 실제 패치 대상 바이트가 같을 수 있습니다. 따라서 IPS가 실제로 건드리는 영역의 `before` / `after` 해시를 사용합니다.

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

표준 IPS에는 출력 오프셋과 새 데이터/RLE 데이터만 있으며 **원본(preimage) 바이트는 들어 있지 않습니다.** Patchy88은 지원 원본에서 패치 대상 영역의 검증값을 별도로 준비하여 이 문제를 해결합니다.

자세한 내용은 [`docs/VALIDATION.md`](docs/VALIDATION.md)를 참조하십시오.

## 바리스 1에서 확인한 실제 사례

서로 전체 SHA-256이 다른 두 Disk A 원본을 확인했습니다.

```text
ae9e0d57219763cc575e66d38e92c78e7f3fc7a6acdeba0e5f13d7f7dd920a44
7404998ee7e94e14d065a11e55bc26f7f8733202eec6774610a20a6d0b5a1fdf
```

Disk A Ver.1.02 IPS가 실제로 건드리는 `23,330 bytes`를 두 원본에서 전부 비교한 결과 차이는 `0 bytes`였습니다. 그래서 Patchy88은 두 이미지를 모두 정상 원본으로 안전하게 인정할 수 있습니다.

세부 검증 기록: [`docs/Valis1_PC88_Validation.md`](docs/Valis1_PC88_Validation.md)

## 안전 적용 절차

```text
입력 선택
→ IPS/검증정보 검사
→ 패치 대상 영역 사전검증
→ 상태 판정
→ 원본 백업
→ 임시 결과 생성
→ IPS 적용
→ 패치 대상 영역 사후검증
→ 검증 성공 시에만 (K) 결과 확정
```

- `ORIGINAL` — 정상 원본, 패치 가능
- `ALREADY_PATCHED` — 이미 정상 패치됨, 재적용 안 함
- `PARTIAL` — 부분 패치 상태, 적용 거부
- `INCOMPATIBLE` — 지원하지 않는 데이터, 적용 거부

### 백업 규칙

백업 폴더를 따로 만들지 않습니다.

```text
game.d88       → game.d88.bak
기존 .bak 존재 → game.d88.1.bak
또 존재        → game.d88.2.bak
```

기존 백업은 덮어쓰지 않습니다.

## 바리스 II 일괄 패치

바리스 II판은 한 파일씩 즉시 수정하지 않습니다.

1. 지정 폴더 바로 아래의 `.d88`, `.rom` 검사
2. 내용으로 Disk A~G / KANJI1 자동 식별
3. 8개 모두 유일하게 식별되는지 확인
4. 누락/중복/부분패치/비호환/IPS 변조가 하나라도 있으면 전체 중단
5. 8개 임시 패치 결과 생성
6. 8개 결과 모두 사후검증
7. 검증 완료 후 일괄 확정

파일명을 `one.d88`, `two.d88`, `font.rom`처럼 바꿔도 패치 대상 원본 데이터가 맞으면 내용으로 식별합니다.

## D88 처리 범위

현재 Patchy88은 IPS의 **D88 파일 오프셋을 직접 사용**합니다.

- NDC 미사용
- 파일시스템 추출 미사용
- 논리 섹터 정규화 미사용
- 지원 원본과 동일/호환되는 D88 직렬화에서 패치 대상 영역이 일치해야 함

논리적으로 같은 섹터 내용이라도 D88 내부 오프셋 배치가 달라지면 `INCOMPATIBLE`로 거부할 수 있습니다.

## 검증 자료

- [`docs/VALIDATION.md`](docs/VALIDATION.md) — Patchy88 검증 설계
- [`docs/Valis1_PC88_Validation.md`](docs/Valis1_PC88_Validation.md) — 바리스 1 실제 호환 원본 검증
- [`docs/Valis2_PC88_KOR_Hash_List.md`](docs/Valis2_PC88_KOR_Hash_List.md) — 바리스 II A~G/KANJI1 원본 및 패치 결과 CRC32/CRC64/SHA-256/SHA-1
- [`patches/README.md`](patches/README.md) — 현재 사용 IPS의 파일명·크기·SHA-256·Git blob SHA-1
- [`docs/RELEASE_ARTIFACTS.md`](docs/RELEASE_ARTIFACTS.md) — 현재 생성한 배포 ZIP의 크기와 SHA-256

CRC64는 **CRC-64/ECMA-182** 기준입니다.

## 바리스 II 패치 결과 SHA-256

| 대상 | 패치 결과 SHA-256 |
|---|---|
| Disk A | `B96D108D5BA75023E67AD7C66BE370E5E1AAC8A6A73E8F8D5BD248F65C0ABF9A` |
| Disk B | `2368B443A1A536AFAB78F28C582DA268196A2C28F7305B2BAEDF11CCE2A161B4` |
| Disk C | `C6652CF326B9C5250C3C91F11BE9936BE430FD1182D1736CB1132AFC7CF36570` |
| Disk D | `5964758193523428947B0681E2A2B2F7B7ADD379FD86F3A5933022507AF16331` |
| Disk E | `A872EB595CD8CD69377749EDC89786D5076939997DC7AC4A8927EC6E17962219` |
| Disk F | `E82286A4E9C3F8356DA7750BF9F80D504420022758666F7DA44E10A30C2750E4` |
| Disk G | `51679430371D46474C21C6C7CC0F37FA0C05288E0EBE3719A4E18A8A837987E7` |
| KANJI1 | `30435619932143104C6E7D38E282372DDD011F4B7444D49E566AB995192491DC` |

## 소스

### Valis II Windows판

Go 소스는 [`src/valis2-go/`](src/valis2-go/)에 있습니다.

```bash
go test ./...
GOOS=windows GOARCH=amd64 go build -ldflags="-H windowsgui" -o Patchy88-Valis2-x64.exe
GOOS=windows GOARCH=386   go build -ldflags="-H windowsgui" -o Patchy88-Valis2-x86.exe
```

### Valis I Python판

Python 소스는 [`src/valis1-python/`](src/valis1-python/)에 있습니다. 외부 pip 패키지는 필요하지 않으며 GUI는 표준 `tkinter`를 사용합니다.

```bash
python src/valis1-python/Patchy88_Valis1.py
```

## 배포 산출물

현재 생성한 배포판:

- `Patchy88_Valis1_PC88_v1.0.7.zip`
- `Patchy88_Valis1_PC88_Python_v1.0.7.zip`
- `Patchy88_Valis2_PC88_v1.0.0.zip`

정확한 크기와 SHA-256은 [`docs/RELEASE_ARTIFACTS.md`](docs/RELEASE_ARTIFACTS.md)에 기록합니다.

배포 ZIP은 소스 트리에 중복 저장하지 않고 GitHub Releases 자산으로 배포하는 것을 원칙으로 합니다.

## 라이선스 / 파생판 고지

Patchy88은 Pachy98/romtools의 패치 배포 개념에서 출발한 PC-8801용 파생 프로그램입니다.

- Original project: 46OkuMen / romtools / Pachy98
- Original repository: https://github.com/46OkuMen/romtools
- Original license: Apache License 2.0

NDC/xdelta 기반 처리 대신 PC-88 D88/ROM 직접 IPS 적용, 패치 대상 영역 검증, 자동 식별, 트랜잭션식 일괄 적용 등을 추가했습니다.

자세한 변경 고지는 [`NOTICE.md`](NOTICE.md)를 참조하십시오.

## 저작권 관련

이 저장소는 **원본 게임 디스크 이미지, 원본 KANJI ROM, 패치 완료 전체 게임 이미지**를 제공하지 않습니다.

사용자는 적법하게 준비한 원본을 사용해야 합니다.
