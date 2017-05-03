FROM golang:1.8

ADD src /go/src
RUN cd /go/src/meanbee.com/slack/jira-bot/ && go install

ENTRYPOINT /go/bin/jira-bot
