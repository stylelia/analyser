# Analyser

Check commit sha &amp; cookstyle against cache

Run command:

```bash
docker build . -t foo

go build . ;docker run --rm --network=analyser_redis -e REDIS_HOST="redis" -e REDIS_PORT="6379" -e REDIS_PASSWORD="${REDIS_PASSWORD}" -e ORGANISATION=stylelia -e GITHUB_TOKEN="${GITHUB_TOKEN}" -e NAME=snort -v "$PWD":/var/task:ro,delegated foo analyser '{"some": "event"}'
```
