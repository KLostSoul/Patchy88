# 몽환전사 바리스 PC-8801 Patchy88 검증 기록

## Disk A Ver.1.02 IPS

사용 패치:

```text
Valis_Korean_Disk_A_Patch_Ver_1.02.ips
SHA-256: 7f14c7b5d6961e234f702aa3e6007944ad3ec8231af225f6296ef3491e1eff53
```

IPS 구조 확인값:

- IPS 레코드: `1,282`
- RLE 레코드: `0`
- 실제 패치 대상 바이트: `23,330 bytes`
- 최소 오프셋: `12,651 (0x316B)`
- 최대 오프셋: `121,684 (0x1DB54, inclusive)`

## 전체 SHA가 다른 두 원본의 호환성

확인한 두 Disk A 원본의 전체 SHA-256은 서로 다릅니다.

```text
원본 1
SHA-256: ae9e0d57219763cc575e66d38e92c78e7f3fc7a6acdeba0e5f13d7f7dd920a44

원본 2
SHA-256: 7404998ee7e94e14d065a11e55bc26f7f8733202eec6774610a20a6d0b5a1fdf
```

하지만 Ver.1.02 IPS가 실제로 읽고 덮어쓰는 모든 대상 범위를 두 원본에서 바이트 단위로 비교한 결과:

```text
target_byte_diffs = 0
```

즉 패치 대상 `23,330 bytes`의 preimage는 완전히 동일합니다.

이 때문에 Patchy88은 전체 D88 SHA 하나를 강제하지 않고 **IPS가 실제로 건드리는 원본 영역을 검증**합니다.

## 적용 결과

원본 1에 적용:

```text
SHA-256: 08b389a69858cc3244799567cf985f9c473f2344623c14d2f249fc9900b4a93a
```

원본 2에 적용:

```text
SHA-256: 18e274dc730902f90e4d3939ad3ac2853c927d19baf896cee88e5b22321427b8
```

두 결과의 전체 SHA가 다른 이유는 IPS 대상 밖의 원본 바이트 차이가 그대로 남기 때문입니다. 실제 IPS 대상 패치 결과 영역은 동일합니다.

## KANJI1 ROM

확인한 재현 기준 원본:

```text
SHA-256: f6c1c5022fe5935f6dfa3eb919e51441e75191270b639edcb7938b3bce41f6a3
```

`VALIS_KANJI1_ROM_Patch_Ver_1.02.ips` 적용 결과:

```text
SHA-256: 3a4ce60dc4a23d7918a8726b99c2192c9420313bab40c50880eea3a387243f45
```

## Patchy88 판정 규칙

패치 대상 각 영역은 매니페스트에 기록된 `before`/`after` SHA-256으로 판정합니다.

- 모든 영역이 `before` → `ORIGINAL`
- 모든 영역이 `after` → `ALREADY_PATCHED`
- `before`와 `after` 혼재 → `PARTIAL`
- 어느 쪽도 아닌 영역 존재 → `INCOMPATIBLE`

`PARTIAL`, `INCOMPATIBLE` 상태에서는 원본 파일을 변경하지 않습니다.
