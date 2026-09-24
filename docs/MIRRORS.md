# Mirrors용 Patchy88 — Mirrors_Kor1.00

대상 기종: **PC-8801** | 패치 방식: **xdelta3 (VCDIFF)** | 지원 원본: 일본판 및 영문판

**Mirrors용 Patchy88**은 PC-8801 《Mirrors》의 일본판 또는 영문판 이미지에 해당 판본의 전용 xdelta 패치 3개를 적용하여 동일한 한글판 CCD/IMG/SUB를 만듭니다. Valis용 IPS Patchy88과는 패치 엔진·원본 검사 방식·출력 방식이 다릅니다.

## 사용 방법

1. 패처 ZIP을 풀고 `Mirrors_Kor1.00-x64.exe` 또는 `Mirrors_Kor1.00-x86.exe`를 실행합니다. `assets/` 폴더 및 그 안의 `xdelta3.exe`를 실행파일과 함께 유지합니다. xdelta3를 별도로 실행할 필요는 없습니다.
2. 일본판 또는 영문판의 `.ccd`, `.img`, `.sub`가 있는 **원본 폴더**를 선택합니다. 판본 식별은 파일명이 아닌 전체 MD5 + SHA-256으로 수행합니다.
3. 한 판본만 완전하게 검출되면 자동으로 선택합니다. 일본판과 영문판이 **둘 다 검출되면 선택 창이 나타납니다**. 예(Y)는 일본판, 아니오(N)는 영문판, 취소는 중단입니다.
4. 선택한 판본을 확인한 뒤 패치를 적용합니다. 다른 판본의 원본은 사용하거나 변경하지 않습니다.
5. xdelta 패치 결과 CCD/IMG/SUB의 MD5 + SHA-256을 모두 검증한 다음, **같은 원본 폴더**에 결과와 CUE, D88 두 개를 추가합니다. 오류 또는 기존 파일명 충돌이 있으면 확정하지 않습니다.

## 출력 파일

| 파일명 | 생성 방법 |
|---|---|
| `Mirrors_Kor1.00.ccd` | 선택한 판본의 CCD에 전용 xdelta 적용 |
| `Mirrors_Kor1.00.img` | 선택한 판본의 IMG에 전용 xdelta 적용 |
| `Mirrors_Kor1.00.sub` | 선택한 판본의 SUB에 전용 xdelta 적용 |
| `Mirrors_Kor1.00.cue` | 동봉 CUE를 검증 후 복사, 내부 `FILE "Mirrors_Kor1.00.img" BINARY` |
| `disk1main.d88` | 동봉 이미지 검증 후 복사 |
| `disk2game.d88` | 동봉 이미지 검증 후 복사 |

**주의:** 제공된 `disk1main.d88`과 `disk2game.d88`은 명칭은 다르지만 SHA-256이 동일합니다. 패처는 사용자 요청에 따라 두 파일을 모두 추가합니다. 서로 다른 디스크 내용임을 입증하는 자료로 해석하지 마십시오.

## 원본 판본 구분 (전체 MD5)

| 파일 | 일본판 MD5 | 영문판 MD5 |
|---|---|---|
| CCD | `80273154A2DAF2D107A282B353A86246` | `35C733769D60277FCCE522E121AF82AE` |
| IMG | `84E583EC62372CA9DD0D737DCE38B5D8` | `F3BC8D9E2650D9D9C6EBAEF5E2353144` |
| SUB | `E45923398B1D150F71E6056BD0974D15` | `F6D739F1B66082F7F06F403D30CF35EB` |

## 한글판 결과 MD5 / SHA-256

| 파일 | MD5 | SHA-256 |
|---|---|---|
| CCD | `80273154A2DAF2D107A282B353A86246` | `5512EFCEC20279753E1C6FEF0C5AAE37F2F3F05BB3AE3900B5DB0F688D4A67AC` |
| IMG | `C637C340C048D4209A5B47AC8DF0E276` | `8EE893AFA18D98991ED78CA9F7F7259186A91FDED222058C2F4AF85E647A53B2` |
| SUB | `E45923398B1D150F71E6056BD0974D15` | `C20C93FF3A2CF3A1E3A18B5F8046EC6F6198A916B417A01678B0B9D10EED9046` |

일본판의 CCD와 SUB는 한글판 결과와 해시가 동일합니다. 영문판에서는 CCD·IMG·SUB 모두 결과 해시가 다릅니다. 두 판본별로 제공된 xdelta 패치 3개를 각각 사용하며, 원본과 결과 검증에 MD5와 SHA-256을 함께 사용합니다.

## 안전성 및 검증 범위

원본은 변경하지 않으며, 다른 내용의 결과 파일이 이미 있다면 덮어쓰지 않습니다. 임시 결과를 만든 뒤 전체 MD5 + SHA-256 검증에 성공한 경우에만 파일명을 확정합니다. 이미 같은 내용의 한글판 CCD/IMG/SUB가 있으면 추가 파일만 확인·복사할 수 있습니다.

소스 코드와 매니페스트는 [src/mirrors-go](../src/mirrors-go/README.ko.md)에 보관합니다. 사용자 배포 ZIP에는 Go 소스를 넣지 않습니다. 현재 검증 기록에는 판본 식별·선택·충돌 방지·자산 변조 검사 테스트가 포함되며, **실제 일본판·영문판 전체 CD 원본으로 수행하는 최종 xdelta 적용 시험은 아직 완료되지 않았습니다.**

## 배포판 라이선스 및 xdelta 출처

`Mirrors_Kor1.00.zip`에는 `LICENSE`, `NOTICE_MODIFICATIONS.txt`, `THIRD_PARTY_LICENSE.txt`가 들어 있습니다. 두 라이선스 전문은 [xdelta 공식 저장소의 Apache 2.0 원문](https://github.com/jmacd/xdelta/blob/release3_0_apl/xdelta3/LICENSE)과 바이트 단위로 동일합니다 (Git blob SHA-1: `7a774156a6820befeca5067631b3a50b86481495`).

공식 xdelta 저장소는 3.0.11 기반의 [`release3_0_apl` 브랜치](https://github.com/jmacd/xdelta/tree/release3_0_apl)를 Apache License 2.0으로 재라이선스한 소스로 명시합니다. 별도로 원래 GPL 라이선스를 적용한 [xdelta-gpl](https://github.com/jmacd/xdelta-gpl) 저장소도 존재합니다.

동봉 `xdelta3.exe`는 제공된 Pachy98 v0.20.1 배포본과 동일한 바이너리 (SHA-256: `9bf8d067de9448e521afe1f8108caa0f85b4b7c7933641efd44bc43533920565`)입니다. **그 실행파일의 실제 빌드 출처는 확인되지 않았으므로 공식 Apache 2.0 브랜치에서 빌드됐다고 단정하지 않습니다.** 이 제한을 배포 ZIP의 수정 고지문에도 명시했습니다.

소스 코드는 [Mirrors용 Patchy88 저장소](../src/mirrors-go/README.ko.md)에서 관리하며 배포 ZIP에는 포함하지 않습니다.
