FROM debian:bullseye-slim

WORKDIR /terraform

COPY ./terra-kube-jobs .

ENTRYPOINT [ "terra-kube-jobs" ]