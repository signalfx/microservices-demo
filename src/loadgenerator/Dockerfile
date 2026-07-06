FROM python:3.13.14-alpine3.24 AS builder

RUN apk add --no-cache \
        build-base \
        libffi-dev \
        zeromq-dev

COPY requirements.txt .
RUN pip install --no-cache-dir --prefix=/install -r requirements.txt

FROM python:3.13.14-alpine3.24

RUN apk add --no-cache bash curl && \
    addgroup -S app && adduser -S -G app app

COPY --from=builder /install /usr/local

WORKDIR /app
COPY --chown=app:app . .
RUN sed -i 's/\r$//' ./loadgen.sh && chmod +x ./loadgen.sh

ENV PYTHONUNBUFFERED=1 \
    PYTHONDONTWRITEBYTECODE=1
USER app
ENTRYPOINT ["./loadgen.sh"]
