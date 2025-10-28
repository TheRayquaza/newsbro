FROM python:3.12-slim

ARG model

WORKDIR /app

COPY ${model}/requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY abstract ./abstract
COPY ${model} ./model

ENV MODEL_NAME=${model}

EXPOSE 8000

CMD ["uvicorn", "${MODEL_NAME}.main:app", "--host", "0.0.0.0", "--port", "8000"]
