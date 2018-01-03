FROM golang

COPY ./app /go/src/github.com/user/myProject/app
WORKDIR /go/src/github.com/user/myProject/app

RUN go get github.com/codegangsta/gin
RUN go-wrapper download
RUN go-wrapper install

EXPOSE 3000
EXPOSE 3001
