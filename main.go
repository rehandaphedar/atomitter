package main

import (
	"bytes"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/mmcdole/gofeed"
	"github.com/spf13/viper"
)

func main() {
	initialiseConfig()
	tweetTemplate := getTweetTemplate()

	articles := getArticles()
	client := getTwitterClient()
	tweets := getTweets(client)
	untweetedArticles := getUntweetedArticles(articles, tweets, tweetTemplate)

	tweetArticles(client, untweetedArticles, tweetTemplate)
}

func initialiseConfig() {
	viper.SetDefault("consumer_key", "CONSUMER_KEY")
	viper.SetDefault("consumer_secret", "CONSUMER_SECRET")
	viper.SetDefault("token", "TOKEN")
	viper.SetDefault("token_secret", "TOKEN_SECRET")
	viper.SetDefault("feed_url", "https://example.com/feed.atom")
	viper.SetDefault("username", "USERNAME")
	viper.SetDefault("format", "Read my post \"{{ .Title }}\": {{ .Link }}")

	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		log.Fatalf("Error reading user config directory: %v", err)
	}
	configDir := filepath.Join(userConfigDir, "atomitter")
	os.Mkdir(configDir, os.ModePerm)

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(configDir)

	err = viper.SafeWriteConfig()
	if err != nil {
		_err := *new(viper.ConfigFileAlreadyExistsError)
		if !errors.As(err, &_err) {
			log.Fatalf("Error writing config file: %v", err)
		}
	}

	err = viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	err = viper.WriteConfig()
	if err != nil {
		log.Fatalf("Error writing config file: %v", err)
	}
}

func getTweetTemplate() (tweetTemplate *template.Template) {
	tweetTemplate, err := template.New("main").Parse(viper.GetString("format"))
	if err != nil {
		log.Fatalf("Error initialising tweet template: %v", err)
	}
	return
}

func getArticles() (articles []*gofeed.Item) {
	fp := gofeed.NewParser()
	feed, _ := fp.ParseURL(viper.GetString("feed_url"))
	articles = feed.Items
	return
}

func getTwitterClient() (client *twitter.Client) {
	var config = oauth1.NewConfig(viper.GetString("consumer_key"), viper.GetString("consumer_secret"))
	var token = oauth1.NewToken(viper.GetString("token"), viper.GetString("token_secret"))
	var httpClient = config.Client(oauth1.NoContext, token)
	client = twitter.NewClient(httpClient)
	return
}

func getTweets(client *twitter.Client) (tweets []twitter.Tweet) {
	trimUser := true
	excludeReplies := true
	includeRetweets := false
	tweets, _, err := client.Timelines.UserTimeline(&twitter.UserTimelineParams{
		ScreenName:      viper.GetString("username"),
		TrimUser:        &trimUser,
		ExcludeReplies:  &excludeReplies,
		IncludeRetweets: &includeRetweets,
	})
	if err != nil {
		log.Fatalf("Error retrieving tweets: %v", err)
	}
	return
}

func getUntweetedArticles(articles []*gofeed.Item, tweets []twitter.Tweet, tweetTemplate *template.Template) (untweetedArticles []*gofeed.Item) {
	untweetedArticles = articles

	for _, tweet := range tweets {
		tweetText := getTweetTextWithFullURLs(tweet)

		for _, article := range articles {
			articleText := getArticleText(article, tweetTemplate)

			if articleText == tweetText {
				untweetedArticles = filter(untweetedArticles, func(cur *gofeed.Item) bool {
					return cur.Link != article.Link
				})
			}
		}
	}
	return
}

func getTweetTextWithFullURLs(tweet twitter.Tweet) (tweetText string) {
	tweetText = tweet.Text
	for _, link := range tweet.Entities.Urls {
		tweetText = strings.ReplaceAll(tweetText, link.URL, link.ExpandedURL)
	}
	return tweetText
}

func getArticleText(article *gofeed.Item, tweetTemplate *template.Template) (articleText string) {
	var b bytes.Buffer
	err := tweetTemplate.Execute(&b, article)
	if err != nil {
		log.Fatalf("Error applying tweet template: %v", err)
	}
	articleText = b.String()
	return
}

func tweetArticles(client *twitter.Client, articles []*gofeed.Item, tweetTemplate *template.Template) {
	for _, article := range articles {
		articleText := getArticleText(article, tweetTemplate)
		log.Printf("Tweeting: %s", articleText)

		_, _, err := client.Statuses.Update(articleText, &twitter.StatusUpdateParams{})
		if err != nil {
			log.Fatalf("Error tweeting: %v", err)
		}
	}
}

func filter(source []*gofeed.Item, predicate func(*gofeed.Item) bool) (filtered []*gofeed.Item) {
	for _, s := range source {
		if predicate(s) {
			filtered = append(filtered, s)
		}
	}
	return
}
