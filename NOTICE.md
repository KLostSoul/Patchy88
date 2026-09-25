# Patchy88 프로젝트 고지

Patchy88은 PC-8801용 패치 프로그램입니다. Valis용 IPS 패처와 Mirrors용 xdelta 패처는 엔진과 검증 방식이 다릅니다. Pachy98/romtools (46OkuMen)의 패치 배포 개념을 참고했습니다.

## Mirrors용 Patchy88 v1.01

일본판·영문판 각 3개씩 총 6개 xdelta와 CUE를 사용합니다. 두 판본의 원본과 **패치 후 CCD·IMG·SUB 해시가 서로 다르므로** 결과 기대값을 판본별로 분리했습니다. 일본판 IMG MD5 `56E768F7CE3315A8172338CB10CE153E`, 영문판 IMG MD5 `B737FADF8EAB6C4712F763E142E08DAB`를 기준으로 검사합니다. 기존 다른 판본의 결과가 있을 때는 파일 충돌을 거부하고 원본을 변경하지 않습니다. 결과 파일의 전체 MD5·SHA-256 및 크기를 필수로 검사합니다.

디코더: 공식 [jmacd/xdelta v3.2.0](https://github.com/jmacd/xdelta/tree/v3.2.0) Windows x64 및 동일 공식 소스의 Win32 빌드. 라이선스: Apache License 2.0. 배포 ZIP에는 라이선스·수정 고지문을 포함합니다. 원본 디스크와 패치 완료 전체 이미지는 배포하지 않습니다. 실제 원본을 통한 종단 간 출력 검증은 아직 수행하지 못했습니다.

[Mirrors 사용법과 검증 해시](docs/MIRRORS.md)
