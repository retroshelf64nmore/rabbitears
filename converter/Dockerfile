# vhsify drop-folder converter — runs on a Raspberry Pi (arm64/armhf) or any x86 box.
FROM debian:bookworm-slim
RUN apt-get update \
 && apt-get install -y --no-install-recommends ffmpeg ca-certificates \
 && rm -rf /var/lib/apt/lists/*
COPY vhsify watch.sh /usr/local/bin/
RUN chmod +x /usr/local/bin/vhsify /usr/local/bin/watch.sh
VOLUME ["/input", "/output"]
ENTRYPOINT ["/usr/local/bin/watch.sh"]
