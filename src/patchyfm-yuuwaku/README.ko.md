# PatchyFM_Yuuwaku_v1.01 — FM TOWNS 《誘惑 / Yuuwaku》

FM TOWNS판 《誘惑》(Yuuwaku / Temptation) 원본 CD의 데이터 트랙 IMG/BIN에 **수정된 v1.01 xdelta**를 적용하여 한글판 BIN/CUE를 만드는 Windows 패처입니다. 기존 v1.0의 **폴더 지정·하위 폴더 자동 검색** 방식을 유지합니다.

## 사용법

1. 배포 ZIP을 폴더 전체에 압축 해제합니다. EXE만 단독 실행하지 마십시오.
2. Windows 64비트: `PatchyFM-Yuuwaku-x64.exe`, 32비트: `PatchyFM-Yuuwaku-x86.exe`를 실행합니다.
3. **폴더 선택...**에서 원본 IMG/BIN 이미지가 있는 폴더를 지정합니다.
4. 선택 폴더와 하위 폴더에서 IMG/BIN을 검색하고 **크기·MD5·SHA-256**으로 원본을 자동 인식합니다. CUE/CCD/SUB를 직접 고를 필요가 없습니다.
5. **한글패치 적용**을 누르면 발견한 원본 IMG/BIN과 같은 폴더에 **`Temptation-Korean-v1.01.bin`**, **`Temptation-Korean-v1.01.cue`**가 생성됩니다.
6. Tsugaru 등에서는 생성된 CUE를 CD-ROM으로 선택합니다. 별도의 게임 User Disk가 필요할 수 있으며 이 배포물에 포함하지 않습니다.

원본 IMG/BIN/CCD/CUE/SUB는 수정하지 않습니다. 지정 폴더에 복수의 동일 원본 이미지가 있으면 작업을 중단합니다. 이미 검증된 결과 BIN/CUE가 있다면 재사용하고, 같은 이름의 **내용이 다른 파일은 덮어쓰지 않습니다.** 기존 v1.0 결과는 새 파일명과 충돌하지 않습니다.

## 파일 무결성

| 대상 | 크기 (바이트) | MD5 | SHA-256 |
|---|---:|---|---|
| 원본 Yuuwaku.img | 6,703,200 | `000b8a4e832d952fb29ea6b831ed8cd0` | `09095579184db0d9e1de7b17071550090782e134058b1d93aa432409dfc1557b` |
| 수정 xdelta `Temptation-Korean-v1.01.xdelta` | 330,728 | `1c49f1bed50839abd768e8c74b684f43` | `c936529c182cd2854644b5a839406937628437325ef161f9af6c77cf36919e09` |
| 한글판 `Temptation-Korean-v1.01.bin` | 7,599,312 | `c1f300af9d099e0c61234762c7fa96e6` | `1b456c4ff7d78c7f144986c6714f07a6fd62ecdd5b5c29c403656850653447c2` |
| `Temptation-Korean-v1.01.cue` | 89 | `470955c61b88294f7da020ec7903d0e1` | `bf7336f007be2f8a0c1b3cb2d1a2acda8677200fab41ee501271d58e914c5998` |

새 패치는 업로드된 **`Temptation-Korean-v1.0(1).xdelta`**의 바이트를 수정하지 않고 배포용으로 `Temptation-Korean-v1.01.xdelta`이라고 명명한 것입니다. 패치의 원본 파일명에 v1.0이 있었지만, **패처의 버전 및 새 결과물을 v1.01로 구분**합니다.

## v1.0 대비

- 제공받은 새 xdelta: **330,728바이트** (기존 330,715바이트).
- 결과 BIN 크기: **7,599,312바이트**, 동일.
- 기존 v1.0 결과와 새 결과 BIN: **41,809바이트가 변경됨**.
- 새 BIN의 MD5/SHA-256을 재측정하고 매니페스트에 반영함.
- 결과물을 `v1.01.bin/.cue`로 구분하여 기존 v1.0 결과와 충돌하지 않음.
- 원본 자동 탐색, 안전한 임시 파일 생성, 입력·출력 무결성 검증, 중복 검색 방지 등 기존 로직 유지.

## 검증 범위

수정 xdelta를 직접 해독한 VCDIFF 결과는 7,599,312바이트이며, VCDIFF 윈도우의 **Adler-32 무결성 검사 통과** 및 최종 MD5/SHA-256을 확인했습니다. 같은 자체 디코더로 기존 v1.0 xdelta를 해독하면 업로드된 v1.0 한글판 BIN과 완전히 일치함을 교차 확인했습니다. Go 단위 테스트, Windows x64/x86 빌드 및 별도 디코더를 사용하는 9가지 Linux 통합 검사(정상 생성·하위 폴더·재실행·변조/오류 거부 등)를 수행했습니다.

**이 Linux 환경에서 Windows GUI와 xdelta3.exe의 실제 실행 및 Tsugaru 부팅은 검증하지 않았습니다.** 따라서 검증된 범위와 실기 동작을 혼동하지 않습니다.

## 배포 구성

```text
PatchyFM_Yuuwaku_v1.01/
├─ PatchyFM-Yuuwaku-x64.exe
├─ PatchyFM-Yuuwaku-x86.exe
├─ README.ko.md
├─ LICENSE
├─ THIRD_PARTY_LICENSE.txt
├─ NOTICE_MODIFICATIONS.txt
├─ UPSTREAM_XDELTA.txt
├─ SHA256SUMS.txt
└─ assets/
   ├─ Yuuwaku_Kor_v1.01.json
   ├─ Temptation-Korean-v1.01.cue
   ├─ xdelta3.exe
   ├─ xdelta3-x64.exe
   └─ patches/
      └─ Temptation-Korean-v1.01.xdelta
```

원본 게임 이미지 및 완성 한글판 이미지/BIN은 배포하지 않습니다. 패치 툴 소스 코드는 Apache-2.0을 사용하며, 게임 원작과 번역 패치의 권리는 별개입니다. xdelta3는 jmacd/xdelta v3.2.0 (Apache-2.0)이며 FreeTOWNSOS는 포함하지 않습니다.
