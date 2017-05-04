FROM golang:1.8

ADD src /go/src
RUN cd /go/src/github.com/meanbee/slack-jira-bot/ && go install

ENTRYPOINT /go/bin/slack-jira-bot
