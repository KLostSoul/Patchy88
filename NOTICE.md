# Patchy88 프로젝트 고지

Patchy88은 46OkuMen의 `romtools` 저장소에 포함된 Pachy98의 PC 게임 패치 배포 개념에서 출발한 PC-8801용 파생 프로그램입니다.

- Original project: 46OkuMen / romtools / Pachy98
- Original repository: https://github.com/46OkuMen/romtools
- Original license: Apache License 2.0

## Valis용 IPS Patchy88의 변경/추가 사항

- NDC를 이용한 파일 추출/삭제/재삽입 방식 사용 안 함
- xdelta3 대신 IPS 직접 적용 (Valis용 IPS 패처에 한함)
- PC-8801 D88/ROM 파일 오프셋 직접 패치
- IPS 대상 레코드별 적용 전/후 SHA-256 검증
- `ORIGINAL / ALREADY_PATCHED / PARTIAL / INCOMPATIBLE` 상태 판정
- 파일명 대신 패치 대상 데이터로 원본 식별
- 같은 폴더 `.bak` 백업
- 결과 파일명에 `(K)` 접미사 적용
- 임시파일 패치 및 사후검증 후 결과 확정
- 몽환전사 바리스 Disk A/KANJI1 지원
- 몽환전사 바리스 II Disk A~G/KANJI1 폴더 자동 식별 및 일괄 트랜잭션 처리
- Windows x64/x86 네이티브 배포판 및 바리스 1 Python판 제공

## Mirrors용 Patchy88 (PC-8801)

Mirrors_Kor1.00은 일본판과 영문판의 CCD/IMG/SUB에 각각 전용 xdelta 패치를 적용하는 별도 패처입니다.

- 일본판과 영문판을 파일 전체 MD5 및 SHA-256으로 식별합니다.
- 두 판본이 같은 폴더에 있으면 사용자가 패치할 판본을 선택합니다.
- 결과 파일명: `Mirrors_Kor1.00.ccd`, `Mirrors_Kor1.00.img`, `Mirrors_Kor1.00.sub`, `Mirrors_Kor1.00.cue`.
- xdelta3.exe를 함께 사용하며 기존 IPS 자체 적용 엔진과는 다릅니다.
- 결과 검증 성공 후 D88 2개 및 CUE를 입력 폴더에 추가하며 원본은 덮어쓰지 않습니다.

공식 xdelta: https://github.com/jmacd/xdelta
Apache 2.0으로 재라이선스된 3.0.11 소스: https://github.com/jmacd/xdelta/tree/release3_0_apl
라이선스 원문: https://github.com/jmacd/xdelta/blob/release3_0_apl/xdelta3/LICENSE

동봉된 xdelta3.exe는 Pachy98 v0.20.1 제공본과 일치하지만, 공식 Apache 브랜치에서 빌드된 실행파일인지까지는 확인되지 않았습니다. 배포 ZIP에는 라이선스 원문과 수정·출처 고지문을 포함합니다.

[Mirrors용 Patchy88 문서](docs/MIRRORS.md)를 참조하십시오.

## 배포 데이터

Patchy88 저장소 및 배포물에는 다음을 포함하지 않습니다.

- 원본 게임 D88 이미지
- 원본 KANJI ROM
- 패치 완료 전체 게임 이미지

IPS/xdelta 패치, 검증 매니페스트, Patchy88 소스 및 배포 실행파일을 별도로 관리합니다. 게임 원본과 패치 완료 전체 이미지는 배포하지 않습니다.
