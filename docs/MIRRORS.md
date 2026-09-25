# Mirrors용 Patchy88 — Mirrors_Kor1.01 (PC-8801)

일본판·영문판 PC-8801 《Mirrors》의 원본 CCD·IMG·SUB를 MD5·SHA-256으로 식별하고, **선택한 판본 전용 xdelta 3개**를 적용합니다. 두 판본의 패치 결과 해시는 동일하지 않습니다.

## 사용 방법

배포 ZIP 전체를 압축 해제해 Windows x64용 `Mirrors_Kor1.01-x64.exe` 또는 x86용 `Mirrors_Kor1.01-x86.exe`를 실행합니다. 옆의 `assets/` 폴더를 그대로 유지하십시오. 원본 CCD·IMG·SUB 폴더를 선택하고 검사 후 적용합니다. 일본판과 영문판이 동시에 있으면 **예=일본판, 아니오=영문판, 취소=중단**으로 선택합니다.

결과는 동일한 원본 폴더의 `Mirrors_Kor1.01.ccd`, `Mirrors_Kor1.01.img`, `Mirrors_Kor1.01.sub`, `Mirrors_Kor1.01.cue`, `disk1main.d88`, `disk2game.d88`에 생성됩니다. CUE는 `Mirrors_Kor1.01.img`를 참조합니다. 추가 D88 두 개는 명칭이 다르지만 SHA-256은 같습니다.

**일본판과 영문판을 모두 패치할 때는 별도 폴더를 사용하십시오.** 기존에 다른 판본의 결과가 있으면 출력 파일명 충돌을 안전하게 거부하며 원본을 덮어쓰지 않습니다.

## 일본판 원본 해시

| 파일 | 크기 | MD5 | SHA-256 |
|---|---:|---|---|
| CCD | 3,500 B | `80273154A2DAF2D107A282B353A86246` | `5512EFCEC20279753E1C6FEF0C5AAE37F2F3F05BB3AE3900B5DB0F688D4A67AC` |
| IMG | 551,779,200 B | `84E583EC62372CA9DD0D737DCE38B5D8` | `258533B4AC5FD8B16170DDF9509DCF2D6C4A0959BAC61E999B3D14C9C0B48D65` |
| SUB | 22,521,600 B | `E45923398B1D150F71E6056BD0974D15` | `C20C93FF3A2CF3A1E3A18B5F8046EC6F6198A916B417A01678B0B9D10EED9046` |

## 영문판 원본 해시

| 파일 | 크기 | MD5 | SHA-256 |
|---|---:|---|---|
| CCD | 3,500 B | `35C733769D60277FCCE522E121AF82AE` | `2DAEAAF64FD4C206CC28C2438A5FA480272DB0888CC19106380D13950DE1C102` |
| IMG | 551,779,200 B | `F3BC8D9E2650D9D9C6EBAEF5E2353144` | `294EE461A2BEA8745394C71CA7EEC64960A865673D72F8B60557C0A29A9CA7AA` |
| SUB | 22,521,600 B | `F6D739F1B66082F7F06F403D30CF35EB` | `3C05EFEC20E9BE2CDA81500178FF1F8A7657908FD699D2070827E58D700122D2` |

## 일본판 패치 결과

| 파일 | 크기 | MD5 | SHA-256 |
|---|---:|---|---|
| CCD | 3,500 B | `80273154A2DAF2D107A282B353A86246` | `5512EFCEC20279753E1C6FEF0C5AAE37F2F3F05BB3AE3900B5DB0F688D4A67AC` |
| IMG | 551,779,200 B | `56E768F7CE3315A8172338CB10CE153E` | `FDCF60364815ADF0E85C2B796533276E2C72210F1024425C02757BDF88333B10` |
| SUB | 22,521,600 B | `E45923398B1D150F71E6056BD0974D15` | `C20C93FF3A2CF3A1E3A18B5F8046EC6F6198A916B417A01678B0B9D10EED9046` |

## 영문판 패치 결과 (제공받은 CloneCD 측정값)

| 파일 | 크기 | MD5 | SHA-256 |
|---|---:|---|---|
| CCD | 3,500 B | `35C733769D60277FCCE522E121AF82AE` | `2DAEAAF64FD4C206CC28C2438A5FA480272DB0888CC19106380D13950DE1C102` |
| IMG | 551,779,200 B | `B737FADF8EAB6C4712F763E142E08DAB` | `6AB276B6B3A79B64DD8F513F528750EC19ED7999737FEEB5855CC0CFEFF27AAE` |
| SUB | 22,521,600 B | `F6D739F1B66082F7F06F403D30CF35EB` | `3C05EFEC20E9BE2CDA81500178FF1F8A7657908FD699D2070827E58D700122D2` |

원본은 판본별 전체 MD5·SHA-256으로 검사하며, 결과 CCD·IMG·SUB 각각도 **선택한 판본의 전체 MD5·SHA-256 및 크기**가 모두 일치해야 파일을 확정합니다. 이전 IMG Adler-32 대체 검증은 제거했습니다. 기존 결과도 일본판인지 영문판인지 검증하며 서로 바꿔 쓰지 않습니다.

## 배포 정보와 검증 범위

사용자 제공 `1.01.zip`의 일본판·영문판 전용 xdelta 6개와 수정 CUE를 사용합니다. 일본판 IMG 패치 SHA-256은 `dbc804ecc343d7bf19e493208fb88449c957c82c7767eb85ac3c16412eea2064`, 영문판은 `f0912110ff2a932ebe4379f45a95ff003ae9bea7114a7a55dc66b77d9feebbae`입니다. 배포 내 영문 IMG 패치명은 `English_IMG_v1.01.xdelta`로 정리했습니다.

xdelta3는 [jmacd/xdelta v3.2.0](https://github.com/jmacd/xdelta/releases/tag/v3.2.0) 공식 Windows x64 배포본과 동일 소스의 Win32 빌드입니다. Apache License 2.0 라이선스와 수정 고지를 ZIP에 포함하며 Go 소스는 [src/mirrors-go](../src/mirrors-go/README.ko.md)에 관리합니다.

`Mirrors_Kor1.01.zip`: 18,919,281바이트, SHA-256 `4561080e082bc9006da8da789af1b79545c8a72c6f96a4e4ee8f872e4cc60646`. 로컬 Go 테스트 20개, Windows x64·x86 빌드, 패치·디코더·추가 자산 SHA-256 및 ZIP CRC 검사 통과. **실제 일본판·영문판 전체 원본 이미지가 없어 새 패치 결과의 종단 간 시험은 수행하지 못했습니다.** 위 패치 후 해시는 사용자에게 제공받은 CloneCD 측정값입니다.
