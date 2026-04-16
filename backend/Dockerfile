FROM golang:1.22-alpine AS build

WORKDIR /app
COPY go.mod ./
COPY main.go ./

RUN go build -o /out/opsdash-backend .

FROM alpine:3.20
RUN adduser -D -H -u 10001 appuser

WORKDIR /app
COPY --from=build /out/opsdash-backend /app/opsdash-backend

USER 10001
ENV PORT=8080
EXPOSE 8080

ENTRYPOINT ["/app/opsdash-backend"]
