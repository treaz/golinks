<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->
**Table of Contents**  *generated with [DocToc](https://github.com/thlorenz/doctoc)*

- [Golinks](#golinks)
  - [History of Golinks](#history-of-golinks)
  - [Why](#why)
  - [Setup](#setup)
    - [Install](#install)
    - [Database](#database)
    - [Run](#run)
    - [Run At Startup](#run-at-startup)
    - [Docker](#docker)
    - [Browser Extension Redirect Setup (Recommended)](#browser-extension-redirect-setup-recommended)
      - [Using Redirector Extension](#using-redirector-extension)
    - [DNS Setup (Manual)](#dns-setup-manual)
    - [Port Redirection Setup (Manual)](#port-redirection-setup-manual)
  - [FAQ](#faq)
  - [Troubleshooting](#troubleshooting)
  - [Developing](#developing)
  - [Roadmap](#roadmap)
  - [Contributing](#contributing)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->


# Golinks

Golinks is an internal URL shortener that organizes your company links into easily rememberable keywords. If you’re on the company network, you can type in <code>go/keyword</code> in your browser, and that will redirect you to the expanded url.


## History of Golinks

Benjamin Staffin at Google developed a golink system that introduced the "go/" domain and allowed Googlers to simply use the shortlink “go/link” in their browser. Benjamin described golinks as "AOL keywords for the corporate network."

## Why

I developed this to scratch my own itch mostly and to learn Go. It was built intending to be run locally on localhost using a sqllite database. It is meant to be lightweight and simple. I was inspired by
@thesephist's [tools](https://thesephist.com/posts/tools/) and the concept of [building software for yourself](https://changelog.com/podcast/455).
The backend API is written in Go and the frontend in Vue.js as a single page app.


## Setup

### Install

Go to the [releases](https://github.com/crhuber/golinks/releases) page and download the latest release.
Or, use my own tool: [kelp](https://github.com/crhuber/kelp)

```bash
kelp add crhuber/golinks
kelp install golinks
```

### Database

Setup a path where you want your golinks sqllite database to live and set the environment variable

```bash
mkdir ~/.golinks
export GOLINKS_DB="/Users/username/.golinks/golinks.db"
```

You can also use postgres or mysql database using a valid DSN like:

```bash
export GOLINKS_DBTYPE="mysql"
export GOLINKS_DB="user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
```

### Run

Run

```bash
golinks serve
```

Use the following flags to configure database, port and static folder

```
Flags:
  -d, --db string       DB DSN or SQLLite location path. (default "~/.golinks/golinks.db")
  -t, --dbtype string   Database type (default "sqllite")
  -h, --help            help for serve
  -p, --port int        Port to run Application server on (default 8998)
```

All the flags can also be set via environment variables

```
GOLINKS_DB
GOLINKS_DBTYPE
GOLINKS_PORT
```

### Run At Startup
To run as an Agent on boot for mac edit and copy the `io.intra.golinks.plist` file to `~/Library/LaunchAgents`  directory.
See [launchd.info](https://www.launchd.info/)

```
launchctl load ~/Library/LaunchAgents/io.intra.golinks.plist
launchctl start io.intra.golinks
tail -f /tmp/golinks.log
tail -f /private/var/log/com.apple.xpc.launchd/launchd.log
```

### Docker

Build image and run
```
docker build . -t crhuber/golinks:latest
docker run -p 8998:8998 crhuber/golinks
```

### Browser Extension Redirect Setup (Recommended)

For a true go/alias experience in your browser, you can use a redirect browser extension instead of complex DNS and port configuration.

#### Using Redirector Extension

1.  **Install the Extension:**
    *   Redirector for Chrome
    *   Redirector for Firefox

2.  **Configure a Redirect Rule:**
    *   Open extension settings
    *   Add a new redirect with these settings:
        *   **Description:** `GoLink Redirector`
        *   **Example URL:** `http://go/docs`
        *   **Include pattern:** `^http://go/(.*)$`
        *   **Redirect to:** `http://go/$1` or if you didn't update the `/etc/hosts` file `http://localhost:8998/$1`
        *   **Pattern type:** `Regular Expression`

3.  **Usage:**
    *   Simply type `go/docs` (or any alias) in your browser's address bar.
    *   The extension will redirect to your local GoLink server.
    *   **Note:** You must first open the link using `http://go/{alias}` for each link before the browser will recognize this as a valid path. (You can use the `open` command to do this quickly). This is because the browser will try search first if the url is not recognized.


## FAQ

* How can I see all the links available

    http://go:8998/


* How do variable links work?

    You can create dynamic links that accept parameters by adding `{*}` to your destination URL.
    
    For example:
    *   Keyword: `gh`
    *   Destination: `https://github.com/{*}`
    
    When a user types `go/gh/torvalds`, the `{*}` will be replaced with `torvalds`, redirecting to `https://github.com/torvalds`.
    
    You can also use multiple variables:
    *   Keyword: `jira`
    *   Destination: `https://{*}.jira.com/browse/{*}`
    
    Typing `go/jira/github/PROJ-123` will redirect to `https://github.jira.com/browse/PROJ-123`.

## Troubleshooting

- If you change the port of the API. Be sure that you change the frontend index.html to connect to the same port

## Developing

I use [air](https://github.com/cosmtrek/air) for live reloading Go apps.
Just run

```
> air

watching .
building...
running...
INFO[0000] Starting server on port :8998
```

## Roadmap

- Add CLI interface to adding/removing/searching links from command line
- UI Refactor

## Contributing

If you find bugs, please open an issue first. If you have feature requests, I probably will not honor it because this project is being built mostly to suit my personal workflow and preferences.
