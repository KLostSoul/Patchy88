# Mirrors_Kor1.00 — xdelta 패처 소스

PC Engine CD 《Mirrors》 일본판 또는 영문판의 CCD/IMG/SUB 3개 파일을 MD5 + SHA-256으로 식별하고, 각 판본 전용 xdelta 패치를 적용하는 Patchy88 계열 패처 소스입니다.

## 동작

- 일본판 / 영문판 CCD·IMG·SUB 자동 식별
- 판본별 xdelta 3개 적용
- 결과 CCD·IMG·SUB의 MD5 + SHA-256 검증
- 성공 후 같은 폴더에 `Kor.cue`, `disk1main.d88`, `disk2game.d88` 추가
- 원본 파일은 수정하지 않음
- 기존 동일 이름 파일은 덮어쓰지 않음
- 실패 시 이번 실행에서 생성한 임시/결과 파일 정리

## 배포물과 소스 분리

사용자 배포 ZIP에는 Go 소스를 넣지 않습니다. 이 디렉터리가 패처 소스의 기준 위치입니다.

실제 배포에는 별도로 다음 자산이 필요합니다.

- `Mirrors_Kor1.00.json`
- 일본판/영문판 xdelta 패치 6개
- `xdelta3.exe`
- `Kor.cue`
- `disk1main.d88`
- `disk2game.d88`

게임 원본 CCD/IMG/SUB는 배포하지 않습니다.

## 빌드

```bash
go test ./...
GOOS=windows GOARCH=amd64 go build -ldflags="-H windowsgui" -o Mirrors_Kor1.00-x64.exe
GOOS=windows GOARCH=386   go build -ldflags="-H windowsgui" -o Mirrors_Kor1.00-x86.exe
```

Windows 이외에서는 CLI 빌드가 사용됩니다.

```bash
go run . scan /path/to/source-folder
go run . apply /path/to/source-folder
```

## 검증 한계

현재 소스의 자동 식별, 자산 변조 검출, 충돌 방지, 실패 롤백 동작은 테스트되어 있습니다. 일본판/영문판 실제 원본 CD 이미지 전체는 저장소에 포함하지 않으므로 실제 원본에 대한 xdelta 최종 적용 검증은 별도로 수행해야 합니다.
