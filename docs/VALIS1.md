# 몽환전사 바리스용 Patchy88

PC-8801 《몽환전사 바리스》 한글패치용 Patchy88 문서입니다.

## 현재 버전

- Windows x64/x86 독립 실행판: **Patchy88 Valis1 v1.0.7**
- Python판: **Patchy88 Valis1 Python v1.0.7**
- 대상
  - Disk A D88
  - KANJI1 ROM

## 사용 방법

Windows판은 자신의 환경에 맞는 실행파일을 실행합니다.

```text
Patchy88-Valis1-x64.exe
Patchy88-Valis1-x86.exe
```

Python판은 다음 파일을 사용할 수 있습니다.

```text
Patchy88-Valis1.bat
Patchy88-Valis1.pyw
Patchy88_Valis1.py
```

GUI에서 Disk A와 KANJI1 ROM을 각각 선택한 뒤 패치를 적용합니다.

Patchy88은 **파일명을 검사 기준으로 사용하지 않습니다.** 파일 이름이 달라도 패치 대상 원본 데이터가 지원 기준과 일치하면 정상 원본으로 판정합니다.

## 출력 규칙

패치 결과는 원본과 같은 폴더에 `(K)`를 붙여 생성합니다.

```text
VALIS.D88
→ VALIS.D88.bak
→ VALIS(K).D88

KANJI1.ROM
→ KANJI1.ROM.bak
→ KANJI1(K).ROM
```

기존 백업이 있으면 덮어쓰지 않습니다.

```text
VALIS.D88.bak
VALIS.D88.1.bak
VALIS.D88.2.bak
```

## 검증 방식

Patchy88은 전체 D88 SHA-256 하나만으로 호환성을 결정하지 않습니다.

표준 IPS에는 원본 바이트가 없으므로, 지원 원본에서 IPS가 실제로 건드리는 영역의 적용 전/후 SHA-256을 별도로 기록해 판정합니다.

- 모든 대상 영역이 적용 전 값과 일치 → `ORIGINAL`
- 모든 대상 영역이 적용 후 값과 일치 → `ALREADY_PATCHED`
- 적용 전/후 상태가 혼재 → `PARTIAL`
- 어느 쪽에도 맞지 않는 대상 영역 존재 → `INCOMPATIBLE`

`PARTIAL`, `INCOMPATIBLE` 상태에서는 원본을 변경하지 않습니다.

## Disk A Ver.1.02 검증 결과

사용 IPS:

```text
Valis_Korean_Disk_A_Patch_Ver_1.02.ips
SHA-256: 7f14c7b5d6961e234f702aa3e6007944ad3ec8231af225f6296ef3491e1eff53
```

IPS 구조:

- 레코드: `1,282`
- RLE 레코드: `0`
- 실제 대상 바이트: `23,330 bytes`
- 최소 오프셋: `0x316B`
- 최대 오프셋: `0x1DB54` inclusive

### 전체 SHA가 다른 호환 원본

확인한 두 Disk A 원본은 전체 SHA-256이 서로 다릅니다.

```text
ae9e0d57219763cc575e66d38e92c78e7f3fc7a6acdeba0e5f13d7f7dd920a44
7404998ee7e94e14d065a11e55bc26f7f8733202eec6774610a20a6d0b5a1fdf
```

그러나 Ver.1.02 IPS가 건드리는 `23,330 bytes`를 두 원본에서 비교한 결과 차이는 없었습니다.

```text
target_byte_diffs = 0
```

따라서 두 원본 모두 패치 대상 영역 기준으로는 같은 원본이며 Patchy88이 안전하게 허용할 수 있습니다.

적용 결과 전체 SHA-256은 원본의 패치 대상 밖 데이터 차이 때문에 서로 달라집니다.

```text
원본 1 적용 결과
08b389a69858cc3244799567cf985f9c473f2344623c14d2f249fc9900b4a93a

원본 2 적용 결과
18e274dc730902f90e4d3939ad3ac2853c927d19baf896cee88e5b22321427b8
```

## KANJI1 ROM 검증 결과

확인한 원본:

```text
SHA-256: f6c1c5022fe5935f6dfa3eb919e51441e75191270b639edcb7938b3bce41f6a3
```

`VALIS_KANJI1_ROM_Patch_Ver_1.02.ips` 적용 결과:

```text
SHA-256: 3a4ce60dc4a23d7918a8726b99c2192c9420313bab40c50880eea3a387243f45
```

## 관련 자료

- [바리스 1 상세 검증 기록](Valis1_PC88_Validation.md)
- [Patchy88 공통 검증 구조](VALIDATION.md)
- [IPS 식별값](../patches/README.md)

## 원본 데이터 배포 정책

Patchy88 저장소와 배포물에는 원본 게임 D88, 원본 KANJI ROM, 패치 완료 전체 게임 이미지를 포함하지 않습니다. 사용자는 적법하게 준비한 원본을 사용해야 합니다.
