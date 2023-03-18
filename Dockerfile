FROM golang@sha256:48f336ef8366b9d6246293e3047259d0f614ee167db1869bdbc343d6e09aed8a AS builder

WORKDIR /go/src/github.com/ukane-philemon/go-all-the-way

COPY . .

RUN go install

EXPOSE 8080

CMD ["go-app"]
