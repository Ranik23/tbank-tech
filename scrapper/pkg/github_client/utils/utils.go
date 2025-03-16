package utils

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)


var (
	ErrPrefixNotValid = errors.New("link does not start with the prefix `https://github.com/`")
	ErrFailedToSplitTheLink = errors.New("failed to split the link")
)


func GetLinkParams(link string) (owner string, repoName string, err error) {

	newLink := strings.TrimPrefix(link, "https://github.com/")

	if newLink == link {
		return "", "", ErrPrefixNotValid
	}

	fmt.Println(newLink)

	content := strings.Split(newLink, "/")

	if len(content) < 2 {
		return "", "", ErrFailedToSplitTheLink
	}

	user, repoName := content[0], content[1]

	return user, repoName, nil
}

func IsGitHubRepo(url string) bool {
	re := `^https:\/\/github\.com\/[A-Za-z0-9_-]+\/[A-Za-z0-9_-]+$`

	match, err := regexp.MatchString(re, url)
	if err != nil {
		return false
	}

	return match
}