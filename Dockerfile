FROM golang:alpine
# install git
RUN apk --update add \
	git openssl \
	&& rm /var/cache/apk/*

WORKDIR /go/src/chaos-stack-tesco
COPY ./ /go/src/chaos-stack-tesco

ADD . /go/src


#ENV GO111MODULE=on
RUN go get -u github.com/golang/dep/cmd/dep
RUN dep ensure

#RUN go build -o main .
#
#CMD ["./main"]
CMD ["go", "run", "main.go"]

EXPOSE 5500

#RUN go get -u github.com/golang/dep/cmd/dep
#
#ADD ./main.go /go/src/app
#COPY ./Gopkg.toml /go/src/app
#
#WORKDIR /go/src/app
#
#RUN dep ensure
#RUN go test -v
#RUN go build
#
#CMD ["./app"]