FROM gobuffalo/buffalo:v0.12.3 as builder

RUN mkdir -p $GOPATH/src/github.com/joshgav/vanpool-manager/services/frontend
WORKDIR $GOPATH/src/github.com/joshgav/vanpool-manager/services/frontend

ADD package.json .
ADD yarn.lock .
RUN yarn install --no-progress
ADD . .
RUN dep ensure
RUN buffalo build --static -o /bin/app

FROM alpine
RUN apk add --no-cache bash
RUN apk add --no-cache ca-certificates

WORKDIR /bin/

COPY --from=builder /bin/app .

ENV GO_ENV=development
# ENV GO_ENV=production

ENV ADDR=0.0.0.0

EXPOSE 3000

# CMD /bin/app migrate; /bin/app
CMD exec /bin/app
