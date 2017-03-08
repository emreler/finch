# Finch - Persist Alerts

This package is used to avoid lost messages from Redis when the actual Finch binary is not working, during deployment etc.

Working mechanism of Finch is based on expired keys on Redis, and it is using pub/sub feature of Redis for getting notified for expired keys. Since with pub/sub when there is no subscriber attached to a key, messages are simply lost. Sole purpose of this package is to listen to that key, and persist them to a list, which is to be consumed by Finch.

## Running

It should be running alongside Finch and they must be connected to the same Redis instance.

See `docker-compose.yml` in root directory.
