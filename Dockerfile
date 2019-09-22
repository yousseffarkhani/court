FROM golang

ARG app_env
ENV APP_ENV $app_env

WORKDIR /go/src/github.com/yousseffarkhani/court
ADD . .

# Downloads all dependecies
RUN go get ./
RUN go install

# Launches app if production mode otherwise launches a server with hot reload
CMD if [ "${APP_ENV}" = "production" ]; then app; \
	else go get github.com/gravityblast/fresh && fresh; \
	fi

EXPOSE 8080