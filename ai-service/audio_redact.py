from pydub import AudioSegment


def find_silence_ranges(words: list[dict], pii_events: list[dict]) -> list[tuple[float, float]]:
    char_offset = 0
    word_positions = []
    for w in words:
        word_positions.append({
            "start_time": w["start"],
            "end_time": w["end"],
            "start_char": char_offset,
            "end_char": char_offset + len(w["word"]),
        })
        char_offset += len(w["word"]) + 1

    ranges = []
    for event in pii_events:
        event_start = event["start_char"]
        event_end = event["end_char"]

        overlapping = [
            wp for wp in word_positions
            if wp["start_char"] < event_end and wp["end_char"] > event_start
        ]

        if overlapping:
            t_start = min(wp["start_time"] for wp in overlapping)
            t_end = max(wp["end_time"] for wp in overlapping)
            ranges.append((t_start, t_end))

    return ranges


def silence_audio(input_path: str, output_path: str, silence_ranges: list[tuple[float, float]], fmt: str = "wav"):
    audio = AudioSegment.from_file(input_path)

    for t_start, t_end in silence_ranges:
        start_ms = int(t_start * 1000)
        end_ms = int(t_end * 1000)

        silence = AudioSegment.silent(duration=end_ms - start_ms)
        audio = audio[:start_ms] + silence + audio[end_ms:]

    audio.export(output_path, format=fmt)