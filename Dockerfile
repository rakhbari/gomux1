# syntax=docker/dockerfile:1

#
# Larger Golang image as "build"
#
FROM golang:1.19-alpine AS build

# add a user here because addgroup and adduser are not available in scratch
RUN addgroup -g 10000 gomux1 && \
    adduser -D --uid=10000 --ingroup=gomux1 gomux1

# add make tool
RUN apk add --update make

WORKDIR /build

COPY . .

RUN make build

#
# Much smaller "final" image to run the app
#
FROM alpine:latest AS final

LABEL maintainer="rakhbari"

# copy users from build image
COPY --from=build /etc/passwd /etc/passwd

WORKDIR /app
USER gomux1

COPY --from=build --chown=gomux1:gomux1 /build/gomux1 ./
COPY --from=build --chown=gomux1:gomux1 /build/*.json ./

EXPOSE 8080
EXPOSE 8443

ENTRYPOINT [ "./gomux1" ]
