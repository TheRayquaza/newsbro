FROM rustlang/rust:nightly as builder

RUN apt-get update && apt-get install -y musl-tools && \
    rustup target add x86_64-unknown-linux-musl

WORKDIR /app
COPY Cargo.toml Cargo.lock ./
COPY src ./src

ENV CARGO_UNSTABLE_EDITION2024=1

RUN cargo +nightly build --release --target x86_64-unknown-linux-musl

FROM alpine:3.18

RUN apk add --no-cache \
    ca-certificates \
    curl

WORKDIR /app

COPY --from=builder /app/target/x86_64-unknown-linux-musl/release/srvc-drift /app/srvc-drift

RUN adduser -D -u 1000 driftuser && \
    chown -R driftuser:driftuser /app

USER driftuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD sh -c "/usr/bin/curl -f http://localhost:8080/health || exit 1"

CMD ["/app/srvc-drift"]
