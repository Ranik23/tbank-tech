package utils

import (
	"errors"
	"testing"
)

func TestGetTheLinkParamsSucces(t *testing.T) {
	url := "https://github.com/google/go-github"

	owner, repoName, err := GetLinkParams(url)
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	if owner != "google" {
		t.Fatalf("wrong username")
	}

	if repoName != "go-github" {
		t.Fatalf("wrong repo name")
	}
}


func TestGetTheLinkParamsFail(t *testing.T) {
	url := "https://github/google/go-github"

	owner, repoName, err := GetLinkParams(url)
	if err != nil {
		if !errors.Is(err, ErrPrefixNotValid) {
			t.Fatalf("wrong type of error occr")
		}
	}

	if owner != "google" {
		t.Fatalf("owner name is not valid")
	}

	if repoName != "go-github" {
		t.Fatalf("repo name is not valid")
	}
}

func TestIsGitHubRepoSucces(t *testing.T) {

	testURLS := []string{
		"https://github.com/MAtveyka12/tbank",
		"https://github.com/OpenAPITools/openapi-generator",
		"https://github.com/epchamp001/avito-tech-merch",
	}

	for _, testURL := range testURLS {
		if !IsGitHubRepo(testURL) {
			t.Fatalf("%s is a repo link, but functions says no", testURL)
		}
	}
}

func TestIsGitHubRepoFail(t *testing.T) {

	testURLS := []string{
		"https://github.com/MAtveyka12/",
		"https://github.com/OpenAPIToo",
		"https:/pchamp001/avito-tech-merch",
		"https://github.com/MAtveyka12/tbank/tree/test/.github/workflows",
	}

	for _, testURL := range testURLS {
		if IsGitHubRepo(testURL) {
			t.Fatalf("%s is not repo link, but functions says yes", testURL)
		}
	}
}

