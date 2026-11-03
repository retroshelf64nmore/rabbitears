#!/usr/bin/env bash
# Drop-folder service: anything that lands in /input is converted with vhsify
# into /output, keeping sub-folders (a folder = one channel that airs its
# videos in order). Polls instead of using inotify so it also works on
# network shares and USB drives.
set -uo pipefail

IN=${INPUT_DIR:-/input}
OUT=${OUTPUT_DIR:-/output}
INTERVAL=${INTERVAL:-60}          # seconds between scans
SETTLE=${SETTLE:-2}               # scans a file must stay unchanged before converting
WEAR=${WEAR:-1}
FIT=${FIT:-letterbox}
SIZE=${SIZE:-320x240}
CRF=${CRF:-24}
PRESET=${PRESET:-medium}
STATE=${STATE_DIR:-/output/.vhsify}
EXTS="mp4 m4v mkv avi mov wmv flv webm mpg mpeg ts m2ts vob 3gp ogv divx"

mkdir -p "$OUT" "$STATE"
declare -A seen count
echo "vhsify drop folder: $IN → $OUT (wear $WEAR, $FIT, $SIZE, every ${INTERVAL}s)"

opts=(--wear "$WEAR" --size "$SIZE" --crf "$CRF" --preset "$PRESET")
[ "$FIT" = fill ] && opts+=(--fill)

while :; do
  while IFS= read -r -d '' f; do
    rel="${f#"$IN"/}"
    name="$(basename "$rel")"
    case "$name" in .*|*.part|*.part.mp4|*.tmp|*.!qB) continue ;; esac
    ext="${name##*.}"; ext="${ext,,}"
    [[ " $EXTS " == *" $ext "* ]] || continue
    out="$OUT/${rel%.*}.mp4"
    [ -s "$out" ] && [ "$out" -nt "$f" ] && continue
    # skip files that failed before until they change
    sig="$(stat -c '%s:%Y' "$f" 2>/dev/null)" || continue
    failmark="$STATE/$(printf '%s' "$rel" | md5sum | cut -c1-16).failed"
    [ -f "$failmark" ] && [ "$(cat "$failmark")" = "$sig" ] && continue
    # wait until the file has stopped growing (still copying/downloading)
    if [ "${seen[$rel]:-}" = "$sig" ]; then
      count[$rel]=$(( ${count[$rel]:-0} + 1 ))
    else
      seen[$rel]="$sig"; count[$rel]=0
    fi
    [ "${count[$rel]}" -ge "$SETTLE" ] || continue
    echo "$(date '+%F %T') converting $rel"
    if vhsify -q "${opts[@]}" "$f" --to "$out"; then
      rm -f "$failmark"
      echo "$(date '+%F %T') done $out"
    else
      echo "$sig" > "$failmark"
      echo "$(date '+%F %T') FAILED $rel (will retry if the file changes)"
    fi
  done < <(find "$IN" -type f -print0 2>/dev/null | sort -z)
  sleep "$INTERVAL"
done
