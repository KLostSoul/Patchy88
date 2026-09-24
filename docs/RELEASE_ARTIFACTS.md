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

## Mirrors용 Patchy88 (PC-8801)

- 제품명: **Mirrors용 Patchy88** / 패키지 버전: `Mirrors_Kor1.00`.
- 수정판 ZIP 이름: `Mirrors_Kor1.00_PC8801_SELECT.zip` (별도 배포 파일; GitHub 소스 트리에 실행파일·게임 원본을 넣지 않음).
- Windows x64/x86 실행파일, xdelta3.exe, 일본판·영문판 전용 xdelta 6개, D88 2개, `Mirrors_Kor1.00.cue` 포함. 배포 ZIP에는 Go 소스를 넣지 않음.
- 두 판본 동시 검출 시 사용자가 적용할 판본 선택. 원본 및 결과 전체 MD5 + SHA-256 검증.
- 한글판 파일명은 `Mirrors_Kor1.00.ccd`, `Mirrors_Kor1.00.img`, `Mirrors_Kor1.00.sub`, `Mirrors_Kor1.00.cue`로 통일.
- 이번 GitHub 문서 갱신에서는 ZIP의 실제 크기와 SHA-256을 다시 측정하지 않았으므로 미기재. 저장소에 업로드되지 않은 Releases 자산이 있다고 가정하지 않음.
- 일본판 및 영문판 실제 원본 CD 이미지 전체로 수행하는 최종 xdelta 적용 시험은 아직 완료되지 않음.

자세한 내용은 [Mirrors용 Patchy88 사용법](MIRRORS.md)을 참조하십시오.
