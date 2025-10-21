FROM python:3.13-slim

ARG model

WORKDIR /app

COPY apps/srvc-inference/${model}/requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY apps/srvc-inference/abstract ./abstract
COPY apps/srvc-inference/${model} ./model

ENV MODEL_NAME=${model}

EXPOSE 8000

CMD ["uvicorn", "${MODEL_NAME}.main:app", "--host", "0.0.0.0", "--port", "8000"]
