FROM gcr.io/distroless/static
COPY markscribe /usr/local/bin/markscribe
ENTRYPOINT [ "/usr/local/bin/markscribe" ]
