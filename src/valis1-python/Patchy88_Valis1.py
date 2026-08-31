#!/usr/bin/env python3
# SPDX-License-Identifier: Apache-2.0
# Patchy88 - Valis 1 PC-8801 Python edition.
from __future__ import annotations

import argparse
import hashlib
import json
import os
import shutil
import struct
import sys
import tempfile
from dataclasses import dataclass
from pathlib import Path

APP_NAME = "Patchy88 - Valis 1 PC-88"
APP_VERSION = "1.0.7"
MANIFEST_NAME = "Patchy88-Valis1.json"


class PatcherError(Exception):
    pass


@dataclass(frozen=True)
class IPSRecord:
    offset: int
    data: bytes


def app_dir() -> Path:
    if getattr(sys, "frozen", False) and hasattr(sys, "_MEIPASS"):
        return Path(sys._MEIPASS)
    return Path(__file__).resolve().parent


def sha256(data: bytes) -> str:
    return hashlib.sha256(data).hexdigest()


def parse_ips(blob: bytes) -> list[IPSRecord]:
    if not blob.startswith(b"PATCH"):
        raise PatcherError("IPS 헤더(PATCH)가 없습니다.")
    cursor = 5
    records: list[IPSRecord] = []
    while True:
        if cursor + 3 > len(blob):
            raise PatcherError("IPS가 EOF 전에 잘렸습니다.")
        if blob[cursor:cursor + 3] == b"EOF":
            cursor += 3
            if cursor != len(blob):
                raise PatcherError("이 버전은 IPS EOF 뒤의 확장 데이터를 허용하지 않습니다.")
            if not records:
                raise PatcherError("IPS에 패치 레코드가 없습니다.")
            return records

        if cursor + 5 > len(blob):
            raise PatcherError("IPS 레코드 헤더가 잘렸습니다.")
        offset = int.from_bytes(blob[cursor:cursor + 3], "big")
        size = int.from_bytes(blob[cursor + 3:cursor + 5], "big")
        cursor += 5
        if size:
            end = cursor + size
            if end > len(blob):
                raise PatcherError(f"IPS 레코드 0x{offset:06X} 데이터가 잘렸습니다.")
            data = blob[cursor:end]
            cursor = end
        else:
            if cursor + 3 > len(blob):
                raise PatcherError(f"IPS RLE 레코드 0x{offset:06X}가 잘렸습니다.")
            repeat = int.from_bytes(blob[cursor:cursor + 2], "big")
            value = blob[cursor + 2]
            cursor += 3
            if repeat == 0:
                raise PatcherError(f"IPS RLE 레코드 0x{offset:06X} 길이가 0입니다.")
            data = bytes([value]) * repeat
        records.append(IPSRecord(offset, data))


def load_manifest() -> dict:
    path = app_dir() / MANIFEST_NAME
    try:
        manifest = json.loads(path.read_text(encoding="utf-8"))
    except OSError as exc:
        raise PatcherError(f"매니페스트를 열 수 없습니다: {path}\n{exc}") from exc
    except json.JSONDecodeError as exc:
        raise PatcherError(f"매니페스트 JSON이 손상되었습니다: {exc}") from exc
    if manifest.get("schema") != "pachy98.pc88.valis1.ips_manifest":
        raise PatcherError("지원하지 않는 매니페스트 형식입니다.")
    return manifest


def target_by_id(manifest: dict, target_id: str) -> dict:
    for target in manifest.get("targets", []):
        if target.get("id") == target_id:
            return target
    raise PatcherError(f"매니페스트에 대상이 없습니다: {target_id}")


def load_and_check_patch(target: dict) -> list[IPSRecord]:
    patch_path = app_dir() / target["patch"]
    try:
        blob = patch_path.read_bytes()
    except OSError as exc:
        raise PatcherError(f"IPS를 열 수 없습니다: {patch_path}\n{exc}") from exc
    actual_sha = sha256(blob)
    if actual_sha != target["patch_sha256"]:
        raise PatcherError(
            "IPS SHA-256이 매니페스트와 다릅니다.\n"
            f"파일: {patch_path.name}\n현재: {actual_sha}\n기준: {target['patch_sha256']}"
        )
    records = parse_ips(blob)
    expected = target["records"]
    if len(records) != len(expected):
        raise PatcherError("IPS 레코드 수가 매니페스트와 다릅니다.")
    for record, meta in zip(records, expected):
        if record.offset != meta["offset"] or len(record.data) != meta["length"]:
            raise PatcherError(f"IPS 레코드 구조가 매니페스트와 다릅니다: #{meta['index']}")
        if sha256(record.data) != meta["after_sha256"]:
            raise PatcherError(f"IPS 레코드 데이터가 매니페스트와 다릅니다: #{meta['index']}")
    return records


def d88_structure_ok(data: bytes) -> tuple[bool, str]:
    if len(data) < 0x2B0:
        return False, "D88 헤더보다 파일이 짧습니다."
    declared_size = struct.unpack_from("<I", data, 0x1C)[0]
    if declared_size != len(data):
        return False, f"D88 헤더 파일 크기({declared_size})와 실제 크기({len(data)})가 다릅니다."
    pointers = [struct.unpack_from("<I", data, 0x20 + i * 4)[0] for i in range(164)]
    nonzero = [(i, p) for i, p in enumerate(pointers) if p]
    if [i for i, _ in nonzero] != list(range(80)):
        return False, "Valis Disk A 기준 80개 트랙 포인터 구조가 아닙니다."
    try:
        for pos, (track_index, pointer) in enumerate(nonzero):
            limit = nonzero[pos + 1][1] if pos + 1 < len(nonzero) else len(data)
            if pointer < 0x2B0 or pointer + 16 > limit or limit > len(data):
                return False, f"트랙 {track_index} 경계가 비정상입니다."
            sector_count = struct.unpack_from("<H", data, pointer + 4)[0]
            expected_count = 16 if track_index < 2 else 5
            expected_sector_size = 256 if track_index < 2 else 1024
            if sector_count != expected_count:
                return False, f"트랙 {track_index} 섹터 수가 {sector_count}입니다."
            cursor = pointer
            for _ in range(sector_count):
                if cursor + 16 > limit:
                    return False, f"트랙 {track_index} 섹터 헤더가 잘렸습니다."
                declared_count = struct.unpack_from("<H", data, cursor + 4)[0]
                sector_size = struct.unpack_from("<H", data, cursor + 14)[0]
                if declared_count != sector_count or sector_size != expected_sector_size:
                    return False, f"트랙 {track_index} 섹터 메타데이터가 기준과 다릅니다."
                cursor += 16 + sector_size
                if cursor > limit:
                    return False, f"트랙 {track_index} 섹터가 트랙 경계를 넘습니다."
            if cursor != limit:
                return False, f"트랙 {track_index} 끝에 해석되지 않은 바이트가 있습니다."
    except (struct.error, IndexError) as exc:
        return False, f"D88 구조 해석 오류: {exc}"
    return True, "D88 구조 정상"


def classify_state(data: bytes, target: dict) -> tuple[str, dict]:
    before = after = unknown = 0
    first_unknown: list[dict] = []
    for meta in target["records"]:
        off = meta["offset"]
        length = meta["length"]
        end = off + length
        if end > len(data):
            unknown += 1
            if len(first_unknown) < 5:
                first_unknown.append({"index": meta["index"], "offset": off, "reason": "past_eof"})
            continue
        h = sha256(data[off:end])
        if h == meta["before_sha256"]:
            before += 1
        elif h == meta["after_sha256"]:
            after += 1
        else:
            unknown += 1
            if len(first_unknown) < 5:
                first_unknown.append({"index": meta["index"], "offset": off, "reason": "hash_mismatch"})

    if unknown:
        state = "INCOMPATIBLE"
    elif before and not after:
        state = "ORIGINAL"
    elif after and not before:
        state = "ALREADY_PATCHED"
    elif before and after:
        state = "PARTIAL"
    else:
        state = "INCOMPATIBLE"
    return state, {
        "before_records": before,
        "after_records": after,
        "unknown_records": unknown,
        "unknown_samples": first_unknown,
    }


def inspect_file(target: dict, path: Path) -> dict:
    if not path.is_file():
        raise PatcherError(f"파일을 찾을 수 없습니다: {path}")
    data = path.read_bytes()
    if len(data) != target["expected_size"]:
        raise PatcherError(
            f"{target['display_name']} 크기가 다릅니다. 현재 {len(data):,} bytes / 기준 {target['expected_size']:,} bytes"
        )
    load_and_check_patch(target)
    if target["kind"] == "d88":
        ok, detail = d88_structure_ok(data)
        if not ok:
            raise PatcherError(detail)
    else:
        detail = "ROM 크기 정상"
    state, counts = classify_state(data, target)
    return {
        "path": str(path),
        "display_name": target["display_name"],
        "size": len(data),
        "full_sha256": sha256(data),
        "state": state,
        "structure": detail,
        **counts,
    }


def apply_records(base: bytes, records: list[IPSRecord]) -> bytes:
    out = bytearray(base)
    for record in records:
        end = record.offset + len(record.data)
        if record.offset < 0 or end > len(out):
            raise PatcherError(f"IPS가 입력 파일 끝을 넘어 씁니다: 0x{record.offset:06X}")
        out[record.offset:end] = record.data
    return bytes(out)


def output_path_for(path: Path) -> Path:
    return path.with_name(f"{path.stem}(K){path.suffix}")


def backup_path_for(path: Path) -> Path:
    candidate = path.with_name(path.name + ".bak")
    n = 1
    while candidate.exists():
        candidate = path.with_name(f"{path.name}.{n}.bak")
        n += 1
    return candidate


def _restore_original(backup: Path, original: Path) -> None:
    if original.exists():
        original.unlink()
    os.replace(backup, original)


def patch_file(target: dict, path: Path) -> dict:
    pre = inspect_file(target, path)
    if pre["state"] == "ALREADY_PATCHED":
        return {"status": "already_patched", "pre": pre, "backup": None, "output": str(path), "post": pre}
    if pre["state"] == "PARTIAL":
        raise PatcherError(f"{target['display_name']}가 부분 패치 상태입니다. 안전을 위해 적용하지 않습니다.")
    if pre["state"] != "ORIGINAL":
        raise PatcherError(f"{target['display_name']}가 지원 원본과 일치하지 않습니다. 안전을 위해 적용하지 않습니다.")

    output = output_path_for(path)
    if output.resolve() == path.resolve():
        raise PatcherError("출력 파일 이름을 만들 수 없습니다.")
    if output.exists():
        try:
            existing = inspect_file(target, output)
        except Exception:
            existing = None
        if existing and existing["state"] == "ALREADY_PATCHED":
            return {"status": "output_exists", "pre": pre, "backup": None, "output": str(output), "post": existing}
        raise PatcherError(f"출력 파일이 이미 존재합니다: {output}")

    records = load_and_check_patch(target)
    original = path.read_bytes()
    patched = apply_records(original, records)
    post_state, post_counts = classify_state(patched, target)
    if post_state != "ALREADY_PATCHED":
        raise PatcherError(f"메모리상 IPS 적용 결과의 사후검증에 실패했습니다: {post_state}, {post_counts}")
    if target["kind"] == "d88":
        ok, detail = d88_structure_ok(patched)
        if not ok:
            raise PatcherError(f"IPS 적용 결과의 D88 구조검증 실패: {detail}")

    backup = backup_path_for(path)
    fd, tmp_name = tempfile.mkstemp(prefix=f".{output.name}.", suffix=".patchy88.tmp", dir=str(output.parent))
    tmp = Path(tmp_name)
    tmp_consumed = False
    try:
        with os.fdopen(fd, "wb") as f:
            f.write(patched)
            f.flush()
            os.fsync(f.fileno())

        temp_data = tmp.read_bytes()
        temp_state, temp_counts = classify_state(temp_data, target)
        if temp_state != "ALREADY_PATCHED":
            raise PatcherError(f"임시파일 사후검증 실패: {temp_state}, {temp_counts}")
        if target["kind"] == "d88":
            ok, detail = d88_structure_ok(temp_data)
            if not ok:
                raise PatcherError(f"임시 D88 구조검증 실패: {detail}")

        # Current native Patchy88 behavior: move the original to .bak, then
        # promote the already-verified temporary file to the (K) output name.
        os.rename(path, backup)
        try:
            os.rename(tmp, output)
            tmp_consumed = True
        except Exception:
            _restore_original(backup, path)
            raise

        try:
            post = inspect_file(target, output)
            if post["state"] != "ALREADY_PATCHED":
                raise PatcherError("최종 파일 사후검증 실패")
        except Exception as exc:
            try:
                output.unlink(missing_ok=True)
            finally:
                _restore_original(backup, path)
            raise PatcherError(f"최종 파일 사후검증 실패. 원본을 복구했습니다: {exc}") from exc

        return {
            "status": "patched",
            "pre": pre,
            "backup": str(backup),
            "output": str(output),
            "post": post,
        }
    finally:
        if not tmp_consumed:
            tmp.unlink(missing_ok=True)


def state_korean(state: str) -> str:
    return {
        "ORIGINAL": "정상 원본",
        "ALREADY_PATCHED": "이미 패치됨",
        "PARTIAL": "부분 패치",
        "INCOMPATIBLE": "호환되지 않음",
    }.get(state, state)


def cli(argv: list[str]) -> int:
    parser = argparse.ArgumentParser(description=f"{APP_NAME} v{APP_VERSION} Python")
    parser.add_argument("--check", choices=("disk_a", "kanji1"), help="패치하지 않고 상태만 검사")
    parser.add_argument("--apply", choices=("disk_a", "kanji1"), help="검증 후 IPS 적용")
    parser.add_argument("file", nargs="?", type=Path)
    parser.add_argument("--json", action="store_true", help="결과를 JSON으로 출력")
    args = parser.parse_args(argv)
    if not (args.check or args.apply):
        return gui()
    if args.check and args.apply:
        parser.error("--check와 --apply를 동시에 사용할 수 없습니다.")
    if not args.file:
        parser.error("대상 파일을 지정해야 합니다.")

    manifest = load_manifest()
    target = target_by_id(manifest, args.check or args.apply)
    try:
        result = inspect_file(target, args.file) if args.check else patch_file(target, args.file)
        if args.json:
            print(json.dumps(result, ensure_ascii=False, indent=2))
        elif args.check:
            print(f"{target['display_name']}: {state_korean(result['state'])}")
            print(f"SHA-256: {result['full_sha256']}")
        else:
            status = result["status"]
            if status == "patched":
                print("패치 완료")
                print(f"백업: {result['backup']}")
                print(f"출력: {result['output']}")
            elif status == "output_exists":
                print("정상 패치 결과 파일이 이미 존재합니다. 변경하지 않았습니다.")
                print(f"출력: {result['output']}")
            else:
                print("선택한 파일은 이미 패치되어 있습니다. 변경하지 않았습니다.")
            print(f"결과 SHA-256: {result['post']['full_sha256']}")
        return 0
    except (OSError, PatcherError) as exc:
        print(f"오류: {exc}", file=sys.stderr)
        return 2


def gui() -> int:
    try:
        import tkinter as tk
        from tkinter import filedialog, messagebox, ttk
    except Exception as exc:
        print(f"GUI를 시작할 수 없습니다: {exc}", file=sys.stderr)
        return 2

    manifest = load_manifest()
    disk_target = target_by_id(manifest, "disk_a")
    kanji_target = target_by_id(manifest, "kanji1")

    root = tk.Tk()
    root.title(f"Patchy88 - 무한전기 바리스 한국어 패치 v{APP_VERSION} (Python)")
    root.geometry("840x680")
    root.minsize(780, 620)

    style = ttk.Style(root)
    try:
        style.theme_use("vista")
    except tk.TclError:
        pass

    outer = ttk.Frame(root, padding=(18, 16, 18, 14))
    outer.pack(fill="both", expand=True)

    title = ttk.Label(outer, text="무한전기 바리스 한국어 패치", font=("Segoe UI", 16, "bold"))
    title.pack(anchor="w")
    ttk.Label(
        outer,
        text="Disk A와 KANJI1 ROM의 실제 패치 대상 데이터를 검증한 뒤 안전하게 IPS를 적용합니다.",
        font=("Segoe UI", 9),
        wraplength=780,
    ).pack(anchor="w", pady=(4, 14))

    vars_by_id = {"disk_a": tk.StringVar(), "kanji1": tk.StringVar()}

    def choose(target_id: str, title_text: str, filetypes):
        filename = filedialog.askopenfilename(parent=root, title=title_text, filetypes=filetypes)
        if filename:
            vars_by_id[target_id].set(filename)

    def target_row(target: dict, filetypes):
        frame = ttk.LabelFrame(outer, text=target["display_name"], padding=(10, 9))
        frame.pack(fill="x", pady=5)
        entry = ttk.Entry(frame, textvariable=vars_by_id[target["id"]], font=("Segoe UI", 9))
        entry.pack(side="left", fill="x", expand=True)
        ttk.Button(
            frame,
            text="찾아보기...",
            command=lambda: choose(target["id"], f"{target['display_name']} 선택", filetypes),
            width=12,
        ).pack(side="left", padx=(8, 0))

    target_row(disk_target, [("D88 image", "*.d88"), ("All files", "*.*")])
    target_row(kanji_target, [("ROM", "*.rom *.ROM"), ("All files", "*.*")])

    log_frame = ttk.LabelFrame(outer, text="상태 / 로그", padding=8)
    log_frame.pack(fill="both", expand=True, pady=(10, 8))
    log = tk.Text(log_frame, height=13, wrap="word", state="disabled", font=("Consolas", 9), relief="flat")
    scrollbar = ttk.Scrollbar(log_frame, orient="vertical", command=log.yview)
    log.configure(yscrollcommand=scrollbar.set)
    scrollbar.pack(side="right", fill="y")
    log.pack(side="left", fill="both", expand=True)

    def write_log(text: str):
        log.configure(state="normal")
        log.insert("end", text.rstrip() + "\n")
        log.see("end")
        log.configure(state="disabled")
        root.update_idletasks()

    def inspect_selected():
        any_selected = False
        for target in (disk_target, kanji_target):
            value = vars_by_id[target["id"]].get().strip()
            if not value:
                continue
            any_selected = True
            try:
                result = inspect_file(target, Path(value))
                write_log(
                    f"[{target['display_name']}] {state_korean(result['state'])} | "
                    f"before={result['before_records']} after={result['after_records']} unknown={result['unknown_records']}\n"
                    f"  SHA-256 {result['full_sha256']}"
                )
            except Exception as exc:
                write_log(f"[{target['display_name']}] 오류: {exc}")
        if not any_selected:
            messagebox.showinfo(APP_NAME, "검사할 Disk A 또는 KANJI1 ROM을 선택해 주세요.", parent=root)

    def apply_selected():
        selected = [(t, vars_by_id[t["id"]].get().strip()) for t in (disk_target, kanji_target)]
        selected = [(t, v) for t, v in selected if v]
        if not selected:
            messagebox.showinfo(APP_NAME, "패치할 Disk A 또는 KANJI1 ROM을 선택해 주세요.", parent=root)
            return
        summary = []
        for target, value in selected:
            try:
                write_log(f"[{target['display_name']}] 검증 시작: {value}")
                result = patch_file(target, Path(value))
                if result["status"] == "patched":
                    msg = f"패치 완료 | 백업: {result['backup']} | 출력: {result['output']}"
                elif result["status"] == "output_exists":
                    msg = f"정상 `(K)` 결과가 이미 존재함 - 변경하지 않음 | 출력: {result['output']}"
                else:
                    msg = "이미 패치됨 - 변경하지 않음"
                write_log(f"[{target['display_name']}] {msg}\n  결과 SHA-256 {result['post']['full_sha256']}")
                summary.append(f"{target['display_name']}: {msg}")
            except Exception as exc:
                write_log(f"[{target['display_name']}] 실패: {exc}")
                summary.append(f"{target['display_name']}: 실패 - {exc}")
        messagebox.showinfo(APP_NAME, "\n".join(summary), parent=root)

    buttons = ttk.Frame(outer)
    buttons.pack(fill="x", pady=(2, 0))
    ttk.Button(buttons, text="상태 검사", command=inspect_selected, width=14).pack(side="left")
    ttk.Button(buttons, text="한글 패치 적용", command=apply_selected, width=18).pack(side="left", padx=8)
    ttk.Button(buttons, text="종료", command=root.destroy, width=12).pack(side="right")

    ttk.Separator(outer, orient="horizontal").pack(fill="x", pady=(12, 9))
    ttk.Label(
        outer,
        text=(
            "패치 성공 시 원본은 같은 폴더의 <원본파일명>.bak 으로 보관되고, "
            "패치 결과는 <원본명>(K)<확장자> 형식으로 생성됩니다. "
            "기존 백업은 덮어쓰지 않고 .1.bak, .2.bak 순으로 보존합니다."
        ),
        font=("Segoe UI", 9),
        wraplength=790,
    ).pack(anchor="w")

    root.mainloop()
    return 0


def main(argv: list[str] | None = None) -> int:
    try:
        return cli(sys.argv[1:] if argv is None else argv)
    except PatcherError as exc:
        print(f"오류: {exc}", file=sys.stderr)
        return 2


if __name__ == "__main__":
    raise SystemExit(main())
