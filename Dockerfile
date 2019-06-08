FROM golang as dev

# dep
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

# csc tool
RUN go get github.com/rexray/gocsi
RUN make -C $GOPATH/src/github.com/rexray/gocsi

RUN mkdir -p $GOPATH/src/github.com/mingfang
COPY . $GOPATH/src/github.com/mingfang/alluxio-csi

FROM dev as build
RUN cd $GOPATH/src/github.com/mingfang/alluxio-csi && \
    CGO_ENABLED=0 go build -o /usr/local/bin/alluxio-csi

FROM alluxio/alluxio-fuse:2.0.0-RC3 as final
COPY --from=build /usr/local/bin/alluxio-csi /usr/local/bin/
COPY --from=build /go/bin/csc /usr/local/bin/
