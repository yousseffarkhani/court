FROM golang

ARG app_env
ENV APP_ENV $app_env

WORKDIR /go/src/app
ADD . .

# Adds local dependencies
# courtdb package
RUN mkdir /go/src/courtdb
RUN ln -s /go/src/app/courtdb/courtdb.go /go/src/courtdb/courtdb.go
# server package
RUN mkdir /go/src/server
RUN ln -s /go/src/app/server/server.go /go/src/server/server.go
# views package
RUN mkdir /go/src/views
RUN ln -s /go/src/app/views/views.go /go/src/views/views.go
# model package
RUN mkdir /go/src/model
RUN ln -s /go/src/app/model/model.go /go/src/model/model.go
# middlewares package
RUN mkdir /go/src/middlewares
RUN ln -s /go/src/app/middlewares/middlewares.go /go/src/middlewares/middlewares.go
# contactMail package
RUN mkdir /go/src/contactMail
RUN ln -s /go/src/app/contactMail/contactMail.go /go/src/contactMail/contactMail.go
# handlers package
RUN mkdir /go/src/handlers
RUN ln -s /go/src/app/handlers/handlers.go /go/src/handlers/handlers.go

# Downloads all dependecies
RUN go get ./
RUN go install

CMD if [ "${APP_ENV}" = "production" ]; then app; \
	else go get github.com/gravityblast/fresh && fresh; \
	fi
    
EXPOSE 8080