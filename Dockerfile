FROM golang:1.7

RUN curl -sSL https://get.docker.com | bash


RUN go get -u github.com/govend/govend
COPY . /go/src/github.com/sevoma/SeriousApiarist
WORKDIR /go/src/github.com/sevoma/SeriousApiarist
RUN govend -v
RUN go install
RUN rm -rf /go/src/github.com/sevoma/SeriousApiarist
COPY config.yaml /

#RUN adduser -D -u 59999 -s /usr/sbin/nologin user
#USER 59999

CMD ["SeriousApiarist"]
