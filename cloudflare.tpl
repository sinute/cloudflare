apiVersion: v1
kind: Secret
metadata:
  name: cloudflare-secret
data:
  CF_API_EMAIL: $CF_API_EMAIL
  CF_API_KEY: $CF_API_KEY
  CF_DNS_NAME: $CF_DNS_NAME
  CF_DNS_TTL: $CF_DNS_TTL
  CF_ZONE_NAME: $CF_ZONE_NAME
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: cloudflare
  labels:
    app: cloudflare
spec:
  replicas: 1
  selector:
    matchLabels:
      app: cloudflare
  template:
    metadata:
      labels:
        app: cloudflare
    spec:
      containers:
      - name: cloudflare
        image: $CI_REGISTRY_IMAGE:latest
      envFrom:
      - secretRef:
        name: cloudflare-secret
      restartPolicy: Always
