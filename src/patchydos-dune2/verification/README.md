# Independent binary verification

Uses Python 3 and standard-library `lzma` to decode the xdelta VCDIFF/LZMA sections and validate each window Adler-32.

```sh
python verification/verify_vcdiff_lzma.py ORIGINAL.img Dune2-Battle-for-Arrakis-Korean-v1.01.bin.xdelta TEMP_OUTPUT.img
```

Expected output length: 16,859,136 bytes. Expected output SHA-256: `d8c4e5025ad35886cd29b164ab69134eb4fce502b507ac3af15fc7e647d263e5`.

This independent verification script is only for auditing and tests. The Windows GUI uses the bundled official xdelta3 decoder.

Updated 2026-10-27: VCDIFF header indicator 0x01 (application-header flag absent). The original version v1.01 output was 16,835,616 bytes; the new updated patch produces 16,859,136 bytes. Use the new patch with the updated manifest rather than mixing editions.
