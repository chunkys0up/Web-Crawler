FROM golang:1.24

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Ensure wait-for-it.sh is executable
# RUN chmod +x wait-for-it.sh

RUN go build -o crawler

CMD ["./crawler"]
