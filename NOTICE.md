# Patchy88 프로젝트 고지

Patchy88은 PC-8801용 패치 프로그램 모음이며, 같은 안전 검증 원칙을 이용한 FM TOWNS용 변형 PatchFM도 이 저장소에서 관리합니다. Pachy98/romtools (46OkuMen)의 패치 배포 개념을 참고했습니다.

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

## PatchFM — Alltynex Kor v1.0

PatchFM Alltynex는 사용자가 준비한 원본 FM TOWNS 《Alltynex》 ZIP 파일 하나를 입력으로 받아 xdelta를 적용하고 부팅 가능한 한글판 ISO를 생성합니다.

배포물에는 다음을 **포함하지 않습니다.**

- 원본 `alltynex_fmtowns.zip`
- 완성 `Alltynex (Kor v1.0).iso`

배포물에는 PatchFM 실행파일, xdelta 패치, xdelta3, 매니페스트 및 라이선스/고지 문서만 포함합니다.

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
- 원 프로젝트의 `LICENSE.md` 전문을 PatchFM 배포물과 `src/patchfm-alltynex/FREETOWNSOS_LICENSE.md`에 포함

FreeTOWNSOS 라이선스에 포함된 ORICON 및 Free386 관련 고지도 함께 유지합니다.

[PatchFM Alltynex 사용법과 검증 정보](docs/PATCHFM_ALLTYNEX.md)
