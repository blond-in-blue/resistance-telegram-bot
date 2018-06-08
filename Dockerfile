FROM golang:1.9

ENV MAGICK_URL "https://www.imagemagick.org/download"
ENV MAGICK_VERSION 7.0.7-38

RUN apt-get update -y \
  && apt-get install -y xz-utils  \
  && apt-get install -y --no-install-recommends \
    libpng-dev libjpeg-dev libtiff-dev \
  && cd /tmp \
  && wget "${MAGICK_URL}/ImageMagick-${MAGICK_VERSION}.tar.xz" \
  && tar xvf "/tmp/ImageMagick-${MAGICK_VERSION}.tar.xz" \
  && cd "ImageMagick-${MAGICK_VERSION}" \
  && ./configure --with-png=yes \
  && make \
  && make install \
  && ldconfig /usr/local/lib

COPY ./app /go/src/github.com/user/myProject/app
WORKDIR /go/src/github.com/user/myProject/app

RUN go get github.com/codegangsta/gin
RUN go-wrapper download
RUN go-wrapper install

EXPOSE 3000
EXPOSE 3001
