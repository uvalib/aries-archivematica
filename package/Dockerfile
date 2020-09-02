#
# build the target application
#
FROM golang:1.15 as builder
RUN apt-get -y update

WORKDIR /build
COPY go.mod go.sum Makefile ./
COPY cmd ./cmd
RUN make build-linux

#
# build the target container
#
FROM alpine:3.12

# update the packages
RUN apk update && apk upgrade && apk add bash tzdata ca-certificates curl && rm -fr /var/cache/apk/*

# Create the run user and group
RUN addgroup webservice && adduser webservice -G webservice -D

# set the timezone appropriatly
ENV TZ=UTC
RUN cp /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

# Specify home 
ENV APP_HOME /aries-archivematica
WORKDIR $APP_HOME

# Create necessary directories
RUN mkdir -p $APP_HOME/bin
RUN chown -R webservice $APP_HOME && chgrp -R webservice $APP_HOME

# port and run command
EXPOSE 8088
CMD scripts/entry.sh

# Move in necessary assets
COPY package/data/container_bash_profile /home/webservice/.profile
COPY package/scripts/entry.sh $APP_HOME/scripts/entry.sh
COPY --from=builder /build/bin/aries-archivematica $APP_HOME/bin/aries-archivematica

# Ensure permissions are correct
RUN chown webservice:webservice /home/webservice/.profile $APP_HOME/scripts/entry.sh $APP_HOME/bin/aries-archivematica && chmod 755 /home/webservice/.profile $APP_HOME/scripts/entry.sh $APP_HOME/bin/aries-archivematica

# Add the build tag
COPY buildtag.* $APP_HOME/

# Specify the user
USER webservice

#
# end of file
#
