# Patchy88 IPS 자산 검증표

이 문서는 Patchy88 개발 및 배포에 사용한 IPS 패치 파일의 식별값을 기록합니다.

> 원본 게임 D88, 원본 KANJI ROM, 패치 완료 게임 이미지는 이 저장소에 포함하지 않습니다.

## 몽환전사 바리스

| 파일 | 크기 | SHA-256 | Git blob SHA-1 |
|---|---:|---|---|
| `Valis_Korean_Disk_A_Patch_Ver_1.02.ips` | 29,748 | `7f14c7b5d6961e234f702aa3e6007944ad3ec8231af225f6296ef3491e1eff53` | `0f2006e22ce15dd8ca864f22b1687c18c904a298` |
| `VALIS_KANJI1_ROM_Patch_Ver_1.02.ips` | 22,759 | `1cb66ed56faf20a29cf0ee860805a14fc7d9132f825c22fa846c3bb81a70bc7c` | `8b96e38d78478a82295de2fd8aa1a7e15bc3715c` |

## 몽환전사 바리스 II

Patchy88 Valis2 v1.0.3부터 사용자가 제공한 **V1.01 IPS 세트**를 기준으로 사용합니다.

| 배포 내 파일 | 원본 제공 파일 | 크기 | SHA-256 | Git blob SHA-1 |
|---|---|---:|---|---|
| `Valis2_KOR_Disk_A.ips` | `Valis2_KOR_Disk_A_V1.01.ips` | 13,065 | `547754b043a6c4c2e70f62af22c045f73337822a7d9e49fd3b6825e9a4e4c47f` | `ba5c1f4aa4560d5676ff8c472f76bf5401d22320` |
| `Valis2_KOR_Disk_B.ips` | `Valis2_KOR_Disk_B_V1.01.ips` | 16,978 | `5151fe86ab6da0bc8d98671be4ea2f9eea248b98c144a79050b26cbacdd48722` | `d05c16b7afe3f1fc8ba5076ceb3ae01548a606eb` |
| `Valis2_KOR_Disk_C.ips` | `Valis2_KOR_Disk_C_V1.01.ips` | 317 | `b315445e42650e0da42edb149cfe367fa8dd2c3d6be49a66a7407c815bf218c6` | `f9d52782be768db1ea78da43dbed34e9e38f8375` |
| `Valis2_KOR_Disk_D.ips` | `Valis2_KOR_Disk_D_V1.01.ips` | 10,065 | `63481fb6055e03f4e5b1bfabd21c2f6db808086beb481f3f4f797c35fdfea436` | `683aa83483b77e27322a7b2595d3324b1a75d614` |
| `Valis2_KOR_Disk_E.ips` | `Valis2_KOR_Disk_E_V1.01.ips` | 4,145 | `59132ac1c21f45aa7f64443a1234e7a38faf2dfae012e0ff831fa28e982b324d` | `18392ec082baef06f73fcfec0dd1bfca8953bef5` |
| `Valis2_KOR_Disk_F.ips` | `Valis2_KOR_Disk_F_V1.01.ips` | 3,833 | `585506a1ffcf17d073a3e0957f53cd6f7c795565905c7d6acfa9cfbc0fc71518` | `2c52194abde6b8b281715990cdf58dddd52c166e` |
| `Valis2_KOR_Disk_G.ips` | `Valis2_KOR_Disk_G_V1.01.ips` | 3,333 | `30cde63b1a960f733d64e465305d08d80afddc0ac2118b4ef76e1076de966eb8` | `2a5e99998451923256b16f4ec38427dc2c57b7fe` |
| `Valis2_KOR_KANJI1.ips` | `Valis2_KOR_KANJI1_V1.01.ips` | 29,269 | `a8fe0820e0152e0164f064b1e6ea938c9117807ca103e45085ac620a0309d473` | `5b64a17c719ad767b9b6adbd638df2ca93b6b711` |

### V1.01 구조 검증

| 대상 | Records | RLE |
|---|---:|---:|
| Disk A | 432 | 0 |
| Disk B | 451 | 0 |
| Disk C | 19 | 0 |
| Disk D | 186 | 0 |
| Disk E | 73 | 0 |
| Disk F | 107 | 0 |
| Disk G | 294 | 0 |
| KANJI1 | 1,337 | 0 |

모든 IPS는 `PATCH` 헤더와 정상 `EOF`를 확인했고, 현재 기준 원본 범위를 벗어나는 레코드·중첩 레코드·원본과 동일한 no-op 레코드는 없음을 확인했습니다.

## Git blob SHA-1 계산법

```text
"blob " + 파일크기 + NUL + 파일바이트
```
