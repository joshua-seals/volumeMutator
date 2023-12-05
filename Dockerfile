FROM golang:1.21 as builder 
ENV CGO_ENABLED 0
ARG CERT_PATH

COPY . /helx 

# Build the Binary, passing in CERT_PATH from the Makefile
WORKDIR /helx/tools 
RUN go build -ldflags "-X main.certPath=${CERT_PATH}" -o generateTLSCerts

# WORKDIR /helx
# RUN go build -o volumeMutator


FROM alpine:3.18
# Keep these ARGS in the final image
ARG BUILD_DATE
ARG BUILD_REF

# Ensure we have a valid user and group
RUN addgroup -g 1000 -S helx && \
    adduser -u 1000 -h /helx -G helx -S helx

# Copy tooling for our initContainer
COPY --from=builder --chown=helx:helx /helx/tools/generateTLSCerts /helx
# Copy main application
# COPY --from=builder --chown=helx:helx /helx/volumeMutator /helx

USER helx

WORKDIR /helx
EXPOSE 8443
CMD ["./generateTLSCerts"]

LABEL org.opencontainers.image.created="${BUILD_DATE}" \
      org.opencontainers.image.title="volumeMutator" \
      org.opencontainers.image.authors="Joshua Seals" \
      org.opencontainers.image.source="https://github.com/helxplatform/volumeMutator" \
      org.opencontainers.image.revision="${BUILD_REF}" \
      org.opencontainers.image.vendor="RENCI - Renaissance Computing Institute"