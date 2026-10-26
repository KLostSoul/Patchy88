# Patchy88 프로젝트 고지

Patchy88은 PC-8801용 패치 프로그램 모음이며, 같은 안전 검증 원칙을 이용한 FM TOWNS용 변형 PatchyFM도 이 저장소에서 관리합니다. Pachy98/romtools (46OkuMen)의 패치 배포 개념을 참고했습니다.

원본 게임 이미지·ROM·ZIP과 패치 완료 전체 게임 이미지/ISO는 저장소 및 배포물에 포함하지 않습니다.

## Mirrors용 Patchy88 — Mirrors_Kor1.01

일본판과 영문판은 서로 다른 원본 CCD·IMG·SUB를 사용하지만, 현재 패치 세트는 **두 판본 모두 동일한 최종 한글판 CCD·IMG·SUB**를 생성합니다.

- 일본판/영문판 원본 전체 MD5·SHA-256·크기 검증
- xdelta 패치 자체 SHA-256 검증
- 생성 결과 CCD·IMG·SUB 전체 MD5·SHA-256·크기 검증
- 결과 CCD: 3,500 bytes
- 결과 IMG MD5: `56E768F7CE3315A8172338CB10CE153E`
- 일본판과 영문판의 최종 결과 해시는 동일
- 실제 영문판 원본을 이용한 종단 간 xdelta 적용 및 최종 결과 검증 완료

이전 문서에 있던 “일본판과 영문판의 패치 후 결과 해시가 서로 다르다”는 설명은 구형 검증 정보이며 현재 구현과 맞지 않습니다.

디코더는 공식 [jmacd/xdelta v3.2.0](https://github.com/jmacd/xdelta/tree/v3.2.0) Windows x64 및 동일 공식 소스의 Win32 빌드를 사용합니다. xdelta3는 Apache License 2.0입니다.

[Mirrors 사용법과 검증 해시](docs/MIRRORS.md)

## PatchyFM — Alltynex Kor v1.0

PatchyFM Alltynex는 사용자가 준비한 원본 FM TOWNS 《Alltynex》 ZIP 파일 하나를 입력으로 받아 xdelta를 적용하고 부팅 가능한 한글판 ISO를 생성합니다.

배포물에는 다음을 **포함하지 않습니다.**

- 원본 `alltynex_fmtowns.zip`
- 완성 `Alltynex (Kor v1.0).iso`

배포물에는 PatchyFM 실행파일, xdelta 패치, xdelta3, 매니페스트 및 라이선스/고지 문서만 포함합니다.

검증 절차:

```text
원본 ZIP 전체 크기/MD5/SHA-256 검증
→ xdelta 및 xdelta3 SHA-256 검증
→ 임시 ISO 생성
→ 결과 ISO 전체 크기/MD5/SHA-256 검증
→ 검증 성공 시에만 결과 ISO 확정
```

지원 원본:

- 크기: `574,208 bytes`
- MD5: `BCA57CBE9A2F23C922118F40FFFAE43A`
- SHA-256: `19B662F038E50E99535F60DBE6E07A77E269F574ED9250D46FC0655BB2AA6B3A`

검증된 결과:

- 파일명: `Alltynex (Kor v1.0).iso`
- 크기: `5,222,400 bytes`
- MD5: `D5A438B7C7AAC9E5B08BAEA6CB361F70`
- SHA-256: `114A9FFA82E32BCC8CACB8B0475810422056560260EAE5697CCF32B17BC91A54`

사용자 실기/에뮬레이션 환경에서 **원본 ZIP 선택 → 패치 적용 → ISO 생성 → 게임 실행**까지 확인했습니다.

### FreeTOWNSOS / TSUGARU OS

생성되는 ISO의 부팅 환경에는 CaptainYS(Soji Yamakawa)의 FreeTOWNSOS / TSUGARU OS를 사용합니다.

- Upstream: https://github.com/captainys/FreeTOWNSOS
- Copyright 2021 Soji Yamakawa
- 원 프로젝트의 `LICENSE.md` 전문을 PatchyFM 배포물과 `src/patchyfm-alltynex/FREETOWNSOS_LICENSE.md`에 포함

FreeTOWNSOS 라이선스에 포함된 ORICON 및 Free386 관련 고지도 함께 유지합니다.

[PatchyFM Alltynex 사용법과 검증 정보](docs/PATCHYFM_ALLTYNEX.md)


## PatchyFM_Yuuwaku_v1.0 — FM TOWNS 《誘惑 / Yuuwaku》 (이전 버전)

사용자가 지정한 폴더와 하위 폴더에서 원본 IMG/BIN 데이터 트랙을 크기·MD5·SHA-256으로 자동 식별하고, 제공된 xdelta로 한글판 BIN/CUE를 생성합니다. CUE/CCD/SUB를 개별 선택하지 않습니다.

- 원본 IMG: 6,703,200 bytes, MD5 `000b8a4e832d952fb29ea6b831ed8cd0`
- 한글판 BIN: 7,599,312 bytes, MD5 `e07f8348b3a3b48f2e02613edeb204db`
- xdelta SHA-256: `01c25b644e222bdb432fe484f6faed8a20714cce571ed2989d8bba79ead42e1e`
- 원본 IMG/CCD/CUE/SUB 및 완성 BIN은 저장소와 배포물에 포함하지 않습니다.
- Windows x64/x86 컴파일, Go 테스트 및 대체 디코더를 사용한 폴더 모드 통합 테스트 8건 통과. Windows GUI 직접 실행 및 Tsugaru 실제 부팅은 미확인.
- FreeTOWNSOS는 사용하지 않으며 포함되지 않습니다.
- 패처 코드는 Apache-2.0; xdelta3의 상위 프로젝트는 jmacd/xdelta v3.2.0입니다.

[사용 설명 및 전체 해시](docs/PATCHYFM_YUUWAKU.md)


## PatchyFM_Yuuwaku_v1.01 — FM TOWNS 《誘惑 / Yuuwaku》 (현재 배포)

- 수정 xdelta `Temptation-Korean-v1.01.xdelta`: 330,728 bytes, SHA-256 `c936529c182cd2854644b5a839406937628437325ef161f9af6c77cf36919e09`.
- 적용 결과 `Temptation-Korean-v1.01.bin`: 7,599,312 bytes, SHA-256 `1b456c4ff7d78c7f144986c6714f07a6fd62ecdd5b5c29c403656850653447c2`.
- v1.0 출력과 비교하여 41,809 bytes 변경. 새 BIN/CUE 파일명으로 이전 출력과 충돌하지 않습니다.
- 사용자 지정 폴더 및 하위 디렉터리에서 크기·MD5·SHA-256으로 원본 IMG/BIN 자동 인식; 원본·결과 전체 검증, 중복·변조·덮어쓰기 보호.
- 제공받은 `Temptation-Korean-v1.0(1).xdelta`를 수정 없이 배포용 이름으로 변경했으며, 디코딩 결과의 Adler-32와 MD5/SHA-256을 확인했습니다.
- Go 단위 테스트, Windows x64/x86 크로스 빌드, Linux 환경의 9개 통합 테스트 통과. Windows GUI 및 실제 Tsugaru 실행은 아직 검증되지 않았습니다.
- 원본 게임 이미지 및 완성 BIN/ISO는 배포하지 않습니다. 패치 파일과 xdelta3 실행파일은 배포 ZIP에만 포함하며 Git 소스 트리에는 두지 않습니다.

[PatchyFM Yuuwaku v1.01 사용법](docs/PATCHYFM_YUUWAKU.md)


## PatchyDOS_Dune2_v1.01 — MS-DOS 《Dune II: The Battle for Arrakis》

PatchyDOS는 폴더 및 하위 폴더에서 지원 원본 IMG/BIN을 크기·MD5·SHA-256으로 식별하고 xdelta를 적용해 한글판 CD 이미지와 CUE를 생성합니다. 원본 이미지를 수정하지 않으며 기존 결과물이 다를 때 덮어쓰지 않습니다.

- 원본: 15,664,320 bytes, SHA-256 `edcdce3de2e4d986421d99cc11be424277317fc3e8cd71cf78f8ed544bce81bf`.
- v1.01 xdelta: 356,319 bytes, SHA-256 `41ef92e9a9cf062109371abce1e591fbe331e7b184b6614f6eedc1f74acef9db`.
- 생성 BIN: 16,835,616 bytes, SHA-256 `5f3877d5d6be7b54b525ec91df19a46b01c4a6f5e83e3df6b683f05203053752`.
- 전체 디코더 실행 제한 5분, 활동 없음 90초 중단, 5초 간격 경과/출력 크기 로그, 디코더 오류 출력 및 사용자 취소, 임시 파일 정리.
- Windows x64/x86 빌드 및 Go 테스트, Linux 기반 xdelta 복원과 통합 검사 통과. Windows GUI 구동과 DOS 게임 실행 자체는 아직 직접 검증하지 않았습니다.
- 게임 원본 IMG/CUE/CCD/SUB와 완성 BIN은 소스/배포 ZIP에 포함하지 않습니다. xdelta 및 디코더 실행파일은 배포 ZIP에만 포함합니다.
- `FONT_ATTRIBUTIONS.md`: [달무리 글꼴](https://github.com/RanolP/dalmoori-font) Apache-2.0, [GNU Unifont](https://unifoundry.com/unifont/index.html) SIL OFL 1.1 또는 GPL v2 이상(폰트 임베딩 예외). 세부 고지 및 라이선스 전문은 `src/patchydos-dune2/LICENSES/` 참조.

[PatchyDOS Dune II 사용법 및 검증 정보](docs/PATCHYDOS_DUNE2.md)
