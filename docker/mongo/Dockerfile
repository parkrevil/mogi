FROM mongo:8.0

COPY mongod.conf /etc/mongod.conf
RUN openssl rand -base64 756 > /data/keyfile && \
  chmod 400 /data/keyfile && \
  chown 999:999 /data/keyfile
