# Build

```sh
go build
```

# Run

```sh
# ddns for test.example.com
./cloudflare -CF_API_EMAIL user@example.com -CF_API_KEY xxxxxx -CF_DNS_NAME test -CF_DNS_TTL 120 -CF_ZONE_NAME example.com
# or
CF_API_EMAIL=user@example.com CF_API_KEY=xxxxxx CF_DNS_NAME=test CF_DNS_TTL=120 CF_ZONE_NAME=example.com ./cloudflare
```

> You can find your API_KEY [here](https://dash.cloudflare.com/profile/api-tokens)
