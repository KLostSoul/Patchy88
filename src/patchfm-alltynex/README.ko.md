# PatchFM — Alltynex Kor v1.0

FM TOWNS용 《Alltynex》 한글판 ISO 생성 패처 소스입니다.

PatchFM은 사용자가 준비한 **원본 Alltynex ZIP 전체 파일**을 MD5·SHA-256·크기로 검증한 다음 xdelta를 적용하여 `Alltynex (Kor v1.0).iso`를 생성합니다.

**원본 게임 ZIP과 완성 ISO는 저장소 및 배포물에 포함하지 않습니다.** 소스 저장소에도 원본 ZIP이나 ISO를 넣지 않습니다.

## 검증값

원본:
- 크기: `574,208 bytes`
- MD5: `BCA57CBE9A2F23C922118F40FFFAE43A`
- SHA-256: `19B662F038E50E99535F60DBE6E07A77E269F574ED9250D46FC0655BB2AA6B3A`

결과 ISO:
- 파일명: `Alltynex (Kor v1.0).iso`
- 크기: `5,222,400 bytes`
- MD5: `D5A438B7C7AAC9E5B08BAEA6CB361F70`
- SHA-256: `114A9FFA82E32BCC8CACB8B0475810422056560260EAE5697CCF32B17BC91A54`

xdelta:
- SHA-256: `D7EF05131F9F2F2809E73E7E77ADF77F4FF19F89A2C59BCF6E604EC8F4CFC99B`

## 안전 적용

```text
원본 ZIP 선택
→ 전체 크기/MD5/SHA-256 검증
→ xdelta 및 xdelta3 SHA-256 검증
→ 임시 ISO 생성
→ 결과 ISO 전체 크기/MD5/SHA-256 검증
→ 검증 성공 시에만 Alltynex (Kor v1.0).iso 확정
```

기존 결과 ISO가 정상 해시라면 재생성하지 않습니다. 같은 이름의 파일이 있지만 검증값이 다르면 덮어쓰지 않습니다.

## 빌드

```sh
go test ./...
GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -trimpath -ldflags="-H windowsgui -s -w" -o PatchFM-Alltynex-x64.exe .
GOOS=windows GOARCH=386 CGO_ENABLED=0 go build -trimpath -ldflags="-H windowsgui -s -w" -o PatchFM-Alltynex-x86.exe .
```

실행 시 실행파일 옆 `assets/`에 다음 파일이 필요합니다.

```text
assets/
├─ Alltynex_Kor_v1.0.json
├─ xdelta3-x64.exe
├─ xdelta3.exe
└─ patches/
   └─ Alltynex_Kor_v1.0.xdelta
```

xdelta와 xdelta3 실행파일은 배포 자산이며 이 소스 디렉터리에는 저장하지 않습니다.

## FreeTOWNSOS / TSUGARU OS

생성되는 ISO의 부팅 환경은 CaptainYS(Soji Yamakawa)의 FreeTOWNSOS / TSUGARU OS를 사용합니다.

- Upstream: https://github.com/captainys/FreeTOWNSOS
- 라이선스: [FREETOWNSOS_LICENSE.md](FREETOWNSOS_LICENSE.md)

원 라이선스 전문을 그대로 포함했습니다.

## xdelta3

- [UPSTREAM_XDELTA.md](UPSTREAM_XDELTA.md)
- Patchy88 저장소의 Apache License 2.0 및 관련 고지를 따릅니다.
