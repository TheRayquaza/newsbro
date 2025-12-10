FROM rustlang/rust:nightly as builder

WORKDIR /app
COPY Cargo.toml Cargo.lock ./
COPY src ./src

ENV RUSTFLAGS=""
ENV CARGO_UNSTABLE_EDITION2024=1

RUN cargo +nightly build --release

FROM alpine:3.18

RUN apk add --no-cache \
    ca-certificates \
    curl \
    libssl3

WORKDIR /app

COPY --from=builder /app/target/release/srvc-drift /app/srvc-drift

RUN adduser -D -u 1000 driftuser && \
    chown -R driftuser:driftuser /app

USER driftuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD ["/usr/bin/curl", "-f", "http://localhost:8080/health"] || exit 1

CMD ["/app/srvc-drift"]
