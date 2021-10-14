# Analyser

Check commit sha &amp; cookstyle against cache

Run command:

```bash
docker run --rm --network=analyser_redis -e REDIS_HOST="redis" -e REDIS_PORT="6379" -e REDIS_PASSWORD="${REDIS_PASSWORD}" -e ORGANISATION=stylelia -e NAME=snort -v "$PWD":/var/task:ro,delegated lambci/lambda:go1.x analyser '{"some": "event"}'
```
