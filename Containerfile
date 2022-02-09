########### ObalkyKnih API availability checker
FROM docker.io/golang:1.17-alpine AS builder
ARG BUILD_DIR=/build-dir
WORKDIR ${BUILD_DIR}/
#RUN #apk add --no-cache git && \
#    #git clone https://github.com/VerosK/obalkyknih-checker.git && \
ADD main.go ${BUILD_DIR}/main.go
RUN go build -o obalkyknih-checker main.go


FROM docker.io/alpine:3.15
ARG BUILD_DIR=/build-dir
ARG COMMIT_ID=latest
LABEL cz.knihovny.version_id=${VERSION_ID}
COPY --from=builder ${BUILD_DIR}/obalkyknih-checker /obalkyknih-checker
ENTRYPOINT ["/obalkyknih-checker", "-listenAddress=:80"]

