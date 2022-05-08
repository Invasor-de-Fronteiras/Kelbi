FROM golang:1.15.7-alpine3.13

ENV GO111MODULE=on

WORKDIR /app/erupe

COPY go.mod .
COPY go.sum .

RUN go mod download

# RUN mkdir www
# RUN mkdir savedata
# RUN mkdir bin

# COPY . .

# # COPY bin .
# COPY common .
# COPY config .
# COPY Gifts .
# COPY migrations .
# COPY network .
# # COPY savedata .
# COPY server .
# # COPY www .

# COPY config.json .
# COPY main.go .

CMD [ "go", "run", "." ]