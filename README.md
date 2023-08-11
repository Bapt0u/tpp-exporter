# venafi-exporter

This exporter exposes [Venafi](https://venafi.com/) metrics in prometheus format such as total of enrolled certificates or active or expire soon ones. 

## Config file

Create a configuration file at `./conf/venafi.yml` with the following variables.

```yml
---
venafi_tpp:
  username: my_user                    # Username
  password: password                   # User password
  client_id: client-id                 # User client id
  scope: certificate:manage,approve,revoke;configuration:manage
  url: https://venafi.dev.localhost    # Venafi URL
  scrapetime: 3                        # Scrape time in minutes
```

## Metrics available

| Metric Name | Description |
|---|---|
| venafi_certificate_total | Total of available certificates. |
| venafi_certificate_expiring_soon | Number of certificate expiring soon. |
| venafi_certificate_per_policy | Number of certificate in each policy. |
| venafi_certificate_valid_per_policy | Number of valid certificate in each policy. |

## Docker build from source

Here is an example of `Dockerfile` you can use. The default port exposed is `2112`. 

```Dockerfile
FROM golang:1.20 AS build

# Build from sources
WORKDIR /app
COPY . ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /venafi_exporter

# Change environment
FROM alpine:latest as run

WORKDIR /app
COPY --from=build /venafi_exporter /app/venafi_exporter
EXPOSE 2112
ENTRYPOINT ["/app/venafi_exporter"]
```

Save this Dockerfile and run `docker build` command.

```bash
docker build --tag registry.dev.localhost/venafi-exporter:1.0 .
```

Finally, run the container with `docker run`.

```bash
docker run -d --restart=unless-stopped  --name=venafi-exporter -v $(pwd)/conf/venafi.yml:/app/conf/venafi.yml -p 2112:2112 registry.dev.localhost/venafi-exporter:1.0
```

It is important to bind config file in the path `/app/conf/venafi.yml` to allow the binary to read its configuration.  

## Todo

- [ ] Configure listening port in config file
- [ ] Scrapetime in seconds/minutes
- [ ] Fetch Expire soon certificate by DN