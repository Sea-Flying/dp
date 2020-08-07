# This file is a template, and might need editing before it works on your project.
FROM golang:1.13-alpine AS builder
ENV GOPROXY=https://goproxy.cn
# We'll likely need to add SSL root certificates
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories \
    && apk --no-cache add ca-certificates
WORKDIR /app
COPY . .
RUN GOOS=linux go build -v -a -o dp .

FROM scratch
# Since we started from scratch, we'll copy the SSL root certificates from the builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
WORKDIR /app/jobHclDir
WORKDIR /app/jobTplDir
WORKDIR /usr/local/bin
COPY --from=builder /app/dp .
COPY --from=builder /app/resources/dp.yml ./dp.yml
CMD ["./dp"]
