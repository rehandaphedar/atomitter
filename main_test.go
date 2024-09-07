package main

import (
	"os"
	"testing"
	"text/template"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/mmcdole/gofeed"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	returnCode := m.Run()
	os.Exit(returnCode)
}

func TestNormalGetArticleText(t *testing.T) {
	require := require.New(t)

	format := "Read my post \"{{ .Title }}\": {{ .Link }}"
	tweetTemplate, err := template.New("main").Parse(format)
	require.NoError(err)

	article := &gofeed.Item{Title: "Example Title", Link: "https://example.com"}

	articleText := getArticleText(article, tweetTemplate)

	require.Equal("Read my post \"Example Title\": https://example.com", articleText)
}

func TestIncompleteGetArticleText(t *testing.T) {
	require := require.New(t)

	format := "Read my post \"{{ .Title }}\": {{ .Link }}"
	tweetTemplate, err := template.New("main").Parse(format)
	require.NoError(err)

	article := &gofeed.Item{Title: "Example Title"}

	articleText := getArticleText(article, tweetTemplate)

	require.NotEqual("Read my post \"Example Title\": https://example.com", articleText)
}

func TestIncorrectGetArticleText(t *testing.T) {
	require := require.New(t)

	format := "Read my posts \"{{ .Title }}\": {{ .Link }}"
	tweetTemplate, err := template.New("main").Parse(format)
	require.NoError(err)

	article := &gofeed.Item{Title: "Example Title", Link: "https://example.com"}

	articleText := getArticleText(article, tweetTemplate)

	require.NotEqual("Read my post \"Example Title\": https://example.com", articleText)
}

func TestGetTweetTextWithFullURLs(t *testing.T) {
	require := require.New(t)

	tweet := twitter.Tweet{
		Text: "Read my post \"Example Title\": https://t.co/mtXLLfYOYE",
		Entities: &twitter.Entities{
			Urls: []twitter.URLEntity{
				{
					URL:         "https://t.co/mtXLLfYOYE",
					ExpandedURL: "https://www.bbc.co.uk/news/blogs-trending-47975564",
				},
			},
		},
	}

	tweetText := getTweetTextWithFullURLs(tweet)

	require.Equal("Read my post \"Example Title\": https://www.bbc.co.uk/news/blogs-trending-47975564", tweetText)
}

func TestGetUntweetedArticles(t *testing.T) {
	require := require.New(t)

	format := "Read my post \"{{ .Title }}\": {{ .Link }}"
	tweetTemplate, err := template.New("main").Parse(format)
	require.NoError(err)

	articles := []*gofeed.Item{
		{
			Title: "Title 1",
			Link:  "https://example.com/1",
		},
		{
			Title: "Title 2",
			Link:  "https://example.com/2",
		},
		{
			Title: "Title 3",
			Link:  "https://example.com/3",
		},
	}

	tweets := []twitter.Tweet{
		{
			Text: "Read my post \"Title 1\": https://example.com/1",
			Entities: &twitter.Entities{
				Urls: []twitter.URLEntity{},
			},
		},
	}

	expectedUntweetedArticles := []*gofeed.Item{
		{
			Title: "Title 2",
			Link:  "https://example.com/2",
		},
		{
			Title: "Title 3",
			Link:  "https://example.com/3",
		},
	}

	untweetedArticles := getUntweetedArticles(articles, tweets, tweetTemplate)

	require.Equal(expectedUntweetedArticles, untweetedArticles)

}
