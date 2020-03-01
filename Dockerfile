FROM golang:1.13 as builder
#RUN apt install gcc g++
WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN ./run.sh && go build -o worker .

FROM debian:buster
RUN apt update && apt install -y curl nmap traceroute
WORKDIR /worker
COPY --from=builder /app/ .
RUN find . -name '*.go' -exec rm '{}' \;
CMD ["/worker/worker"]
EXPOSE 8080
