FROM golang:1.25 AS build

WORKDIR /
COPY . /

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-s' -o main . && chmod +x main

# Use scratch image to reduce the size of the image
FROM scratch

COPY --from=build main /

CMD [ "/main" ]
