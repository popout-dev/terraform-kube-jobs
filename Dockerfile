FROM debian:bullseye-slim

WORKDIR /terraform

RUN apt-get update -y && \
apt-get upgrade -y && \
apt-get install -y ca-certificates

COPY ./terra-kube-jobs .

ENTRYPOINT [ "./terra-kube-jobs" ]