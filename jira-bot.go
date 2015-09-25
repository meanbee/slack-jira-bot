package main

import (
	"fmt"
	"log"
	"os"
	"regexp"

	"github.com/nlopes/slack"
	gojira "github.com/plouc/go-jira-client"
)

// Configuration for the bot
type BotConfig struct {
	Username     string
	SlackAPIKey  string
	JiraUsername string
	JiraPassword string
	JiraBaseURL  string
}

func main() {
	api := getSlackAPI()

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	log.Printf("main: Now listening for events")

	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.MessageEvent:
				handleMessage(ev.Msg)
			case *slack.LatencyReport:
				log.Printf("main: Current latency: %v\n", ev.Value)
			case *slack.RTMError:
				log.Printf("main: Error: %s\n", ev.Error())
			case *slack.InvalidAuthEvent:
				log.Printf("main: Invalid credentials")
			default:
				// Ignore other events..
			}
		}
	}
}

func handleMessage(message slack.Msg) {
	api := getSlackAPI()

	messageFrom := message.Username
	messageText := message.Text

	if messageFrom == getConfig().Username {
		log.Print("handleMessage: Ignoring message, it was from me")
		return
	}

	re := regexp.MustCompile(`\b(\w+)-(\d+)\b`)
	matches := re.FindAllString(messageText, -1)

	params := slack.PostMessageParameters{
		Username: getConfig().Username,
		Markdown: true,
	}

	for i := 0; i < len(matches); i++ {
		issueID := matches[i]
		log.Printf("handleMessage: Identified " + issueID + " in message")

		issueData := getJiraIssue(issueID)

		api.PostMessage(message.Channel, formatMessage(issueData), params)
	}
}

func getSlackAPI() *slack.Client {
	return slack.New(getConfig().SlackAPIKey)
}

func getChannel(channelID string) (*slack.Channel, error) {
	api := getSlackAPI()

	return api.GetChannelInfo(channelID)
}

func formatMessage(issue gojira.Issue) string {
	message := fmt.Sprintf(
		"*%s: %s* _Reported by %s_ - %s",
		issue.Key,
		issue.Fields.Summary,
		issue.Fields.Reporter.DisplayName,
		getJiraURL(issue.Key),
	)

	return message
}

func getJiraURL(issueKey string) string {
	return getConfig().JiraBaseURL + "/browse/" + issueKey
}

func getJiraIssue(issueID string) gojira.Issue {
	jiraAPIPath := "/rest/api/latest"
	jiraActivityPath := ""

	jira := gojira.NewJira(
		getConfig().JiraBaseURL,
		jiraAPIPath,
		jiraActivityPath,
		&gojira.Auth{
			Login:    getConfig().JiraUsername,
			Password: getConfig().JiraPassword,
		},
	)

	issueData := jira.Issue(issueID)

	return issueData
}

func getConfig() BotConfig {
	return BotConfig{
		Username:     "jirabot",
		SlackAPIKey:  os.Getenv("SLACK_API_KEY"),
		JiraBaseURL:  os.Getenv("JIRA_BASEURL"),
		JiraUsername: os.Getenv("JIRA_USERNAME"),
		JiraPassword: os.Getenv("JIRA_PASSWORD"),
	}
}
