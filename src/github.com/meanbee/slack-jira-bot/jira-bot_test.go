package main

import (
	"testing"

	"github.com/nlopes/slack"
)

func TestExtractIssueID(t *testing.T) {
	result := extractIssueIDs("ABC-123")

	if len(result) != 1 {
		t.Errorf("Expected one result, got %v", len(result))
	}

	if result[0] != "ABC-123" {
		t.Errorf("Expected ABC-123, got %v", result[0])
	}
}

func TestExtractMultipleIssueIDs(t *testing.T) {
	result := extractIssueIDs("ABC-123 DEF-345")

	if len(result) != 2 {
		t.Errorf("Expected two result, got %v", len(result))
	}

	if result[0] != "ABC-123" {
		t.Errorf("Expected ABC-123, got %v", result[0])
	}

	if result[1] != "DEF-345" {
		t.Errorf("Expected DEF-345, got %v", result[0])
	}
}

func TestExtractIssueIDUniques(t *testing.T) {
	result := extractIssueIDs("ABC-123 ABC-123 ABC-123 ABC-123 ABC-123 ABC-123 abc-123 aBc-123")

	if len(result) != 1 {
		t.Errorf("Expected one result, got %v", len(result))
	}

	if result[0] != "ABC-123" {
		t.Errorf("Expected ABC-123, got %v", result[0])
	}
}

func TestIgnoresMessageFromBot(t *testing.T) {
	messageBot := slack.Msg{
		SubType: "bot_message",
	}

	if !shouldIgnoreMessage(messageBot) {
		t.Errorf("Message was from a bot, expected to ignore")
	}
}

func TestDoesNotIgnoreMessageFromUser(t *testing.T) {
	messageUser := slack.Msg{
		Username: "Eric",
		SubType: "",
	}

	if shouldIgnoreMessage(messageUser) {
		t.Errorf("Message was from a user, expected to not ignore")
	}
}
