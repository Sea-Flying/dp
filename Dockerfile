# This file is a template, and might need editing before it works on your project.
FROM golang:1.14 AS builder
ENV GOPROXY=https://goproxy.cn
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -v -a -o dp .

FROM scratch
# Since we started from scratch, we'll copy the SSL root certificates from the builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
WORKDIR /app/jobHclDir
WORKDIR /app/jobTplDir
WORKDIR /usr/local/bin
COPY --from=builder /app/dp .
COPY --from=builder /app/resources/dp.yml ./dp.yml
CMD ["./dp"]
