# kevmo314/tinywhip

A simple, lightweight [WHIP](https://datatracker.ietf.org/doc/html/draft-ietf-wish-whip)/[WHEP](https://www.ietf.org/id/draft-murillo-whep-01.html) server. It accepts video as WHIP and publishes video as WHEP.

This server is intended to make testing and development of WHIP/WHEP clients easier as well as to provide a minimal, easy to understand and hack implementation of the protocol. It is not intended for production use.

Looking for a production-ready WHIP/WHEP server? Check out [LiveKit](https://livekit.io/)! It's a full-featured, scalable, and production-ready video conferencing platform that supports WHIP/WHEP.

## Usage

### Docker

```bash
docker run -p 8080:8080 kevmo314/tinywhip
```

### Binary

```bash
go run cmd/main.go
```

Optionally,

```bash
PORT=8080 go run cmd/main.go
```

## Example

Run the server

```bash
go run cmd/main.go
```

In a separate terminal, run the client

```bash
go run examples/publisher/main.go http://localhost:8080 testvideo
```

This will print out a track id in the server:

```
2023/01/13 12:09:47 Adding track: 3e3c9c71-3a2d-4a72-a21f-8f9cae60a4ed
```

Then in one more terminal, run the subscriber with this track id

```bash
go run examples/subscriber/main.go http://localhost:8080/3e3c9c71-3a2d-4a72-a21f-8f9cae60a4ed
```

For a more fun (and realistic) experience, check out the [OBS WHIP](https://github.com/obsproject/obs-studio/pull/7926) support or [ffmpeg-whip](https://github.com/kevmo314/ffmpeg-whip) to publish real video instead of a test pattern.

## API

The WHIP endpoint is `/`. An SDP posted to that endpoint will be ingested according to the WHIP specification.

The WHEP endpoints are `/<id>`, where `<id>` is the stream ID. The WHEP stream will be published according to the WHEP specification.
