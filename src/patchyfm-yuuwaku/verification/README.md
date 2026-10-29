# Yuuwaku v1.02 independent VCDIFF verification

An independent Python 3 decoder is provided for source and release auditing.
It is **not required** to run the Windows patcher, and it does not include any game images.

```sh
python3 verification/verify_vcdiff_lzma.py /path/to/Yuuwaku.img /path/to/Temptation-Korean-v1.02.xdelta /tmp/Temptation-Korean-v1.02.bin
```

Expected results:
- 1 VCDIFF window, Adler-32 `0x45ed3812`.
- output size `7,599,312` bytes, MD5 `6a83336e6b2b85f63567712e3144895d`.
- output SHA-256 `746c23498bbac452c314ca4f07962dc65cd6c0dfcb2dd90441ee7a2de904b5f5`.

The verification script uses in-memory intermediate data. It is a developer-only auditor; the Windows patcher uses bounded streaming hash validation and bounded external-decoder diagnostics. It does **not** establish that the FM TOWNS game runs without memory errors.
