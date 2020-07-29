# Stage 0: Builder
FROM golang:1.14-alpine as builder
COPY . /build
WORKDIR /build
RUN apk add --update make git && make

# Stage 1: Final image
FROM scratch
COPY --from=builder /build/stitch /stitch
CMD ["/stitch"]
