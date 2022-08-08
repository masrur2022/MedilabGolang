FROM golang
COPY . /new
WORKDIR /perent
COPY go.mod /perent/
COPY go.sum /perent/
RUN go mod download 
COPY . /perent/
EXPOSE 4400
RUN go build -o /main
CMD [ "/main" ]
