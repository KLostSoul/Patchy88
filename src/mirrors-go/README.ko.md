# Mirrors용 Patchy88 — Mirrors_Kor1.01 (PC-8801)

사용자용 설명과 원본/결과 전체 해시는 [Mirrors용 Patchy88 설명서](../../docs/MIRRORS.md)에 있습니다. **배포 ZIP에는 Go 소스를 포함하지 않습니다.**

Go 패처는 원본 CCD·IMG·SUB 전체 MD5·SHA-256으로 일본판과 영문판을 식별합니다. 양쪽이 모두 있으면 선택을 요구합니다. 두 판본의 최종 CCD·IMG·SUB는 동일한 한글판으로 검증합니다. 결과 파일명과 해시가 같으며, 기존의 정상 결과를 재검사할 수 있습니다.

## 공통 한글판 결과

| 파일 | 크기 | MD5 | SHA-256 |
|---|---:|---|---|
| CCD | 3,500 B | `80273154A2DAF2D107A282B353A86246` | `5512EFCEC20279753E1C6FEF0C5AAE37F2F3F05BB3AE3900B5DB0F688D4A67AC` |
| IMG | 551,779,200 B | `56E768F7CE3315A8172338CB10CE153E` | `FDCF60364815ADF0E85C2B796533276E2C72210F1024425C02757BDF88333B10` |
| SUB | 22,521,600 B | `E45923398B1D150F71E6056BD0974D15` | `C20C93FF3A2CF3A1E3A18B5F8046EC6F6198A916B417A01678B0B9D10EED9046` |

`core.go`: 판본 선택, 결과 판본 구분, 안전한 적용·정리. `verify_img.go`: CCD·IMG·SUB 해시 및 크기 검사. `core_test.go`: 판본별 출력·기존 결과·충돌·부분 재실행 테스트. `main_windows.go`: 사용자 폴더 및 판본 선택 GUI. 디코더 출처는 [UPSTREAM_XDELTA.md](UPSTREAM_XDELTA.md)를 참조하십시오.

```sh
cd src/mirrors-go
go test -count=1 ./...
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags='-H windowsgui -s -w' -o Mirrors_Kor1.01-x64.exe .
GOOS=windows GOARCH=386 CGO_ENABLED=0 go build -trimpath -ldflags='-H windowsgui -s -w' -o Mirrors_Kor1.01-x86.exe .
```

비Windows CLI: `go run . scan FOLDER`, `go run . apply FOLDER Japanese` 또는 `English`. 공식 xdelta3 바이너리와 패치·D88·CUE는 실행파일 옆의 `assets/`에서 읽으며 Git 소스에는 넣지 않습니다.

실제 영문판 원본 CCD·IMG·SUB를 이용한 종단 간 적용 및 최종 해시 검증을 완료했습니다.

## CCD/SUB 처리

일본판 CCD·SUB는 원본이 한글판 결과와 동일하므로 검증 후 복사합니다. 영문판 CCD·IMG·SUB는 각각의 xdelta를 적용합니다. 모든 출력은 동일한 한글판 해시와 크기로 검증합니다.
