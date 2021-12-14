FROM golang:1.17-alpine as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main . 

FROM golang:1.17-alpine as builderrtpforwarder
WORKDIR /app
COPY rtpforwarder/go.mod rtpforwarder/go.sum ./
RUN go mod download
COPY rtpforwarder/ .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o rtpforwarder .

FROM alpine:latest
WORKDIR /app
RUN apk add busybox-extras && apk add ffmpeg
ENV RTP_SESSION "eyJkdW1teXVzZXIiIDogIjAwMDAifQ=="
ENV PORT_USED "WyI0MDA4Il0="
ENV SDP_DIRECTORY "sdpcollection"
ENV DASH_SERVER "https://bifrost.inlive.app"
ENV MANIFEST_FILENAME "manifest.mpd"
ENV SDP_FILENAME "rtpforwarder.sdp"
ENV FFMPEG_INSTANCES "WyI0MDA4Il0="
COPY .env .
COPY --from=builder /app/main .
COPY --from=builderrtpforwarder /app/rtpforwarder .
RUN mkdir logs && chmod -R 777 logs
RUN mkdir sdpcollection && chmod -R 777 sdpcollection


EXPOSE 9090
CMD ["./main"]
