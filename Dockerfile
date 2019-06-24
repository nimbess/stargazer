#FROM golang:alpine
FROM scratch

ADD bin/stargazer /usr/bin/stargazer
ADD stargazer.yaml /etc/stargazer.yaml

ENTRYPOINT ["/bin/stargazer -v"]
#ENTRYPOINT ["/bin/sh"]
#CMD ["stargazer", "-config-path", "/etc"]
