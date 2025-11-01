FROM python:3.13-slim

ARG model
ENV MODEL_NAME=${model}
ENV HOST=0.0.0.0
ENV PORT=8000

WORKDIR /app

COPY ${model}/requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY abstract ./abstract
COPY ${model} ./model

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 8000

ENTRYPOINT ["/entrypoint.sh"]
CMD ["uvicorn"]
