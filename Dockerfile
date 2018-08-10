FROM node:latest as js-builder
WORKDIR /web
COPY ./web/package.json ./package.json
RUN npm install -g npm@latest
RUN npm install
COPY ./web .
RUN npm run build

FROM golang:latest as api-builder
ENV PACKAGE_NAME=github.com/joshgav/vanpool-manager
ENV DEP_VERSION=0.5.0
RUN curl -fsSL -o /usr/local/bin/dep \
    https://github.com/golang/dep/releases/download/v${DEP_VERSION}/dep-linux-amd64 \
 && chmod +x /usr/local/bin/dep
WORKDIR /go/src/${PACKAGE_NAME}
COPY . .
RUN dep ensure -v
RUN go build -a -v -o server

FROM golang:latest
ENV PACKAGE_NAME=github.com/joshgav/vanpool-manager
ENV PORT=8080
EXPOSE ${PORT}
WORKDIR /app
COPY --from=api-builder /go/src/${PACKAGE_NAME}/server .
COPY --from=js-builder /web/dist ./web/dist
RUN chmod +x ./server
CMD ["./server"]
