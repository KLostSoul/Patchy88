# Mirrors용 Patchy88 — Mirrors_Kor1.01 (PC-8801)

사용자용 설명과 원본/결과 전체 해시는 [Mirrors용 Patchy88 설명서](../../docs/MIRRORS.md)에 있습니다. **배포 ZIP에는 Go 소스를 포함하지 않습니다.**

Go 패처는 원본 CCD·IMG·SUB 전체 MD5·SHA-256으로 일본판과 영문판을 식별합니다. 양쪽이 모두 있으면 선택을 요구합니다. `config/Mirrors_Kor1.01.json`의 `source.Japanese`/`source.English`, `target.Japanese`/`target.English`에 기록된 **별도의 결과 해시**를 이용합니다. 결과 CCD·IMG·SUB 모두 전체 MD5·SHA-256과 크기를 검사하며 이전의 Adler-32 대체 검증 경로는 없습니다. 결과 파일명은 양쪽이 같으므로 다른 판본 결과가 이미 있는 경우에는 덮어쓰지 않습니다.

## 일본판 패치 결과

| 파일 | 크기 | MD5 | SHA-256 |
|---|---:|---|---|
| CCD | 3,500 B | `80273154A2DAF2D107A282B353A86246` | `5512EFCEC20279753E1C6FEF0C5AAE37F2F3F05BB3AE3900B5DB0F688D4A67AC` |
| IMG | 551,779,200 B | `56E768F7CE3315A8172338CB10CE153E` | `FDCF60364815ADF0E85C2B796533276E2C72210F1024425C02757BDF88333B10` |
| SUB | 22,521,600 B | `E45923398B1D150F71E6056BD0974D15` | `C20C93FF3A2CF3A1E3A18B5F8046EC6F6198A916B417A01678B0B9D10EED9046` |

## 영문판 패치 결과 (제공받은 CloneCD 측정값)

| 파일 | 크기 | MD5 | SHA-256 |
|---|---:|---|---|
| CCD | 3,500 B | `35C733769D60277FCCE522E121AF82AE` | `2DAEAAF64FD4C206CC28C2438A5FA480272DB0888CC19106380D13950DE1C102` |
| IMG | 551,779,200 B | `B737FADF8EAB6C4712F763E142E08DAB` | `6AB276B6B3A79B64DD8F513F528750EC19ED7999737FEEB5855CC0CFEFF27AAE` |
| SUB | 22,521,600 B | `F6D739F1B66082F7F06F403D30CF35EB` | `3C05EFEC20E9BE2CDA81500178FF1F8A7657908FD699D2070827E58D700122D2` |

`core.go`: 판본 선택, 결과 판본 구분, 안전한 적용·정리. `verify_img.go`: CCD·IMG·SUB 해시 및 크기 검사. `core_test.go`: 판본별 출력·기존 결과·충돌·부분 재실행 테스트. `main_windows.go`: 사용자 폴더 및 판본 선택 GUI. 디코더 출처는 [UPSTREAM_XDELTA.md](UPSTREAM_XDELTA.md)를 참조하십시오.

```sh
cd src/mirrors-go
go test -count=1 ./...
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags='-H windowsgui -s -w' -o Mirrors_Kor1.01-x64.exe .
GOOS=windows GOARCH=386 CGO_ENABLED=0 go build -trimpath -ldflags='-H windowsgui -s -w' -o Mirrors_Kor1.01-x86.exe .
```

비Windows CLI: `go run . scan FOLDER`, `go run . apply FOLDER Japanese` 또는 `English`. 공식 xdelta3 바이너리와 패치·D88·CUE는 실행파일 옆의 `assets/`에서 읽으며 Git 소스에는 넣지 않습니다.

로컬 테스트 20개와 Windows x64·x86 빌드 완료. 배포 ZIP SHA-256 `4561080e082bc9006da8da789af1b79545c8a72c6f96a4e4ee8f872e4cc60646`. 실원본 종단 간 패치 시험은 원본 부재로 미실시입니다.
