FROM library/golang

RUN go get github.com/beego/bee

ADD just-tit .
ADD conf/ conf/
ADD static/ static/
ADD views/ views/

# Use the revel CLI to start up our application.
ENTRYPOINT ./just-tit

# Open up the port where the app is running.
EXPOSE 8080
