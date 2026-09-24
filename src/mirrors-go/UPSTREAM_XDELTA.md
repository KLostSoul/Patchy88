# Mirrors용 Patchy88 — 공식 xdelta v3.2.0 재현 정보

공식 소스: https://github.com/jmacd/xdelta/tree/v3.2.0
라이선스: https://github.com/jmacd/xdelta/blob/v3.2.0/xdelta3/LICENSE (Apache 2.0)
공식 Windows x64: https://github.com/jmacd/xdelta/releases/download/v3.2.0/xdelta3-3.2.0-windows-x86_64.zip

| 파일 | SHA-256 | 빌드 출처 |
|---|---|---|
| 공식 x64 ZIP | `af8ef036cb077a48df080c9a8ac1be4a6e7511c32d11f8bec89b6803a9e52576` | jmacd/xdelta v3.2.0 릴리스 |
| `assets/xdelta3-x64.exe` | `53d90226615f217d3380c39892833311b4e24acd863e1ca01f14b5e772e2e6d0` | 위 ZIP에서 추출 |
| `assets/xdelta3.exe` | `232a8e8ac9fb47a54d0ca4d6acdb322e456a8c72cbfe3b212cf2f7e498760481` | 공식 v3.2.0 태그 소스, Windows Win32 빌드 |

[빌드 워크플로](../../.github/workflows/mirrors-upstream-xdelta.yml) · [x64/x86 성공 기록](https://github.com/KLostSoul/Patchy88/actions/runs/36051984391)

패처 실행파일은 호스트 아키텍처에 맞는 디코더를 자동 선택하고 매니페스트에 기록된 SHA-256을 검증합니다. 이전 Pachy98 동봉 실행파일은 배포 ZIP에서 제외했습니다. xdelta 바이너리 자체는 소스 Git 저장소에 넣지 않습니다.
