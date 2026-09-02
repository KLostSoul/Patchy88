# 몽환전사 바리스 II용 Patchy88

PC-8801 《몽환전사 바리스 II》 한글패치용 Patchy88 문서입니다.

## 현재 버전

- Windows x64/x86 독립 실행판: **Patchy88 Valis2 v1.0.3**
- 대상: Disk A~G D88 + KANJI1 ROM
- KANJI1 기준 원본: CRC32 `6178BD43`

## v1.0.3 변경 사항

사용자가 제공한 `V1.01` IPS 세트로 Disk A~G와 KANJI1 패치를 교체했습니다.

- A~F: 이전 Patchy88 패치와 실제 IPS 내용 변경
- G: 이전 패치와 동일
- KANJI1: 이전 패치와 실제 IPS 내용 변경
- 8개 IPS 기준으로 패치 대상 `before` / `after` SHA-256 매니페스트 전면 재생성
- 원본 8개 교차 자동식별 오인 없음 확인
- 8개 원본은 각각 `ORIGINAL`, 새 IPS 적용 결과는 각각 `ALREADY_PATCHED` 판정
- 배포물 내부에서도 패치 버전을 구분할 수 있도록 사용자가 지정한 `_V1.01` 파일명을 그대로 유지

배포물의 IPS 파일명은 다음과 같습니다.

```text
Valis2_KOR_Disk_A_V1.01.ips
Valis2_KOR_Disk_B_V1.01.ips
Valis2_KOR_Disk_C_V1.01.ips
Valis2_KOR_Disk_D_V1.01.ips
Valis2_KOR_Disk_E_V1.01.ips
Valis2_KOR_Disk_F_V1.01.ips
Valis2_KOR_Disk_G_V1.01.ips
Valis2_KOR_KANJI1_V1.01.ips
```

## 사용 방법

사용자는 **Disk A~G와 KANJI1 ROM이 들어 있는 폴더 하나만 지정**합니다.

Patchy88은 파일명이 아니라 IPS가 실제로 건드리는 **원본 데이터 영역**을 검사하여 Disk A~G와 KANJI1을 자동 식별합니다.

전체 D88/ROM SHA-256이 달라도 IPS 대상 원본 데이터가 모두 일치하면 호환 원본으로 인정할 수 있습니다. 현재 구현은 D88 파일 오프셋 기준이므로 대상 데이터의 오프셋 배치가 달라진 이미지는 호환되지 않을 수 있습니다.

## 일괄 패치 절차

1. 지정 폴더의 `.d88`, `.rom` 검사
2. Disk A~G / KANJI1 자동 식별
3. 8개 전체 사전검증
4. 누락·중복·부분패치·비호환·IPS 변조가 하나라도 있으면 전체 중단
5. 8개 임시 패치 결과 생성
6. 8개 결과 모두 사후검증
7. 전부 정상일 때만 실제 결과 확정

## 출력 규칙

```text
Mugen Senshi Valis II (Disk A).d88
→ Mugen Senshi Valis II (Disk A).d88.bak
→ Mugen Senshi Valis II (Disk A)(K).d88

KANJI1.ROM
→ KANJI1.ROM.bak
→ KANJI1(K).ROM
```

기존 백업은 덮어쓰지 않습니다.

## 기준 원본 및 V1.01 패치 결과 SHA-256

| 대상 | 원본 SHA-256 | 패치 결과 SHA-256 |
|---|---|---|
| Disk A | `4E1E3F5A21B66BB9C57AB20F6522C95C8074DC9120F92D395EACDB673D9A3BD8` | `976B59ADBEC8BCB3B5962092F8C903082CE6E456B19092245728C69497957879` |
| Disk B | `6537B681F8C24AA92BE608F99DC61EF8DA9AFA1B6569450DBA1DD5D7D00E3919` | `B586D114D3E35F2ADC57D11E434059D99972AF583756E51D86E7F9072BFCD580` |
| Disk C | `396F79DF59D3263B801126943D03F0C63B232A52A6256A0994254D83031FABDA` | `23EE625EEF0F4C8A1C9770FB28BEB323DECB600D2B2C45A80F751B762F3CF3DE` |
| Disk D | `BD453F4576705A224A6B1775D2E7D0D2D499EFD49075FF0B2B4B0B419549DD6D` | `F5E92F31FB9616516CE26A2B4319CE1E167BEA4BAB79BF1FE1FAE0AD4035B2E4` |
| Disk E | `905C9B0DBC402BB6134B63A67F3615AB0ADF79A191FF87CDF365903138CEFEDE` | `5A7ADEFC5849787E90000BFD9A0ECED27E6A4E52577E9AC67D0AAC69F341AE07` |
| Disk F | `8569EE720C0DF24293066FC3589C1A204EDC102EF57ABC3A6529C16EC8A4041B` | `8630A1CC1E2E032FF103F4CBC3FFEB690A56A5A4F6046C583B01555604051A1C` |
| Disk G | `D933A387F6AFE847345482070CD1E23CA4158F6D49144A99EF9EBBDB864C56F7` | `51679430371D46474C21C6C7CC0F37FA0C05288E0EBE3719A4E18A8A837987E7` |
| KANJI1 | `7608040CFFB1951E5CC567ABB63F75B5746777A1BA96196C1B75606B793BB4BB` | `90F5BD468DB8758950C60646C77538956F8B52008E1919C8BB73EE8A05D08242` |

## 관련 자료

- [원본/패치 결과 전체 해시](Valis2_PC88_KOR_Hash_List.md)
- [Patchy88 공통 검증 구조](VALIDATION.md)
- [IPS 식별값](../patches/README.md)

## 원본 데이터 배포 정책

Patchy88 저장소와 배포물에는 원본 게임 D88, 원본 KANJI ROM, 패치 완료 전체 게임 이미지를 포함하지 않습니다.
