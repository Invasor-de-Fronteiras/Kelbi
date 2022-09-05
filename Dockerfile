FROM golang:1.16-alpine as build

WORKDIR /app

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

COPY . .

RUN go build -o main main.go


# Imagem de Origem
FROM node:13-alpine as web
WORKDIR /app
ENV PATH /app/node_modules/.bin:$PATH
COPY www/kelbi /app
RUN npm ci && npm run build

# stage 2: copy only the application binary file and necessary files to the alpine container
FROM alpine:3.12
RUN apk --update add ca-certificates

WORKDIR /app

COPY --from=build /app/main .
COPY --from=web /app/build ./www/kelbi/build 

EXPOSE 80
EXPOSE 53312
EXPOSE 54001
EXPOSE 54002
EXPOSE 54003
EXPOSE 54004
EXPOSE 53310

# run the service on container startup.
CMD ["/app/main"]
