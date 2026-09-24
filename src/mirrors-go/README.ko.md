# Mirrors_Kor1.00 — PC-8801 xdelta 패처 소스

PC-8801 《Mirrors》의 **일본판 또는 영문판** CCD/IMG/SUB 파일 3개를 MD5 + SHA-256으로 검사하고, 각 판본 전용 xdelta 3개를 적용합니다.

- 한 판본만 검출되면 자동 선택
- **일본판과 영문판이 모두 검출되면 사용자에게 원본 선택 요청** (GUI: 예=일본판 / 아니오=영문판 / 취소)
- 선택한 판본의 xdelta만 적용하며 두 원본 세트는 수정하지 않음
- 검증된 한글판 결과 `Mirrors_Kor1.00.ccd`, `.img`, `.sub` 생성
- 같은 폴더에 `Mirrors_Kor1.00.cue`, `disk1main.d88`, `disk2game.d88` 추가
- CUE 내부 IMG 경로도 `Mirrors_Kor1.00.img`
- 기존 파일 충돌 방지, 결과 MD5 + SHA-256 검증, 실패 시 생성물 복구

## CLI (Windows 이외의 개발 환경)

```sh
go run . scan FOLDER
go run . apply FOLDER Japanese
go run . apply FOLDER English
```

두 판본이 있는 경우 CLI에서도 판본 선택이 필수입니다.

## 빌드

```sh
go test ./...
GOOS=windows GOARCH=amd64 go build -ldflags="-H windowsgui" -o Mirrors_Kor1.00-x64.exe
GOOS=windows GOARCH=386   go build -ldflags="-H windowsgui" -o Mirrors_Kor1.00-x86.exe
```

실제 배포 패키지의 `assets/` 디렉터리는 바이너리 자산이므로 소스 저장소에는 넣지 않습니다. 빌드한 실행파일 옆에 `assets/`가 필요합니다. 자산 이름과 해시는 `config/Mirrors_Kor1.00.json`을 참고하세요. 실제 원본 이미지에 xdelta를 적용한 최종 테스트는 일본판/영문판 원본 전체가 없어 수행하지 못했습니다.
