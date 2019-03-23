FROM golang:latest

WORKDIR /home

COPY . .
RUN go install -mod=vendor .

ENTRYPOINT ["kujo"]
