FROM golang:1.20.5-alpine AS build

WORKDIR /src
COPY . .
RUN go build -v .

FROM scratch

COPY --from=build /src/quokka /usr/bin

ENTRYPOINT [ "quokka" ]
