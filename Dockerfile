FROM golang:alpine

ADD bin/stargazer /usr/bin/stargazer
ADD stargazer.yaml /etc/stargazer.yaml

CMD ["stargazer", "-config-path", "/etc"]
