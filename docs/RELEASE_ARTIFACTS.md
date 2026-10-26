# Patchy88 배포 산출물

현재 개발 세션에서 생성·검증한 배포 ZIP의 식별값입니다.

| 대상 | 파일 | 크기 | SHA-256 |
|---|---|---:|---|
| 몽환전사 바리스 Windows x64/x86 | `Patchy88_Valis1_PC88_v1.0.7.zip` | 2,035,112 | `25537c9c97edebd1d958df0ffeb0067afa587ed227967d118f2f23dfcc0d5703` |
| 몽환전사 바리스 Python | `Patchy88_Valis1_PC88_Python_v1.0.7.zip` | 220,377 | `62b553cb010e86902d7edfca532e7a5a3c51ad5d3e434a2a351dc6ec944902e7` |
| 몽환전사 바리스 II Windows x64/x86 | `Patchy88_Valis2_PC88_v1.0.3.zip` | 2,127,737 | `4e7d3c488cfddabbb15e36c4580c4a974421cb7c52f783270f41c6c3d3365028` |
| Alltynex PatchyFM Windows x64/x86 | `PatchyFM_Alltynex_v1.0.zip` | 3,041,670 | `bfac1cd128d39b4bda6df4a2cb033289f3387df97e2da205ec9357631a89237f` |
| Yuuwaku PatchyFM Windows x64/x86 | `PatchyFM_Yuuwaku_v1.01.zip` | 2,588,171 | `0f7d6d216ccba39123aef43120be4cbcc99007318bbdbe5bca83a141cd334548` |
| Dune II PatchyDOS Windows x64/x86 | `PatchyDOS_Dune2_v1.01.zip` | 2,693,994 | `01e490a50e639fcb43c70392cfdd96af2dddd659f27319468ffab24cfddeb951` |

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


## PatchyFM Alltynex v1.0

- FM TOWNS판 《Alltynex》 원본 ZIP 하나를 입력으로 사용합니다.
- 원본 ZIP과 완성 ISO는 저장소 및 배포 ZIP에 포함하지 않습니다.
- 동봉 xdelta SHA-256: `d7ef05131f9f2f2809e73e7e77adf77f4ff19f89a2c59bcf6e604ec8f4cfc99b`
- 지원 원본 ZIP:
  - 크기 `574,208 bytes`
  - MD5 `bca57cbe9a2f23c922118f40fffae43a`
  - SHA-256 `19b662f038e50e99535f60dbe6e07a77e269f574ed9250d46fc0655bb2aa6b3a`
- 검증된 결과 ISO:
  - 파일명 `Alltynex (Kor v1.0).iso`
  - 크기 `5,222,400 bytes`
  - MD5 `d5a438b7c7aac9e5b08baea6cb361f70`
  - SHA-256 `114a9ffa82e32bcc8cacb8b0475810422056560260eae5697ccf32b17bc91a54`
- Windows x64/x86 실행파일과 공식 xdelta3 v3.2.0 도우미를 함께 배포합니다.
- 생성되는 ISO의 부팅 환경에는 FreeTOWNSOS / TSUGARU OS를 사용하며 원 라이선스 전문을 배포물에 포함합니다.
- ZIP 선택 시 종료되던 Win32 파일 선택 필터 문제를 수정한 v1.0 배포물입니다. 버전 번호는 올리지 않았습니다.
- 사용자 환경에서 **원본 ZIP 선택 → xdelta 적용 → ISO 전체 해시 검증 → 게임 실행**까지 종단 간 확인했습니다.

[PatchyFM Alltynex 상세 문서](PATCHYFM_ALLTYNEX.md)


## 이전 배포: PatchyFM_Yuuwaku_v1.0

- Windows x64/x86 실행파일 `PatchyFM-Yuuwaku-x64.exe`, `PatchyFM-Yuuwaku-x86.exe`를 포함합니다.
- 소스 별도 ZIP: `PatchyFM_Yuuwaku_v1.0_Source.zip` (14,193 bytes, SHA-256 `07cddaa6d6b9b8825253ca3a4ce67ca51756c292e590ba4463cd176c0728360b`).
- 원본 이미지 폴더를 선택하면 하위 폴더까지 검색해 크기, MD5, SHA-256으로 정확한 원본 IMG/BIN을 자동 인식합니다. 개별 파일 선택은 하지 않습니다.
- 한글판 `Temptation-Korean-v1.0.bin` 및 `.cue` 생성.
- 사용자 원본 `Yuuwaku.img`와 사용자 제공 xdelta로 복원한 한글판 BIN이 전체 바이트 및 MD5/SHA-256 일치했습니다.
- Go 단위 테스트, Windows x64/x86 빌드, 독립 디코더 기반 8건 통합 테스트를 통과했습니다. Windows GUI 직접 실행과 게임 부팅은 아직 별도로 확인하지 않았습니다.
- 원본 및 완성 게임 이미지는 미배포. xdelta 및 xdelta3 헬퍼는 배포 ZIP에만 포함하며 소스 트리에는 중복 보관하지 않습니다.

[PatchyFM Yuuwaku 상세 문서](PATCHYFM_YUUWAKU.md)


## PatchyFM Alltynex 이름 정리 (2026-10-19)

- 배포 ZIP: `PatchyFM_Alltynex_v1.0.zip`, 크기 3,041,670 bytes, SHA-256 `bfac1cd128d39b4bda6df4a2cb033289f3387df97e2da205ec9357631a89237f`.
- 소스 ZIP: `PatchyFM_Alltynex_v1.0_Source.zip`, 크기 15,321 bytes, SHA-256 `bf6833105be2f849eccb14b7f38940c105555ced0aa752b6ad23ddf70f409807`.
- Windows 실행파일: `PatchyFM-Alltynex-x64.exe`, `PatchyFM-Alltynex-x86.exe`.
- 기존 xdelta, 지원 원본과 결과 ISO 해시/파일명은 변경하지 않았으며, Go 테스트 및 Windows 양 아키텍처 재빌드 통과.
- 이 환경에서 Windows GUI 실제 구동과 게임 실행은 추가로 확인하지 않았습니다.


## PatchyFM_Yuuwaku_v1.01 (2026-10-20)

- Windows 배포 ZIP: `PatchyFM_Yuuwaku_v1.01.zip`, 2,588,171 bytes, SHA-256 `0f7d6d216ccba39123aef43120be4cbcc99007318bbdbe5bca83a141cd334548`.
- Go 소스 ZIP: `PatchyFM_Yuuwaku_v1.01_Source.zip`, 19,394 bytes, SHA-256 `5ad5643d645554717872c80b2c70517d7cacee9c873e7229b0f0d3d1dbfb942e`.
- 업데이트된 xdelta 330,728 bytes, SHA-256 `c936529c182cd2854644b5a839406937628437325ef161f9af6c77cf36919e09`.
- 생성 BIN: `Temptation-Korean-v1.01.bin` (7,599,312 bytes), SHA-256 `1b456c4ff7d78c7f144986c6714f07a6fd62ecdd5b5c29c403656850653447c2`.
- 생성 CUE: `Temptation-Korean-v1.01.cue`, 89 bytes, SHA-256 `bf7336f007be2f8a0c1b3cb2d1a2acda8677200fab41ee501271d58e914c5998`.
- v1.0 대비 결과 BIN 41,809 bytes 변경. 기존 폴더 지정·자동 검색·전체 검증 로직은 유지합니다.
- Go 단위 테스트와 Windows x64/x86 빌드, 별도 디코더 기반 9개 Linux 통합 테스트 성공. Windows GUI 직접 실행 및 Tsugaru 부팅은 별도로 검증하지 않았습니다.
- 원본/완성 게임 디스크는 Git과 배포물 모두에 포함하지 않으며, xdelta와 디코더는 배포 ZIP에 포함하고 소스 트리에는 중복하지 않습니다.

[상세 사용법](PATCHYFM_YUUWAKU.md)


## PatchyDOS_Dune2_v1.01 (2026-10-26)

- Windows 배포 ZIP: `PatchyDOS_Dune2_v1.01.zip`, 2,693,994 bytes, SHA-256 `01e490a50e639fcb43c70392cfdd96af2dddd659f27319468ffab24cfddeb951`.
- 소스 ZIP: `PatchyDOS_Dune2_v1.01_Source.zip`, 40,531 bytes, SHA-256 `793529354fd51db1743178762041acb7f04689bed21b90fefa7a961545a0a092`.
- 지원 원본 IMG/BIN: 15,664,320 bytes (6,660 MODE1/2352 섹터), SHA-256 `edcdce3de2e4d986421d99cc11be424277317fc3e8cd71cf78f8ed544bce81bf`.
- v1.01 xdelta: `Dune2-Battle-for-Arrakis-Korean-v1.01.bin.xdelta`, 356,319 bytes, SHA-256 `41ef92e9a9cf062109371abce1e591fbe331e7b184b6614f6eedc1f74acef9db`.
- 결과 BIN: `Dune2-Battle-for-Arrakis-Korean-v1.01.bin`, 16,835,616 bytes (7,158 섹터), MD5 `1cd47830e0925747622bb53e12590ae2`, SHA-256 `5f3877d5d6be7b54b525ec91df19a46b01c4a6f5e83e3df6b683f05203053752`.
- 결과 CUE: `Dune2-Battle-for-Arrakis-Korean-v1.01.cue`, 103 bytes, SHA-256 `3d26796ff3eefd265b4f85f80a0ac6bf0ee95ae1f4f0c593fb93c8e009b3834b`.
- 폴더 자동 식별과 전체 무결성 검사 유지. xdelta3 실행 시간 최대 5분, 90초 무활동 감지, 5초 간격 진행 기록, 디코더 stdout/stderr 로깅, GUI 중단 및 실패/취소 시 임시 파일 정리 추가.
- Go 단위 테스트 및 Windows x64/x86 컴파일 통과. 실제 원본으로 Linux CLI 종단 간 검증 및 독립 디코더 결과 전체 바이트 일치 확인. Windows GUI 및 DOS 게임 실행은 아직 별도 검증이 필요합니다.
- 글꼴 고지와 전문: 달무리 Apache-2.0, GNU Unifont OFL-1.1 / GPLv2 이상(폰트 임베딩 예외).
- 원본 게임 CD, 패치 완료 이미지, Windows EXE 및 xdelta는 Git 소스 트리에 저장하지 않습니다. 배포 ZIP과 소스 ZIP을 별도로 보관합니다. GitHub Releases 등록은 별도 작업입니다.

[PatchyDOS Dune II 상세 문서](PATCHYDOS_DUNE2.md)
