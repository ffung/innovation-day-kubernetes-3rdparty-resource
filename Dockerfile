FROM scratch
ADD environment-manager /environment-manager

ENTRYPOINT ["/environment-manager"]

