# Mirrors용 Patchy88 — Mirrors_Kor1.01

대상 기종: **PC-8801** | 패치 방식: **xdelta3 (VCDIFF)** | 지원 원본: 일본판 및 영문판

**Mirrors용 Patchy88**은 PC-8801 《Mirrors》의 일본판 또는 영문판 이미지에 해당 판본의 전용 xdelta 패치 3개를 적용하여 동일한 한글판 CCD/IMG/SUB를 만듭니다. Valis용 IPS Patchy88과는 패치 엔진·원본 검사 방식·출력 방식이 다릅니다.

## 사용 방법

1. 패처 ZIP을 풀고 `Mirrors_Kor1.01-x64.exe` 또는 `Mirrors_Kor1.01-x86.exe`를 실행합니다. `assets/` 폴더 및 그 안의 `xdelta3.exe`를 실행파일과 함께 유지합니다. xdelta3를 별도로 실행할 필요는 없습니다.
2. 일본판 또는 영문판의 `.ccd`, `.img`, `.sub`가 있는 **원본 폴더**를 선택합니다. 판본 식별은 파일명이 아닌 전체 MD5 + SHA-256으로 수행합니다.
3. 한 판본만 완전하게 검출되면 자동으로 선택합니다. 일본판과 영문판이 **둘 다 검출되면 선택 창이 나타납니다**. 예(Y)는 일본판, 아니오(N)는 영문판, 취소는 중단입니다.
4. 선택한 판본을 확인한 뒤 패치를 적용합니다. 다른 판본의 원본은 사용하거나 변경하지 않습니다.
5. xdelta 패치 결과 CCD/SUB는 MD5·SHA-256으로, IMG는 VCDIFF 윈도우 체크섬과 전체 크기로 검증한 다음, **같은 원본 폴더**에 결과와 CUE, D88 두 개를 추가합니다. 오류 또는 기존 파일명 충돌이 있으면 확정하지 않습니다.

## 출력 파일

| 파일명 | 생성 방법 |
|---|---|
| `Mirrors_Kor1.01.ccd` | 선택한 판본의 CCD에 전용 xdelta 적용 |
| `Mirrors_Kor1.01.img` | 선택한 판본의 IMG에 전용 xdelta 적용 |
| `Mirrors_Kor1.01.sub` | 선택한 판본의 SUB에 전용 xdelta 적용 |
| `Mirrors_Kor1.01.cue` | 동봉 CUE를 검증 후 복사, 내부 `FILE "Mirrors_Kor1.01.img" BINARY` |
| `disk1main.d88` | 동봉 이미지 검증 후 복사 |
| `disk2game.d88` | 동봉 이미지 검증 후 복사 |

**주의:** 제공된 `disk1main.d88`과 `disk2game.d88`은 명칭은 다르지만 SHA-256이 동일합니다. 패처는 사용자 요청에 따라 두 파일을 모두 추가합니다. 서로 다른 디스크 내용임을 입증하는 자료로 해석하지 마십시오.

## 원본 판본 구분 (전체 MD5)

| 파일 | 일본판 MD5 | 영문판 MD5 |
|---|---|---|
| CCD | `80273154A2DAF2D107A282B353A86246` | `35C733769D60277FCCE522E121AF82AE` |
| IMG | `84E583EC62372CA9DD0D737DCE38B5D8` | `F3BC8D9E2650D9D9C6EBAEF5E2353144` |
| SUB | `E45923398B1D150F71E6056BD0974D15` | `F6D739F1B66082F7F06F403D30CF35EB` |

## 한글판 결과 MD5 / SHA-256 (CCD/SUB만 기준값 확정)

| 파일 | MD5 | SHA-256 |
|---|---|---|
| CCD | `80273154A2DAF2D107A282B353A86246` | `5512EFCEC20279753E1C6FEF0C5AAE37F2F3F05BB3AE3900B5DB0F688D4A67AC` |
| IMG | 미확정 (1.01) | 미확정 (1.01) |
| SUB | `E45923398B1D150F71E6056BD0974D15` | `C20C93FF3A2CF3A1E3A18B5F8046EC6F6198A916B417A01678B0B9D10EED9046` |

일본판의 CCD와 SUB는 한글판 결과와 해시가 동일합니다. 영문판에서는 CCD·IMG·SUB 모두 결과 해시가 다릅니다. 두 판본별로 제공된 xdelta 패치 3개를 각각 사용합니다. 원본 전체는 MD5·SHA-256으로 검사하지만, 1.01 IMG는 새 전체 해시 기준값이 없어 VCDIFF 윈도우 66개의 Adler-32와 출력 크기로 별도 검증합니다.

## 안전성 및 검증 범위

원본은 변경하지 않으며, 다른 내용의 결과 파일이 이미 있다면 덮어쓰지 않습니다. 임시 결과를 만든 뒤 CCD/SUB의 MD5·SHA-256 및 IMG의 66개 윈도우 Adler-32와 크기를 검증한 뒤 파일명을 확정합니다. 이미 같은 내용의 한글판 CCD/IMG/SUB가 있으면 추가 파일만 확인·복사할 수 있습니다.

소스 코드와 매니페스트는 [src/mirrors-go](../src/mirrors-go/README.ko.md)에 보관합니다. 사용자 배포 ZIP에는 Go 소스를 넣지 않습니다. 현재 검증 기록에는 판본 식별·선택·충돌 방지·자산 변조 검사 테스트가 포함되며, **실제 일본판·영문판 전체 CD 원본으로 수행하는 최종 xdelta 적용 시험은 아직 완료되지 않았습니다.**

## 공식 xdelta3 v3.2.0 및 라이선스

이 배포판은 **기존 Pachy98 동봉 xdelta3.exe를 사용하지 않습니다.** 공식 [jmacd/xdelta v3.2.0](https://github.com/jmacd/xdelta/releases/tag/v3.2.0)의 Apache License 2.0 소스를 사용합니다.

- Windows x64: 공식 [v3.2.0 Windows x64 릴리스 ZIP](https://github.com/jmacd/xdelta/releases/download/v3.2.0/xdelta3-3.2.0-windows-x86_64.zip)을 다운로드하여 SHA-256 `af8ef036cb077a48df080c9a8ac1be4a6e7511c32d11f8bec89b6803a9e52576` 검증 후 실행파일을 추출했습니다. `assets/xdelta3-x64.exe` SHA-256: `53d90226615f217d3380c39892833311b4e24acd863e1ca01f14b5e772e2e6d0`.
- Windows x86: 동일한 공식 `v3.2.0` 태그를 사용해 [GitHub Actions Windows Win32 빌드](https://github.com/KLostSoul/Patchy88/actions/runs/36051984391)를 수행했습니다. `assets/xdelta3.exe` SHA-256: `232a8e8ac9fb47a54d0ca4d6acdb322e456a8c72cbfe3b212cf2f7e498760481`.
- x64/x86 실행파일은 매니페스트에 별도의 SHA-256이 기록되며 패처가 실행 환경에 맞춰 자동 선택합니다. 별도로 xdelta를 설치할 필요는 없습니다.

ZIP에는 `LICENSE`, `THIRD_PARTY_LICENSE.txt`, `NOTICE_MODIFICATIONS.txt`, `UPSTREAM_XDELTA.txt`를 포함하며, 소스는 [저장소](../src/mirrors-go/README.ko.md)에 별도 관리합니다. 두 판본 동시 검출·선택 및 디코더 사전검증을 포함하여 Go 테스트 12개를 통과했습니다. 실제 일본판과 영문판 전체 원본에 대한 최종 xdelta 적용 시험은 아직 미실시입니다.

## v1.01 업데이트 (사용자 제공 ZIP 기준)

일본판/영문판 각 CCD·IMG·SUB의 `*_v1.01.xdelta` 여섯 개와 `Mirrors_Kor_v1.01.cue`를 반영했습니다. 게임 원본을 덮어쓰지 않고 출력 파일명을 `Mirrors_Kor1.01`으로 통일합니다. 업로드된 CUE에 남아 있던 예전 IMG 참조도 `Mirrors_Kor1.01.img`로 바꿨습니다.

새 IMG는 이전 결과와 다른 데이터이므로 **1.00 IMG의 MD5·SHA-256을 재사용하지 않습니다.** v1.01의 일본판과 영문판 IMG 패치에서 추출한 66개 출력 윈도우 Adler-32가 모두 동일하며, 전체 출력 크기는 551,779,200바이트입니다. 실제 새 IMG를 만든 뒤 전체 MD5·SHA-256을 계산해 로그에 표시하지만, 업로드에 해당 기준 해시와 실제 원본 CD 이미지가 없으므로 현재는 결과 전체에 대한 암호학적 기준값 검증과 실원본 종단 간 시험을 완료하지 못했습니다.

배포 ZIP: `Mirrors_Kor1.01.zip`. 공식 xdelta v3.2.0 Windows x64 및 Win32 디코더, `LICENSE`, `THIRD_PARTY_LICENSE.txt`, `NOTICE_MODIFICATIONS.txt`를 포함하고 Go 소스는 포함하지 않습니다.
