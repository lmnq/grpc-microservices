FROM golang AS builder
WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o payment ./cmd/main.go

FROM scratch
COPY --from=builder /usr/src/app/payment ./payment
CMD [ "./payment" ]