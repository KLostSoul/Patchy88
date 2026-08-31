# Patchy88 IPS 자산 검증표

이 문서는 Patchy88 개발 및 배포에 사용한 IPS 패치 파일의 식별값을 기록합니다.

> 원본 게임 D88, 원본 KANJI ROM, 패치 완료 게임 이미지는 이 저장소에 포함하지 않습니다.

## 몽환전사 바리스

| 파일 | 크기 | SHA-256 | Git blob SHA-1 |
|---|---:|---|---|
| `Valis_Korean_Disk_A_Patch_Ver_1.02.ips` | 29,748 | `7f14c7b5d6961e234f702aa3e6007944ad3ec8231af225f6296ef3491e1eff53` | `0f2006e22ce15dd8ca864f22b1687c18c904a298` |
| `VALIS_KANJI1_ROM_Patch_Ver_1.02.ips` | 22,759 | `1cb66ed56faf20a29cf0ee860805a14fc7d9132f825c22fa846c3bb81a70bc7c` | `8b96e38d78478a82295de2fd8aa1a7e15bc3715c` |

## 몽환전사 바리스 II

v1.0.1에서 Disk A/B IPS가 최종 확인 완성본 기준으로 교체되었습니다. C~G와 KANJI1은 변경하지 않았습니다.

| 파일 | 크기 | SHA-256 | Git blob SHA-1 |
|---|---:|---|---|
| `Valis2_KOR_Disk_A.ips` | 13,065 | `6db84cbe1a0ff3d7197ae0df143ddd216c892ed0ad3a21732b1dc972043bad31` | `0b662f8e57192fd91064b59186acd3c551a3df3d` |
| `Valis2_KOR_Disk_B.ips` | 16,974 | `f8b855b31757d038786ecfddf3824d52baeff2fc98e2435bc8b6fe6372c33adb` | `f44d50ae8aa6bbb774b89f34c6b852692fe4cd8b` |
| `Valis2_KOR_Disk_C.ips` | 317 | `4ff788902a4580e05d684eb589f194735c9a6db8cc8921779f8aeb38c1f350a8` | `a08b8eb1c325fcda247bd3eb064395f77882ea2a` |
| `Valis2_KOR_Disk_D.ips` | 10,116 | `51f6faac4e74aa7cd23e18775319381af4ae10567dc4bac92554aea7cef36ada` | `816d5480e4339140029868fb2f276fbd7f5a8505` |
| `Valis2_KOR_Disk_E.ips` | 4,164 | `aa43e21df6999d7c284e913a4497697e510e095b62bc57884f2e516e130ac54a` | `81ac73f416e82707800b2d33482a33e938fa2c83` |
| `Valis2_KOR_Disk_F.ips` | 3,805 | `8e9dd7b075e561451c65cd1ca6400d20158e9b61019508b0d32d31a7f18858c7` | `b2c72ee8f5ae77ba0cc4c0d82ec4cc2ec87ef2cd` |
| `Valis2_KOR_Disk_G.ips` | 3,333 | `30cde63b1a960f733d64e465305d08d80afddc0ac2118b4ef76e1076de966eb8` | `2a5e99998451923256b16f4ec38427dc2c57b7fe` |
| `Valis2_KOR_KANJI1.ips` | 21,488 | `c1c3978b19429251e7c86114ed22128f70c567235ddd97467569886bc8de7067` | `d609a19ed7ad9489793f28c965e64a4862b0950d` |

### Disk A/B 교체 검증

```text
Disk A
records: 434
RLE: 0
패치 결과 SHA-256:
4a2cd146c410d01c207670778e47ac6a670c00caf93ef7c52db412b26832040d

Disk B
records: 450
RLE: 0
패치 결과 SHA-256:
5ab755cfe2aeaef68e4cbe4582c979fe6449482d2f09ab72d0e44918cbae1d6b
```

두 결과는 최종 확인된 한국어 완성 D88과 각각 바이트 단위로 동일합니다.

## Git blob SHA-1 계산법

Git blob 식별자는 다음 데이터의 SHA-1입니다.

```text
"blob " + 파일크기 + NUL + 파일바이트
```

따라서 저장소에 바이너리 패치를 추가하거나 교체할 때 위 표의 Git blob SHA-1까지 비교하면 Git 전송 과정에서 패치 바이트가 변하지 않았는지 확인할 수 있습니다.
