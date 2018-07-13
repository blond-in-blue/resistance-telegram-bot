FROM golang:1.10.1

RUN apt-get -y update && apt-get install -y wget nano git build-essential yasm pkg-config

# Compile and install ffmpeg from source
RUN git clone https://github.com/FFmpeg/FFmpeg /root/ffmpeg && \
    cd /root/ffmpeg && \
    ./configure --enable-nonfree --disable-shared --extra-cflags=-I/usr/local/include && \
    make -j8 && make install -j8

# If you want to add some content to this image because the above takes a LONGGG time to build
ARG CACHEBREAK=1

# ENV MAGICK_URL "https://www.imagemagick.org/download"
# ENV MAGICK_VERSION 7.0.8-5

# RUN apt-get update -y \
#   && apt-get install -y xz-utils  \
#   && apt-get install -y --no-install-recommends \
#     libpng-dev \
#     libjpeg-dev libtiff-dev \
#   && cd /tmp \
#   && wget "${MAGICK_URL}/ImageMagick-${MAGICK_VERSION}.tar.xz" \
#   && tar xvf "/tmp/ImageMagick-${MAGICK_VERSION}.tar.xz" \
#   && cd "ImageMagick-${MAGICK_VERSION}" \
#   && ./configure --with-png=yes \
#   && make \
#   && make install \
#   && ldconfig /usr/local/lib

COPY ./app /go/src/github.com/user/myProject/app
WORKDIR /go/src/github.com/user/myProject/app

RUN go get github.com/codegangsta/gin
RUN go get -d -v ./...
RUN go install -v ./...

EXPOSE 3000
EXPOSE 3001
