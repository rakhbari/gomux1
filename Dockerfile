# syntax=docker/dockerfile:1

#
# Build stage
#
FROM golang:1.24-alpine AS build

# Install make and git
RUN apk add --no-cache make git

WORKDIR /build

COPY . .
COPY ./.git /build/.git

RUN make build && \
    ls -la gomux1

#
# Final stage
#
FROM alpine:3.19

LABEL maintainer="rakhbari"

WORKDIR /app

COPY --from=build /build/gomux1 ./
COPY --from=build /build/*.json ./

EXPOSE 8080
EXPOSE 8443

ENTRYPOINT [ "./gomux1" ]
