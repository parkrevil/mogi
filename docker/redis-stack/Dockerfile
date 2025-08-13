FROM redis/redis-stack-server:7.4.5

# Set environment variables for RedisGears
ENV REDISGEARS_V8_PLUGIN_PATH=/opt/redis-stack/lib/libredisgears_v8_plugin.so
ENV LD_LIBRARY_PATH=/opt/redis-stack/lib:$LD_LIBRARY_PATH

COPY redis.conf /etc/redis/redis.conf
COPY init-users.acl /etc/redis/init-users.acl
EXPOSE 6379

CMD ["redis-server", "/etc/redis/redis.conf"]
