FROM alpine:3.17.0

ENV UBI_EMAIL="my@mail.com"
ENV UBI_PASSWORD="v3rys4fep4s5w0rd"
ENV UBI_OBSERVED_USERNAMES="UbiName1,UbiName2"
EXPOSE 2112

COPY app /
ENTRYPOINT ["/app"]
