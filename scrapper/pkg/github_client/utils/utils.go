package utils

import (
	"errors"
	"regexp"
	"strings"
)


var (
	ErrPrefixNotValid = errors.New("link does not start with the prefix `https://github.com/`")
	ErrFailedToSplitTheLink = errors.New("failed to split the link")
)


func GetLinkParams(link string) (owner string, repoName string, err error) {
	link = strings.TrimSpace(link)
	newLink := strings.TrimPrefix(link, "https://github.com/")
	newLink = strings.TrimPrefix(newLink, "https://www.github.com/")

	if newLink == link {
		return "", "", ErrPrefixNotValid
	}

	newLink = strings.TrimSuffix(newLink, ".git")
	content := strings.Split(newLink, "/")

	if len(content) < 2 {
		return "", "", ErrFailedToSplitTheLink
	}

	owner, repoName = content[0], content[1]

	if owner == "" || repoName == "" {
		return "", "", ErrFailedToSplitTheLink
	}

	return owner, repoName, nil
}

func IsGitHubRepo(url string) bool {
	re := `^https:\/\/github\.com\/[A-Za-z0-9_-]+\/[A-Za-z0-9_-]+$`

	match, err := regexp.MatchString(re, url)
	if err != nil {
		return false
	}

	return match
}