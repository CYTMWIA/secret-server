FROM python:3.9-bookworm

WORKDIR /app

COPY config ./config
COPY scripts ./scripts
COPY secret_server ./secret_server
COPY requirements.txt ./

RUN scripts/install.sh && pip cache purge

CMD scripts/run.sh
