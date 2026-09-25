package main

import "fmt"

// Full-image MD5 and SHA-256 are mandatory for every Korean output.
// Size is separately required. Adler-32 fallback is deliberately removed.
func verifyTarget(_ string, got fileHashes, target TargetDef) (string, error) {
	if len(target.MD5) != 32 || len(target.SHA256) != 64 || target.Size <= 0 {
		return "", fmt.Errorf("결과 파일 검증 기준값이 불완전합니다")
	}
	if got.Size != target.Size {
		return "", fmt.Errorf("결과 파일 크기 불일치: got=%d expected=%d", got.Size, target.Size)
	}
	if !hashEqual(got, target.Hashes) {
		return "", fmt.Errorf("MD5/SHA-256 불일치: MD5=%s SHA-256=%s", got.MD5, got.SHA256)
	}
	return "MD5/SHA-256 및 크기", nil
}
