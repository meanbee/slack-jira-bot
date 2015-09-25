FROM golang

ADD . /go/src/meanbee.com/slack/jira-bot/

RUN go get github.com/nlopes/slack
RUN go get github.com/plouc/go-jira-client

RUN cd /go/src/meanbee.com/slack/jira-bot/ && go install

ENTRYPOINT /go/bin/jira-bot
