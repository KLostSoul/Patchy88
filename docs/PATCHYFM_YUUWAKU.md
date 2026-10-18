# PatchyFM_Yuuwaku_v1.0 — FM TOWNS 《誘惑 / Yuuwaku》

이 패처는 사용자가 준비한 《誘惑》(Temptation) FM TOWNS 원본 CD 데이터 트랙 IMG/BIN을 입력으로 사용해 한글판 BIN/CUE를 생성합니다. **폴더 지정·자동 인식 방식**이며 원본 이미지 파일을 직접 선택하지 않습니다.

## 배포 파일

- `PatchyFM_Yuuwaku_v1.0.zip` — Windows x64/x86 GUI 실행파일, 정해진 xdelta, xdelta3 헬퍼, CUE 템플릿, 무결성 매니페스트 및 고지.
- `PatchyFM_Yuuwaku_v1.0_Source.zip` — Go 소스 및 테스트.
- 소스 경로: `src/patchyfm-yuuwaku/` (소스 트리에는 게임 원본, 완성 게임 이미지, 배포용 실행파일, xdelta를 넣지 않음).

## 사용 방법

1. 배포 ZIP을 통째로 압축 해제합니다.
2. 64비트 Windows에서 `PatchyFM-Yuuwaku-x64.exe` 또는 32비트 Windows에서 `PatchyFM-Yuuwaku-x86.exe`를 실행합니다.
3. **폴더 선택...**으로 원본 `.img` 또는 `.bin`이 있는 폴더를 선택합니다.
4. 하위 폴더까지 탐색한 후 해시가 일치하는 원본을 발견하면 **한글패치 적용**을 누릅니다.
5. 원본이 있는 디렉터리에 `Temptation-Korean-v1.0.bin`과 `Temptation-Korean-v1.0.cue`가 생성됩니다.

원본 파일명은 달라도 무방합니다. CUE/CCD/SUB를 직접 선택하거나 수정하지 않습니다. 동일한 원본 이미지가 두 개 이상 있으면 모호한 출력을 방지하고 종료합니다.

## 무결성

| 대상 | 크기 (bytes) | MD5 | SHA-256 |
|---|---:|---|---|
| 원본 `Yuuwaku.img` | 6,703,200 | `000b8a4e832d952fb29ea6b831ed8cd0` | `09095579184db0d9e1de7b17071550090782e134058b1d93aa432409dfc1557b` |
| xdelta | 330,715 | `aec822e8e826e754aefb9fb551d41f86` | `01c25b644e222bdb432fe484f6faed8a20714cce571ed2989d8bba79ead42e1e` |
| 한글판 BIN | 7,599,312 | `e07f8348b3a3b48f2e02613edeb204db` | `b407d4538a3bad2d2a63ef2ea6b036d6ff36e1c7b242c2494c5098dec9fcc9bf` |
| 한글판 CUE | 88 | `8a2606d84daa6495307ae6857977e49a` | `7787ee09485eabccc5699a3eb6e743d73d8796a8563bac5345720ee09ca4ae11` |

배포본의 CUE 템플릿·xdelta·xdelta3는 실행 전에 SHA-256으로 검증됩니다. 실제 적용 직전 원본 해시를 재검사하고, 임시 BIN 전체 해시를 검사한 뒤 파일명을 확정합니다. 기존 결과가 정상이면 재사용하고 다른 파일이면 덮어쓰지 않습니다.

## 빌드·검증 범위

Go 1.23 이상 기준:

```sh
cd src/patchyfm-yuuwaku
go test ./...
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -H windowsgui" -o PatchyFM-Yuuwaku-x64.exe .
GOOS=windows GOARCH=386 CGO_ENABLED=0 go build -trimpath -ldflags="-s -w -H windowsgui" -o PatchyFM-Yuuwaku-x86.exe .
```

실제 업로드 원본 IMG와 xdelta로 BIN을 복원한 결과, 업로드 한글판 BIN과 **전체 바이트 일치**했습니다. 폴더/하위 폴더 자동 인식, 확장자 및 이름 변경 인식, 기존 결과 보호, 중복 원본 거부 등 8건 통합 테스트에 합격했습니다. 단, 이 빌드 환경에서 Windows GUI를 직접 실행하거나 Tsugaru 부팅을 테스트하지 않았습니다.

## 라이선스와 권리

패처 Go 소스: Apache-2.0. xdelta3: [jmacd/xdelta](https://github.com/jmacd/xdelta), Apache-2.0. 원본 게임 이미지와 완성 BIN은 배포하지 않으며, 번역 패치의 권리는 패치 작성자에게 있습니다. FreeTOWNSOS는 이 패처에 포함하지 않습니다.
