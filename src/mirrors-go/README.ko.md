# Mirrors용 Patchy88 — Mirrors_Kor1.00 (PC-8801)

이 디렉터리는 PC-8801 《Mirrors》의 일본판·영문판을 한글판으로 변환하는 **Mirrors용 Patchy88**의 Go 소스입니다. **사용자용 설명은 [docs/MIRRORS.md](../../docs/MIRRORS.md)를 참조하십시오.** 배포 ZIP에는 이 소스를 넣지 않습니다.

## 처리 범위

- CCD/IMG/SUB 세 파일의 원본 MD5 + SHA-256 검증. 파일명으로 판본을 추측하지 않습니다.
- 일본판 또는 영문판 하나만 완전하게 발견되면 자동 선택. 둘 다 발견되면 **Windows GUI에서 사용자 선택** (예: 일본판, 아니오: 영문판, 취소: 중단).
- 선택한 판본 전용 xdelta 3개만 적용. 두 판본의 원본은 덮어쓰지 않습니다.
- 결과 MD5 + SHA-256 검증에 성공하면 같은 폴더에 `Mirrors_Kor1.00.ccd`, `Mirrors_Kor1.00.img`, `Mirrors_Kor1.00.sub`, `Mirrors_Kor1.00.cue`, `disk1main.d88`, `disk2game.d88`을 추가합니다.
- 이미 다른 데이터가 들어 있는 결과 파일은 덮어쓰지 않으며, 오류 시 이번 실행에서 생성한 임시 파일·결과 파일을 정리합니다.

## 소스 및 배포 자산 구분

`core.go` — 원본 판별, 선택 유지, xdelta 실행, 원본·결과 검증 및 안전한 결과 확정.

`main_windows.go` — Windows 폴더 선택 및 두 판본 동시 검출 시 선택 UI.

`main_cli.go` — Windows 이외의 개발 환경에서 검사·패치 실행.

`core_test.go` — 자동 판별, 충돌 검사, 자산 변조, 실패 정리, 두 판본 동시 검출 및 선택 테스트.

`config/Mirrors_Kor1.00.json` — 원본·패치·결과·추가 파일의 파일명과 해시를 기록한 기준 매니페스트.

`config/Mirrors_Kor1.00.cue` — CUE 소스. 내부에서 `Mirrors_Kor1.00.img`를 참조합니다.

실행 시 필요한 `assets/`(매니페스트, xdelta 6개, xdelta3.exe, CUE 및 D88 2개)는 사용자 배포 ZIP에 별도로 포함됩니다. 원본 CCD/IMG/SUB는 저장소 또는 배포 ZIP에 포함하지 않습니다.

## 빌드 및 테스트

```sh
cd src/mirrors-go
go test ./...
GOOS=windows GOARCH=amd64 go build -ldflags="-H windowsgui" -o Mirrors_Kor1.00-x64.exe
GOOS=windows GOARCH=386 go build -ldflags="-H windowsgui" -o Mirrors_Kor1.00-x86.exe
```

비Windows 개발 환경의 CLI:

```sh
go run . scan /path/to/folder
go run . apply /path/to/folder Japanese
go run . apply /path/to/folder English
```

둘 다 발견됐을 때 판본을 명시하지 않은 `apply`는 실행을 거부합니다. 실제 일본판·영문판 전체 원본을 사용한 최종 xdelta 패치 시험은 별도로 필요합니다.

## xdelta 공식 소스 및 라이선스

공식 소스: https://github.com/jmacd/xdelta

3.0.11 기반 Apache 2.0 재라이선스 브랜치: https://github.com/jmacd/xdelta/tree/release3_0_apl

라이선스 원문: https://github.com/jmacd/xdelta/blob/release3_0_apl/xdelta3/LICENSE

배포 ZIP의 `LICENSE`와 `THIRD_PARTY_LICENSE.txt`는 위 원문과 바이트 동일합니다. 동봉 `xdelta3.exe`는 Pachy98 v0.20.1의 파일과 일치하지만 실제 빌드에 사용된 소스 브랜치는 확인되지 않았습니다.

