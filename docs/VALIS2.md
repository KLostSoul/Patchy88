# 몽환전사 바리스 II용 Patchy88

PC-8801 《몽환전사 바리스 II》 한글패치용 Patchy88 문서입니다.

## 현재 버전

- Windows x64/x86 독립 실행판: **Patchy88 Valis2 v1.0.0**
- 대상
  - Disk A~G D88
  - KANJI1 ROM

## 사용 방법

Windows판은 자신의 환경에 맞는 실행파일을 실행합니다.

```text
Patchy88-Valis2-x64.exe
Patchy88-Valis2-x86.exe
```

사용자는 **Disk A~G와 KANJI1 ROM이 들어 있는 폴더 하나만 지정**합니다.

Patchy88은 지정 폴더 바로 아래의 `.d88` / `.rom` 파일을 검사하고, 파일명이 아니라 **패치 대상 원본 데이터**를 기준으로 Disk A~G와 KANJI1을 자동 식별합니다.

예를 들어 파일명이 아래처럼 바뀌어 있어도 내용이 지원 기준과 일치하면 식별할 수 있습니다.

```text
one.d88
two.d88
three.d88
font.rom
```

## 일괄 패치 절차

바리스 II판은 한 파일씩 즉시 수정하지 않습니다.

1. 지정 폴더의 `.d88`, `.rom` 검사
2. Disk A~G / KANJI1 자동 식별
3. 8개가 모두 유일하게 식별되는지 확인
4. 누락, 중복, 부분패치, 비호환, IPS 변조 여부 검사
5. 하나라도 문제가 있으면 **전체 작업 중단**
6. 8개 임시 패치 결과 생성
7. 8개 결과 모두 사후검증
8. 전부 정상일 때만 실제 결과 확정

따라서 Disk A만 패치되고 Disk D에서 실패하는 식의 중간 완료 상태를 피하도록 설계했습니다.

## 출력 규칙

패치 결과는 원본과 같은 폴더에 `(K)`를 붙여 생성합니다.

```text
Mugen Senshi Valis II (Disk A).d88
→ Mugen Senshi Valis II (Disk A).d88.bak
→ Mugen Senshi Valis II (Disk A)(K).d88
```

KANJI1도 동일합니다.

```text
KANJI1.ROM
→ KANJI1.ROM.bak
→ KANJI1(K).ROM
```

기존 백업은 덮어쓰지 않습니다.

```text
game.d88.bak
game.d88.1.bak
game.d88.2.bak
```

## 검증 방식

Patchy88은 각 IPS가 실제로 건드리는 원본 영역의 적용 전/후 SHA-256을 기준으로 상태를 판정합니다.

- 모든 대상 영역이 적용 전 값과 일치 → `ORIGINAL`
- 모든 대상 영역이 적용 후 값과 일치 → `ALREADY_PATCHED`
- 적용 전/후 상태가 혼재 → `PARTIAL`
- 어느 쪽에도 맞지 않는 대상 영역 존재 → `INCOMPATIBLE`

`PARTIAL`, `INCOMPATIBLE` 상태가 하나라도 있으면 8개 전체 패치를 시작하지 않습니다.

## 실제 자동 식별 검증

테스트에서 Disk A~G와 KANJI1의 파일명을 임의로 바꿔도 8개를 내용으로 정확히 식별했습니다.

또한 Disk A에 IPS 레코드 일부만 적용해 `PARTIAL` 상태를 만든 경우에는 전체 패치를 거부하고 `.bak`이나 `(K)` 파일을 하나도 생성하지 않는 것을 확인했습니다.

## 기준 원본 및 패치 결과 해시

CRC64는 **CRC-64/ECMA-182** 기준입니다.

| 대상 | 원본 SHA-256 | IPS 패치 결과 SHA-256 |
|---|---|---|
| Disk A | `4E1E3F5A21B66BB9C57AB20F6522C95C8074DC9120F92D395EACDB673D9A3BD8` | `B96D108D5BA75023E67AD7C66BE370E5E1AAC8A6A73E8F8D5BD248F65C0ABF9A` |
| Disk B | `6537B681F8C24AA92BE608F99DC61EF8DA9AFA1B6569450DBA1DD5D7D00E3919` | `2368B443A1A536AFAB78F28C582DA268196A2C28F7305B2BAEDF11CCE2A161B4` |
| Disk C | `396F79DF59D3263B801126943D03F0C63B232A52A6256A0994254D83031FABDA` | `C6652CF326B9C5250C3C91F11BE9936BE430FD1182D1736CB1132AFC7CF36570` |
| Disk D | `BD453F4576705A224A6B1775D2E7D0D2D499EFD49075FF0B2B4B0B419549DD6D` | `5964758193523428947B0681E2A2B2F7B7ADD379FD86F3A5933022507AF16331` |
| Disk E | `905C9B0DBC402BB6134B63A67F3615AB0ADF79A191FF87CDF365903138CEFEDE` | `A872EB595CD8CD69377749EDC89786D5076939997DC7AC4A8927EC6E17962219` |
| Disk F | `8569EE720C0DF24293066FC3589C1A204EDC102EF57ABC3A6529C16EC8A4041B` | `E82286A4E9C3F8356DA7750BF9F80D504420022758666F7DA44E10A30C2750E4` |
| Disk G | `D933A387F6AFE847345482070CD1E23CA4158F6D49144A99EF9EBBDB864C56F7` | `51679430371D46474C21C6C7CC0F37FA0C05288E0EBE3719A4E18A8A837987E7` |
| KANJI1 | `F6C1C5022FE5935F6DFA3EB919E51441E75191270B639EDCB7938B3BCE41F6A3` | `30435619932143104C6E7D38E282372DDD011F4B7444D49E566AB995192491DC` |

CRC32, CRC64, SHA-1까지 포함한 전체 값은 별도 해시 문서에 기록되어 있습니다.

## 관련 자료

- [바리스 II 원본/패치 결과 전체 해시](Valis2_PC88_KOR_Hash_List.md)
- [Patchy88 공통 검증 구조](VALIDATION.md)
- [IPS 식별값](../patches/README.md)

## 소스

Windows판 Go 소스는 다음 위치에 있습니다.

```text
src/valis2-go/
```

빌드 예:

```bash
go test ./...
GOOS=windows GOARCH=amd64 go build -ldflags="-H windowsgui" -o Patchy88-Valis2-x64.exe
GOOS=windows GOARCH=386   go build -ldflags="-H windowsgui" -o Patchy88-Valis2-x86.exe
```

## 원본 데이터 배포 정책

Patchy88 저장소와 배포물에는 원본 게임 D88, 원본 KANJI ROM, 패치 완료 전체 게임 이미지를 포함하지 않습니다. 사용자는 적법하게 준비한 원본을 사용해야 합니다.
