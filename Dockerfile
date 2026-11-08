FROM alpine:3.24

RUN apk update \
  && apk add --no-cache \
    bash \
    curl \
    ca-certificates \
    git

RUN adduser -D -u 10001 pig-user
USER pig-user

ENV PATH="/home/pig-user/.local/share/mise/shims:/home/pig-user/.local/bin:${PATH}"
ENV MISE_TRUSTED_CONFIG_PATHS=/workspace

RUN curl -fsSL https://mise.run | sh \
    && mise use -g \
        go@latest \
        github:MichaelKinsy/PiG@latest \
    && rm -rf ~/.cache/mise ~/.local/share/mise/downloads

COPY --chown=pig-user:pig-user extensions/plan /opt/pig-extensions/plan
WORKDIR /workspace
ENTRYPOINT ["pig", "--extension", "/opt/pig-extensions/plan"]
