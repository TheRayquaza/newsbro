FROM python:3.13-slim

ARG model
ENV MODEL_NAME=${model}
ENV HOST=0.0.0.0
ENV PORT=8080

WORKDIR /app

# hadolint ignore=DL3008,DL3015
RUN apt-get update && apt-get install -y --no-install-recommends \
    bash=5.2.15-2+b9 \
    ca-certificates=20230311+deb12u1 \
    curl=7.88.1-10+deb12u14 \
    && rm -rf /var/lib/apt/lists/*
RUN useradd -m -s /bin/bash nonroot

COPY ${model}/requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY abstract ./abstract
COPY ${model} ./${model}
COPY ${model}/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh && useradd -m -s /bin/bash appuser

USER appuser

HEALTHCHECK --interval=30s --timeout=5s --start-period=120s --retries=6 \
  CMD curl -f http://${HOST}:${PORT}/health || exit 1

ENTRYPOINT ["/entrypoint.sh"]
CMD ["python3"]
