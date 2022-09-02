# (Deprecated) Slack Jira Issue Expansion

_Notice:_ This project is no longer being maintained and the APIs that is uses are [deprecated](https://github.com/meanbee/slack-jira-bot/issues/12).

[![Build Status](https://travis-ci.org/meanbee/slack-jira-bot.svg?branch=travis)](https://travis-ci.org/meanbee/slack-jira-bot)

> **Update 2017-08-24:** If you're looking to get setup quickly then you should check out https://slack.atlassian.io before using this image.  If that's not for you - come back!

A killer feature of the integration between Hipchat and Jira is the issue expansion.  Whenever a Jira issue is mention in the chat, the Jira integration would pop up some high level information about the issue and a link.

![](https://punkstar.keybase.pub/github/slack-jira-bot/screenshots/hipchat_example.png)

With the Jira integration in Slack being a bit light, we decided to implement a simple bot using the [Slack RTM API](https://api.slack.com/rtm):

![](https://punkstar.keybase.pub/github/slack-jira-bot/screenshots/message_example_with_attachment.png)

# Installation
    
## Docker Image

    docker run -it --restart=always -d \
        -e JIRA_BASEURL=https://yourjirainstall.atlassian.net \
        -e JIRA_USERNAME='yourjirausername' \
        -e JIRA_PASSWORD='yourjirapassword' \
        -e SLACK_API_KEY='yourslackapikey' \
        meanbee/slack-jira-bot:latest
    
# Configuration

The configuration is run of environment variables:

* `SLACK_API_KEY`
* `JIRA_BASEURL`, e.g. `https://yourcompany.atlassian.net`
* `JIRA_USERNAME`
* `JIRA_PASSWORD`
