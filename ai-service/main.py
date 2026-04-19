import tempfile
import os
import json
from faster_whisper import WhisperModel
from fastapi import FastAPI, UploadFile, File, Form
from fastapi.responses import FileResponse
from pydantic import BaseModel
from ner import find_entities, redact_text, Entity
from audio_redact import find_silence_ranges, silence_audio
from llm import find_entities_llm

app = FastAPI(title="Voice Redactor")

device = os.environ.get("WHISPER_DEVICE", "cpu")
model_size = os.environ.get("WHISPER_MODEL", "medium")
whisper_language = os.environ.get("WHISPER_LANGUAGE", "ru")
compute_type = "float16" if device == "cuda" else "int8"

model = WhisperModel(model_size, device=device, compute_type=compute_type)


class ProcessRequest(BaseModel):
    job_id: str
    file_path: str


def get_entities(text: str) -> list[Entity]:
    entities = find_entities(text)
    found_values = {e.value for e in entities}

    for item in find_entities_llm(text):
        value = item.get("value", "")
        if not value or value in found_values:
            continue
        pos = text.find(value)
        if pos == -1:
            continue
        entities.append(Entity(
            type=item.get("type", "unknown"),
            value=value,
            start_char=pos,
            end_char=pos + len(value),
        ))
        found_values.add(value)

    entities.sort(key=lambda e: e.start_char)
    return entities


def transcribe_file(path: str) -> list[dict]:
    segments, _ = model.transcribe(
        path,
        language=whisper_language,
        word_timestamps=True,
        beam_size=5,
        vad_filter=True,
    )
    words = []
    for segment in segments:
        for w in segment.words:
            words.append({
                "word": w.word.strip(),
                "start": round(w.start, 3),
                "end": round(w.end, 3),
            })
    return words


@app.get("/health")
def health():
    return {"status": "ok"}


@app.post("/transcribe")
def transcribe(file: UploadFile = File(...)):
    with tempfile.NamedTemporaryFile(delete=False, suffix=os.path.splitext(file.filename)[1]) as tmp:
        tmp.write(file.file.read())
        tmp_path = tmp.name

    try:
        words = transcribe_file(tmp_path)
    finally:
        os.remove(tmp_path)

    transcript = " ".join(w["word"] for w in words)
    entities = get_entities(transcript)

    return {
        "words": words,
        "transcript": transcript,
        "redacted_transcript": redact_text(transcript, entities),
        "pii_events": [
            {"type": e.type, "value": e.value, "start_char": e.start_char, "end_char": e.end_char}
            for e in entities
        ],
    }


@app.post("/redact-audio")
def redact_audio(
    file: UploadFile = File(...),
    words: str = Form(...),
    pii_events: str = Form(...),
):
    with tempfile.NamedTemporaryFile(delete=False, suffix=".wav") as tmp_in:
        tmp_in.write(file.file.read())
        input_path = tmp_in.name

    output_path = input_path.replace(".wav", "_redacted.wav")

    try:
        silence_ranges = find_silence_ranges(json.loads(words), json.loads(pii_events))
        silence_audio(input_path, output_path, silence_ranges)
    finally:
        os.remove(input_path)

    return FileResponse(output_path, media_type="audio/wav", filename="redacted.wav")


@app.post("/process")
def process(req: ProcessRequest):
    words = transcribe_file(req.file_path)
    transcript = " ".join(w["word"] for w in words)
    entities = get_entities(transcript)
    redacted = redact_text(transcript, entities)

    ext = req.file_path.rsplit(".", 1)[-1]
    redacted_audio_path = req.file_path.rsplit(".", 1)[0] + "_redacted." + ext
    pii_dicts = [
        {"type": e.type, "value": e.value, "start_char": e.start_char, "end_char": e.end_char}
        for e in entities
    ]
    silence_ranges = find_silence_ranges(words, pii_dicts)
    silence_audio(req.file_path, redacted_audio_path, silence_ranges, ext)

    pii_events = []
    for e in entities:
        matched = [w for w in words if w["word"] in e.value]
        pii_events.append({
            "type": e.type,
            "original": e.value,
            "start_sec": matched[0]["start"] if matched else 0.0,
            "end_sec": matched[-1]["end"] if matched else 0.0,
        })

    return {
        "transcript": transcript,
        "redacted_transcript": redacted,
        "redacted_audio_path": redacted_audio_path,
        "pii_events": pii_events,
        "words": words,
    }
