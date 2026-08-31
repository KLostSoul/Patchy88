# Patchy88 파생판 고지

Patchy88은 46OkuMen의 `romtools` 저장소에 포함된 Pachy98의 PC 게임 패치 배포 개념에서 출발한 PC-8801용 파생 프로그램입니다.

- Original project: 46OkuMen / romtools / Pachy98
- Original repository: https://github.com/46OkuMen/romtools
- Original license: Apache License 2.0

## Patchy88에서 변경/추가한 주요 사항

- NDC를 이용한 파일 추출/삭제/재삽입 방식 사용 안 함
- xdelta3 대신 IPS 직접 적용
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

## 배포 데이터

Patchy88 저장소 및 배포물에는 다음을 포함하지 않습니다.

- 원본 게임 D88 이미지
- 원본 KANJI ROM
- 패치 완료 전체 게임 이미지

IPS 패치, 검증 매니페스트, Patchy88 소스 및 배포 실행파일만 보관합니다.
