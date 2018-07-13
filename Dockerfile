FROM golang:1.10.1

RUN apt-get -y update && apt-get install -y \ 
    wget \
    nano \
    git \
    build-essential \
    yasm \
    pkg-config \
    libav-tools \
    libx264-dev \
    libpng-dev \
    libjpeg-dev \
    libtiff-dev

# Compile and install ffmpeg from source
RUN git clone https://github.com/FFmpeg/FFmpeg /root/ffmpeg && \
    cd /root/ffmpeg && \
    ./configure --enable-nonfree --disable-shared --extra-cflags=-I/usr/local/include --enable-gpl --enable-libx264 && \
    make -j8 && make install -j8
# If you want to add some content to this image because the above takes a LONGGG time to build
ARG CACHEBREAK=1

COPY ./app /go/src/github.com/user/myProject/app
WORKDIR /go/src/github.com/user/myProject/app

RUN go get github.com/codegangsta/gin
RUN go get -d -v ./...
RUN go install -v ./...

EXPOSE 3000
EXPOSE 3001
