FROM library/golang

RUN go get github.com/astaxie/beego
RUN go get github.com/astaxie/beego/cache
RUN go get github.com/astaxie/beego/cache/redis

# Recompile the standard library without CGO
RUN CGO_ENABLED=0 go install -a std

ENV APP_DIR $GOPATH/src/github.com/dsmatilla/just-tit
RUN mkdir -p $APP_DIR

# Set the entrypoint
ENTRYPOINT (cd $APP_DIR && ./just-tit)
ADD . $APP_DIR

# Compile the binary and statically link
RUN cd $APP_DIR && CGO_ENABLED=0 go build -ldflags '-w -s'

EXPOSE 8080
