FROM golang

ARG app_env
ENV APP_ENV $app_env

COPY ./app /go/src/github.com/user/myProject/app
WORKDIR /go/src/github.com/user/myProject/app

RUN go get ./
RUN go build

ENV PORT 80
ENV TELE_KEY <insert_key>

EXPOSE 80

CMD app
