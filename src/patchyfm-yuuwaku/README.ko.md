# PatchyFM — 誘惑 / Yuuwaku Kor v1.0

FM TOWNS판 《誘惑》(`Yuuwaku`, 영문명 `Temptation`)의 원본 CD **데이터 트랙**에 xdelta를 적용하여 한글판 **BIN/CUE**를 생성하는 전용 Windows 패처입니다.

**원본 게임 이미지와 완성 한글판 BIN은 배포하지 않습니다.** 사용자가 원본 이미지를 준비해야 하며, 패처는 원본 파일을 수정하지 않습니다.

## 사용 방법

1. 이 ZIP의 **폴더 전체를 압축 해제**합니다. EXE만 따로 옮기면 동작하지 않습니다.
2. Windows 64비트는 `PatchyFM-Yuuwaku-x64.exe`, 32비트는 `PatchyFM-Yuuwaku-x86.exe`를 실행합니다.
3. **`폴더 선택...`**을 눌러 원본 이미지가 들어 있는 폴더를 지정합니다. 개별 IMG/CUE/CCD 파일은 선택할 필요가 없습니다.
4. 패처가 지정 폴더와 하위 폴더에서 **`.img` 및 `.bin` 파일을 자동 검색**하고 원본 이미지의 크기·MD5·SHA-256이 모두 맞는지 검증합니다. 검증되면 `한글패치 적용`을 누릅니다.
5. 자동으로 검색된 원본 IMG/BIN과 같은 폴더에 `Temptation-Korean-v1.0.bin` 및 `Temptation-Korean-v1.0.cue`가 생성됩니다.
6. Tsugaru 등에서는 **생성된 CUE**를 CD로 지정합니다. 실제 게임 부팅에는 별도의 User Disk가 필요할 수 있으며 이 패처에는 포함하지 않습니다.

자동 인식하는 대상은 **2352바이트/섹터인 원본 CD 데이터 트랙**의 `.img` 또는 `.bin` 파일입니다. `Yuuwaku.img`를 비롯해 파일명이 달라도 **크기·MD5·SHA-256이 동일하면 인식**합니다. CUE/CCD/SUB는 원본 폴더에 함께 두어도 되지만 직접 선택하거나 패치하지 않습니다. 원본 SUB 데이터는 한글판 CD 생성에 필요하지 않습니다.

- 지정한 폴더의 **하위 폴더까지 검색**하므로 이미지가 별도 폴더 안에 있어도 찾습니다.
- 동일한 원본 데이터의 IMG/BIN이 **2개 이상** 발견되면 의도하지 않은 위치에 결과를 만들지 않도록 중단합니다. 이때는 원본 1개만 있는 하위 폴더를 선택하세요.
- 결과 BIN과 CUE는 **검증된 원본 이미지가 위치한 폴더**에 저장됩니다. 지정한 상위 폴더가 아닐 수도 있습니다.

## 검증값

| 파일 | 크기 (bytes) | MD5 | SHA-256 |
|---|---:|---|---|
| 원본 `Yuuwaku.img` | 6,703,200 | `000b8a4e832d952fb29ea6b831ed8cd0` | `09095579184db0d9e1de7b17071550090782e134058b1d93aa432409dfc1557b` |
| 패치 `Temptation-Korean-v1.0.xdelta` | 330,715 | `aec822e8e826e754aefb9fb551d41f86` | `01c25b644e222bdb432fe484f6faed8a20714cce571ed2989d8bba79ead42e1e` |
| 한글판 `Temptation-Korean-v1.0.bin` | 7,599,312 | `e07f8348b3a3b48f2e02613edeb204db` | `b407d4538a3bad2d2a63ef2ea6b036d6ff36e1c7b242c2494c5098dec9fcc9bf` |
| 한글판 `Temptation-Korean-v1.0.cue` | 88 | `8a2606d84daa6495307ae6857977e49a` | `7787ee09485eabccc5699a3eb6e743d73d8796a8563bac5345720ee09ca4ae11` |

## 동작 및 보호 정책

- 실행 시 매니페스트, xdelta, 결과 CUE 템플릿, xdelta3 헬퍼의 무결성을 확인합니다.
- 원본 폴더에서 IMG/BIN을 자동 검색하며, 파일명과 무관하게 원본 전체 크기·MD5·SHA-256을 검증합니다. 실제 패치 직전에도 다시 검증합니다.
- xdelta는 **원본 IMG/BIN 데이터 트랙 → 한글판 BIN**에 적용됩니다. CCD/SUB 파일에 패치하지 않습니다.
- 생성 BIN을 임시 파일로 만들어 크기·MD5·SHA-256을 모두 검증한 후 최종 이름으로 확정합니다.
- CUE는 검증된 템플릿에서 생성하며 BIN 파일 이름을 정확히 참조합니다.
- 기존 결과 BIN/CUE가 정상이면 재사용하고, 같은 이름의 **다른 파일이 있으면 덮어쓰지 않고 중단**합니다.
- 원본 IMG/BIN/CCD/CUE/SUB를 수정하거나 삭제하지 않습니다. 임시 파일은 실패 시 제거합니다.

## 검증 범위

업로드된 `Yuuwaku.img` + `Temptation-Korean-v1.0.xdelta`로 별도 VCDIFF 디코더를 통해 **대상 BIN 전체 바이트가 업로드된 한글판과 동일한 것**을 확인했습니다. VCDIFF 윈도우의 Adler-32 및 최종 MD5/SHA-256도 일치했습니다. 패처 소스 단위 테스트 및 **폴더/하위 폴더 자동 탐색, 원본 파일명 변경, 정확한 BIN/CUE 생성, 재실행, 원본 오류, 중복 원본, 기존 결과 보존** 등 8개 통합 검증이 통과했습니다.

**Windows GUI는 컴파일 확인까지 완료했으며, 이 환경에서 Windows 실행화면과 실제 xdelta3.exe 실행을 직접 테스트하지는 못했습니다.** 게임의 부팅·진행 여부도 여기서는 확인하지 않았습니다.

## 배포 구성

- `PatchyFM-Yuuwaku-x64.exe`, `PatchyFM-Yuuwaku-x86.exe`
- `assets/patches/Temptation-Korean-v1.0.xdelta`
- `assets/Temptation-Korean-v1.0.cue` (게임 데이터가 없는 CUE 텍스트)
- `assets/Yuuwaku_Kor_v1.0.json` (검증값)
- `assets/xdelta3-x64.exe`, `assets/xdelta3.exe` (xdelta3 헬퍼)
- 라이선스/고지/검증 체크섬 파일

## 저작권 및 도구

PatchyFM 소프트웨어의 라이선스는 `LICENSE`를 참조하십시오. 이 라이선스가 게임, 번역 데이터, 패치 콘텐츠의 권리까지 부여한다는 의미는 아닙니다. 원본 이미지/완성 이미지의 저작권은 각 권리자에게 있습니다.

xdelta3는 jmacd의 xdelta 프로젝트 v3.2.0입니다. `UPSTREAM_XDELTA.txt`, `THIRD_PARTY_LICENSE.txt` 참조. **FreeTOWNSOS 파일은 이 패처 배포물에 포함하지 않습니다.**
