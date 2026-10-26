# Independent binary verification

Uses Python 3 and standard-library `lzma` to decode the xdelta VCDIFF/LZMA sections and validate each window Adler-32.

```sh
python verification/verify_vcdiff_lzma.py ORIGINAL.img Dune2-Battle-for-Arrakis-Korean-v1.01.bin.xdelta TEMP_OUTPUT.img
```

Expected output length: 16,835,616 bytes. Expected output SHA-256: `5f3877d5d6be7b54b525ec91df19a46b01c4a6f5e83e3df6b683f05203053752`.

This independent verification script is only for auditing and tests. The Windows GUI uses the bundled official xdelta3 decoder.
