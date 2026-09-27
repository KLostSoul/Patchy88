# Patchy88 배포 산출물

현재 개발 세션에서 생성·검증한 배포 ZIP의 식별값입니다.

| 대상 | 파일 | 크기 | SHA-256 |
|---|---|---:|---|
| 몽환전사 바리스 Windows x64/x86 | `Patchy88_Valis1_PC88_v1.0.7.zip` | 2,035,112 | `25537c9c97edebd1d958df0ffeb0067afa587ed227967d118f2f23dfcc0d5703` |
| 몽환전사 바리스 Python | `Patchy88_Valis1_PC88_Python_v1.0.7.zip` | 220,377 | `62b553cb010e86902d7edfca532e7a5a3c51ad5d3e434a2a351dc6ec944902e7` |
| 몽환전사 바리스 II Windows x64/x86 | `Patchy88_Valis2_PC88_v1.0.3.zip` | 2,127,737 | `4e7d3c488cfddabbb15e36c4580c4a974421cb7c52f783270f41c6c3d3365028` |

## 바리스 II v1.0.3

- Disk A~G + KANJI1을 사용자 제공 V1.01 IPS 세트로 교체
- 배포 ZIP 내부에서도 `Valis2_KOR_Disk_A_V1.01.ips` 형식의 원래 파일명을 그대로 유지
- A~F와 KANJI1은 v1.0.2 대비 IPS 내용 변경, G는 동일
- 8개 대상의 before/after 영역 SHA-256 매니페스트 전면 재생성
- KANJI1 기준 원본은 CRC32 `6178BD43` 유지
- 8개 원본 교차 자동식별 오인 없음 확인
- 8개 새 패치 결과 모두 `ALREADY_PATCHED` 판정 확인

## 배포 원칙

- Windows판 ZIP에는 x64/x86 실행파일을 함께 둡니다.
- 패치 파일명에 포함된 버전 표기는 임의로 제거하거나 변경하지 않습니다.
- 원본 게임 D88, 원본 KANJI ROM, 패치 완료 전체 게임 이미지는 배포하지 않습니다.
- 실제 사용자 배포 ZIP은 GitHub Releases 자산으로 두고, 소스 트리에는 중복 보관하지 않는 것을 원칙으로 합니다.

## Mirrors용 Patchy88 v1.01 (공통 한글판 결과)

- 사용자 제공 `1.01.zip`의 일본판·영문판 xdelta 6개 및 CUE를 유지합니다. 일본판 IMG 패치 SHA-256 `dbc804ecc343d7bf19e493208fb88449c957c82c7767eb85ac3c16412eea2064`, 영문판 IMG 패치 SHA-256 `f0912110ff2a932ebe4379f45a95ff003ae9bea7114a7a55dc66b77d9feebbae`.
- 두 원본 판본의 최종 CCD·IMG·SUB 해시는 동일합니다. IMG 기대 MD5는 `56E768F7CE3315A8172338CB10CE153E`입니다.
- 결과 CCD 3,500 B, IMG 551,779,200 B, SUB 22,521,600 B. 각 파일의 전체 MD5·SHA-256과 크기를 필수 검증합니다. 기존 Adler-32 대체 검증은 제거했습니다.
- 일본판 또는 영문판 하나만 검출되면 자동 선택하며 두 판본이 있으면 사용자가 선택합니다. 정상적인 공통 한글판 결과를 재검사할 수 있습니다.
- Windows x64·x86 새 실행파일, 공식 xdelta3 v3.2.0 디코더, 라이선스 및 수정 고지문 포함. 배포 ZIP에 Go 소스는 넣지 않습니다.
- Go 테스트와 Windows x64·x86 빌드를 통과했습니다. 실제 영문판 원본 CCD·IMG·SUB로 세 xdelta 적용과 공통 한글판 결과 해시를 종단 간 확인했습니다.

[Mirrors용 Patchy88 사용법과 모든 판본별 해시](MIRRORS.md)

- 수정 전에는 영문판 최종 해시를 잘못 분리해 IMG 검증이 실패했습니다. 영문판 CCD·IMG·SUB xdelta 출력은 모두 공통 한글판 결과임을 실제 파일로 확인하고, 매니페스트와 기존 결과 판정을 수정했습니다.
