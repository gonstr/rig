FROM golang:alpine as builder
RUN apk update && apk add --no-cache git
RUN mkdir /build 
ADD . /build/
WORKDIR /build 
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o rig .

FROM alpine
RUN apk update && apk add --no-cache git
COPY --from=builder /build/rig /app/
ENV PATH="/app:${PATH}"
CMD ["rig"]
