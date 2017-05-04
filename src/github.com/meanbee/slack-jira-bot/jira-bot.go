package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"errors"

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
	_, err := getConfig()

	if err != nil {
		log.Printf("Error extracting configuration: %v", err)
		os.Exit(1)
	}

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	log.Printf("main: Now listening for events")

	for {
		select {
		case msg := <-rtm.IncomingEvents:
			switch ev := msg.Data.(type) {
			case *slack.MessageEvent:
				handleIncomingMessage(ev.Msg)
			case *slack.LatencyReport:
				log.Printf("main: Current latency: %v\n", ev.Value)
			case *slack.RTMError:
				log.Printf("main: Error: %s\n", ev.Error())
			case *slack.InvalidAuthEvent:
				log.Printf("main: Invalid Slack credentials")
				os.Exit(2)
			default:
			// Ignore other events..
			}
		}
	}
}

func handleIncomingMessage(message slack.Msg) {
	messageText := message.Text

	if shouldIgnoreMessage(message) {
		log.Print("handleMessage: Ignoring message")
		return
	}

	matches := extractIssueIDs(messageText)

	for i := 0; i < len(matches); i++ {
		issueID := matches[i]
		log.Printf("handleMessage: Identified " + issueID + " in message")

		respondToIssueMentioned(message.Channel, issueID)
	}
}

func respondToIssueMentioned(channel string, issueID string) {
	defer func() {
		if e := recover(); e != nil {
			log.Printf("Exception responding to issue %s: %v", issueID, e)
		}
	}()

	api := getSlackAPI()
	config, _ := getConfig()

	issueData := getJiraIssue(issueID)

	params := slack.PostMessageParameters{
		Username: config.Username,
		Markdown: true,
	}

	var assignee = ""

	if issueData.Fields.Assignee != nil {
		assignee = issueData.Fields.Assignee.DisplayName
	} else {
		assignee = "Nobody"
	}

	attachment := slack.Attachment{
		Pretext: formatMessage(issueData),
		MarkdownIn: []string{"pretext", "fields"},
		Fields: []slack.AttachmentField{
			slack.AttachmentField{
				Title: "Status",
				Value: issueData.Fields.Status.Name,
				Short: true,
			},
			slack.AttachmentField{
				Title: "Assigned",
				Value: assignee,
				Short: true,
			},
		},
	}

	params.Attachments= []slack.Attachment{attachment}



	api.PostMessage(channel, "", params)
}

func getSlackAPI() *slack.Client {
	config, _ := getConfig()

	return slack.New(config.SlackAPIKey)
}

func getChannel(channelID string) (*slack.Channel, error) {
	api := getSlackAPI()

	return api.GetChannelInfo(channelID)
}

func formatMessage(issue gojira.Issue) string {
	message := fmt.Sprintf(
		"*%s*: %s - %s",
		issue.Key,
		issue.Fields.Summary,
		getJiraURL(issue.Key),
	)

	return message
}

func getJiraURL(issueKey string) string {
	config, _ := getConfig()

	return config.JiraBaseURL + "/browse/" + issueKey
}

func getJiraIssue(issueID string) gojira.Issue {
	jiraAPIPath := "/rest/api/2"
	jiraActivityPath := "/activity"
	config, _ := getConfig()

	jira := gojira.NewJira(
		config.JiraBaseURL,
		jiraAPIPath,
		jiraActivityPath,
		&gojira.Auth{
			Login:    config.JiraUsername,
			Password: config.JiraPassword,
		},
	)

	issueData := jira.Issue(issueID)

	log.Printf("issueData: %v", issueData)

	return issueData
}

func shouldIgnoreMessage(message slack.Msg) bool {
	config, _ := getConfig()

	return message.Username == config.Username || message.SubType == "bot_message"
}

func extractIssueIDs(message string) []string {
	re := regexp.MustCompile(`\b(\w+)-(\d+)\b`)
	matches := re.FindAllString(message, -1)

	// @see http://www.dotnetperls.com/remove-duplicates-slice
	encountered := map[string]bool{}
	result := []string{}

	for v := range matches {
		// convert all match to upper case.
		matches[v] = strings.ToUpper(matches[v])
		if encountered[matches[v]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[matches[v]] = true
			// Append to result slice.
			result = append(result, matches[v])
		}
	}
	// Return the new slice.
	return result
}

func getConfig() (BotConfig, error) {

	apiKey := os.Getenv("SLACK_API_KEY")
	jiraBaseUrl := os.Getenv("JIRA_BASEURL")
	jiraUsername := os.Getenv("JIRA_USERNAME")
	jiraPassword := os.Getenv("JIRA_PASSWORD")

	if apiKey == "" {
		return BotConfig{}, errors.New("Expected API Key in SLACK_API_KEY environment variable")
	}

	if jiraBaseUrl == "" {
		return BotConfig{}, errors.New("Expected the base URL of your Jira installation in JIRA_BASEURL environment variable")
	}

	if jiraUsername == "" {
		return BotConfig{}, errors.New("Expected a username to access your Jira installation in JIRA_USERNAME environment variable")
	}

	if jiraUsername == "" {
		return BotConfig{}, errors.New("Expected a password to access your Jira installation in JIRA_PASSWORD environment variable")
	}

	return BotConfig{
		Username:     "jirabot",
		SlackAPIKey:  apiKey,
		JiraBaseURL:  jiraBaseUrl,
		JiraUsername: jiraUsername,
		JiraPassword: jiraPassword,
	}, nil
}
