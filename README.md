# Introduction

![Atomitter](https://git.sr.ht/~rehandaphedar/atomitter/blob/main/assets/logo.png)

Atomitter is a Twitter bot that syncs an Atom/RSS/JSON feed to a twitter account.

# Installation

## Dependencies

- `go`

## Building

Just build using standard the golang toolchain:


```shell
git clone https://git.sr.ht/~rehandaphedar/atomitter
cd atomitter
go build .
```

## Usage

Obtain the config file and edit it:

```shell
mkdir -p ~/.config/atomitter
wget https://git.sr.ht/~rehandaphedar/atomitter/blob/main/config.example.yaml -O ~/.config/atomitter/config.yaml
```

Then, simply run `atomitter` to sync.

# Configuration

The config file is stored in `$XDG_CONFIG_HOME/atomitter/config.yaml`.

- `consumer_key`: Your Twitter Consumer Key
- `consumer_secret`: Your Twitter Consumer Secret
- `token`: Your Twitter Token
- `token_secret`: Your Twitter Token Secret
- `feed_url`: Link to the feed to sync
- `username`: Twitter username of the account you wish to sync to
- `format`: String describing how to format your tweet

Also see [config.example.yaml](https://git.sr.ht/~rehandaphedar/atomitter/tree/main/config.example.yaml).

# Auto Sync

You can use cronjobs to automatically sync your feed after a certain time. For example, to sync your feed once every day:

```
0 0 * * * atomitter
```

# Limitations

Currently, Atomitter doesn't work properly for Twitter accounts with more than 200 tweets.
