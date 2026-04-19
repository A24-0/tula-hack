FROM python:3.11-slim

RUN apt-get update && apt-get install -y ffmpeg curl && rm -rf /var/lib/apt/lists/*

WORKDIR /app

ENV HF_HOME=/models \
    XDG_CACHE_HOME=/models \
    HF_HUB_CACHE=/models
RUN mkdir -p /models

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY . .

CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8000"]