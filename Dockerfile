# Use multistage builds.

# Create Golang builder stage.
## Fetch OS updates and app dependencied.
FROM golang:1.14.7-alpine3.11 as builder
RUN apk update && apk add --no-cache git make
## Copy app source to image build directory.
RUN mkdir /build
ADD . /build/
## Build app from source in build directory.
WORKDIR /build
RUN make prepare
RUN GOOS=linux make build

# Create final app image from builder stage.
FROM alpine:3.11
## Add application from builder stage to final image.
COPY --from=builder /build/rest-client /app/
## Set directory to run in.
WORKDIR /app
## Set command to run.
CMD ["./rest-client"]
