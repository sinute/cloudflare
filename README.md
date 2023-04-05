# Build

```sh
make build
```

# Run

```sh
# ddns for test.example.com
./bin/cloudflare -CF_API_EMAIL user@example.com -CF_API_KEY xxxxxx -CF_DNS_NAME foo -CF_DNS_TTL 120 -CF_ZONE_NAME 
example.com
# or
CF_API_EMAIL=user@example.com CF_API_KEY=xxxxxx CF_DNS_NAME=foo CF_DNS_TTL=120 CF_ZONE_NAME=example.com ./bin/cloudflare
```

> You can find your API_KEY [here](https://dash.cloudflare.com/profile/api-tokens)
