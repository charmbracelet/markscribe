# markscribe

[![Latest Release](https://img.shields.io/github/release/charmbracelet/markscribe.svg)](https://github.com/charmbracelet/markscribe/releases)
[![Build Status](https://github.com/charmbracelet/markscribe/workflows/build/badge.svg)](https://github.com/charmbracelet/markscribe/actions)
[![Go ReportCard](https://goreportcard.com/badge/charmbracelet/markscribe)](https://goreportcard.com/report/charmbracelet/markscribe)
[![GoDoc](https://godoc.org/github.com/golang/gddo?status.svg)](https://pkg.go.dev/github.com/charmbracelet/markscribe)

Your personal markdown scribe with template-engine and Git(Hub) & RSS powers 📜

You can run markscribe as a GitHub Action: [readme-scribe](https://github.com/charmbracelet/readme-scribe/)

## Usage

Render a template to stdout:

    markscribe template.tpl

Render to a file:

    markscribe -write /tmp/output.md template.tpl

## Installation

### Packages & Binaries

If you use Brew, you can simply install the package:

    brew install charmbracelet/tap/markscribe

Or download a binary from the [releases](https://github.com/charmbracelet/markscribe/releases)
page. Linux (including ARM) binaries are available, as well as Debian and RPM
packages.

### Build From Source

Alternatively you can also build `markscribe` from source. Make sure you have a
working Go environment (Go 1.16 or higher is required). See the
[install instructions](https://golang.org/doc/install.html).

To install markscribe, simply run:

    go get github.com/charmbracelet/markscribe

## Templates

You can find an example template to generate a GitHub profile README under
[`templates/github-profile.tpl`](templates/github-profile.tpl). Make sure to fill in (or remove) placeholders,
like the RSS-feed or social media URLs.

Rendered it looks a little like my own profile page: https://github.com/charmbracelet

## Functions

### RSS feed

```
{{range rss "https://domain.tld/feed.xml" 5}}
Title: {{.Title}}
URL: {{.URL}}
Published: {{humanize .PublishedAt}}
{{end}}
```

### Your recent contributions

```
{{range recentContributions 10}}
Name: {{.Repo.Name}}
Description: {{.Repo.Description}}
URL: {{.Repo.URL}})
Occurred: {{humanize .OccurredAt}}
{{end}}
```

This function requires GitHub authentication with the following API scopes:
`repo:status`, `public_repo`, `read:user`.

### Your recent pull requests

```
{{range recentPullRequests 10}}
Title: {{.Title}}
URL: {{.URL}}
State: {{.State}}
CreatedAt: {{humanize .CreatedAt}}
Repository name: {{.Repo.Name}}
Repository description: {{.Repo.Description}}
Repository URL: {{.Repo.URL}}
{{end}}
```

This function requires GitHub authentication with the following API scopes:
`repo:status`, `public_repo`, `read:user`.

### Repositories you recently starred

```
{{range recentStars 10}}
Name: {{.Repo.Name}}
Description: {{.Repo.Description}}
URL: {{.Repo.URL}})
Stars: {{.Repo.Stargazers}}
{{end}}
```

This function requires GitHub authentication with the following API scopes:
`repo:status`, `public_repo`, `read:user`.

### Repositories you recently created

```
{{range recentCreatedRepos "charmbracelet" 10}}
Name: {{.Name}}
Description: {{.Description}}
URL: {{.URL}})
Stars: {{.Stargazers}}
{{end}}
```

This function requires GitHub authentication with the following API scopes:
`repo:status`, `public_repo`, `read:user` or `read:org` if you provide an organization name.

### Repositories with the most stars

```
{{range popularRepos "charmbracelet" 10}}
Name: {{.Name}}
NameWithOwner: {{.NameWithOwner}}
Description: {{.Description}}
URL: {{.URL}})
Stars: {{.Stargazers}}
{{end}}
```

This function requires GitHub authentication with the following API scopes:
`read:org`, `public_repo`, `read:user`

> [!TIP]
> Use `{{with repo "charmbracelet .Name"}}` to create a pipeline that grabs additional information about the repo including releases.

### Custom GitHub repository

```
{{with repo "charmbracelet" "markscribe"}}
Name: {{.Name}}
Description: {{.Description}}
URL: {{.URL}}
Stars: {{.Stargazers}}
Is Private: {{.IsPrivate}}
Last Git Tag: {{.LastRelease.TagName}}
Last Release: {{humanize .LastRelease.PublishedAt}}
{{end}}
```

### Recent releases to a given repository

```
{{range recentRepoReleases "charmbracelet" "markscribe" 10}}
Name: {{.Name}}
Git Tag: {{.TagName}}
URL: {{.URL}}
Published: {{humanize .PublishedAt}}
CreatedAt: {{humanize .CreatedAt}}
IsPreRelease: {{.IsPreRelease}}
IsDraft: {{.IsDraft}}
IsLatest: {{.IsLatest}}
{{end}}
```

This function requires GitHub authentication with the following API scopes:
`repo:status`, `public_repo`, `read:user`.

### Forks you recently created

```
{{range recentForkedRepos "charmbracelet" 10}}
Name: {{.Name}}
Description: {{.Description}}
URL: {{.URL}})
Stars: {{.Stargazers}}
{{end}}
```

This function requires GitHub authentication with the following API scopes:
`repo:status`, `public_repo`, `read:user` or `read:org` if you provide an organization name.

### Latest released projects

```
{{range latestReleasedRepos "charmbracelet" 10}}
Name: {{.Name}}
Description: {{.Description}}
URL: {{.URL}})
Stars: {{.Stargazers}}
Last Release Name: {{.LastRelease.TagName}}
Last Release URL: {{.LastRelease.URL}}
Last Release Date: {{humanize .LastRelease.PublishedAt}}
{{end}}
```

This function requires GitHub authentication with the following API scopes:
`repo:status`, `public_repo`, `read:user`, `read:org`.

### Recent releases you contributed to

```
{{range recentReleases 10}}
Name: {{.Name}}
Git Tag: {{.LastRelease.TagName}}
URL: {{.LastRelease.URL}}
Published: {{humanize .LastRelease.PublishedAt}}
{{end}}
```

This function requires GitHub authentication with the following API scopes:
`repo:status`, `public_repo`, `read:user`.

### Recent pushes

```
{{range recentPushedRepos "charmbracelet" 10}}
Name: {{.Name}}
URL: {{.URL}}
Description: {{.Description}}
Stars: {{.Stargazers}}
{{end}}
```

This function requires GitHub authentication with the following API scopes:
`public_repo`, `read:org`.

> [!TIP]
> Use `{{with repo "charmbracelet .Name"}}` to create a pipeline that grabs additional information about the repo including releases.

### Your published gists

```
{{range gists 10}}
Name: {{.Name}}
Description: {{.Description}}
URL: {{.URL}}
Created: {{humanize .CreatedAt}}
{{end}}
```

This function requires GitHub authentication with the following API scopes:
`repo:status`, `public_repo`, `read:user`.

### Your latest followers

```
{{range followers 5}}
Username: {{.Login}}
Name: {{.Name}}
Avatar: {{.AvatarURL}}
URL: {{.URL}}
{{end}}
```

This function requires GitHub authentication with the following API scopes:
`read:user`.

### Your sponsors

```
{{range sponsors 5}}
Username: {{.User.Login}}
Name: {{.User.Name}}
Avatar: {{.User.AvatarURL}}
URL: {{.User.URL}}
Created: {{humanize .CreatedAt}}
{{end}}
```

This function requires GitHub authentication with the following API scopes:
`repo:status`, `public_repo`, `read:user`, `read:org`.

### Your GoodReads reviews

```
{{range goodReadsReviews 5}}
- {{.Book.Title}} - {{.Book.Link}} - {{.Rating}} - {{humanize .DateUpdated}}
{{- end}}
```

This function requires GoodReads API key!

### Your GoodReads currently reading books

```
{{range goodReadsCurrentlyReading 5}}
- {{.Book.Title}} - {{.Book.Link}} - {{humanize .DateUpdated}}
{{- end}}
```

This function requires GoodReads API key!

### Your Literal.club currently reading books

```
{{range literalClubCurrentlyReading 5}}
- {{.Title}} - {{.Subtitle}} - {{.Description}} - https://literal.club/_YOUR_USERNAME_/book/{{.Slug}}
  {{- range .Authors }}{{ .Name }}{{ end }}
{{- end}}
```

This function requires a `LITERAL_EMAIL` and `LITERAL_PASSWORD`.

### Your Wakatime total coding time for the week (human readable)

```
{{ wakatimeData.HumanReadableTotal }}
```

This function requires a `WAKATIME_API_KEY` and potentialy a `WAKATIME_URL`.

### Your top Wakatime languages for the week

```
{{ range wakatimeData.Languages | chunk 5 | first }}
- {{ .Name }}: {{ .Percent }}%
{{- end}}
```

This function requires a `WAKATIME_API_KEY` and potentialy a `WAKATIME_URL`.

### Your top Wakatime Languages formated as a bar

```
{{ wakatimeData.Languages | wakatimeCategoryBar 5 }}
```

```text
TypeScript   5h 6m 42s    █░░░░░░░░░░░░░░░░░░░░░░░░  45.56%
Go           3h 18m 11s   █░░░░░░░░░░░░░░░░░░░░░░░░  29.44%
Astro        53m 35s      █░░░░░░░░░░░░░░░░░░░░░░░░  7.96%
Markdown     31m 42s      █░░░░░░░░░░░░░░░░░░░░░░░░  4.71%
Bash         25m 51s      █░░░░░░░░░░░░░░░░░░░░░░░░  3.84%
```

This can also be done with `Categories`, `Projects`, `Editors`, `OperatingSystems`, `Machines`, and `Dependencies`

This function requires a `WAKATIME_API_KEY` and potentialy a `WAKATIME_URL`.

### Other Wakatime data

You can find the full list of the data you can get at [wakatimeTypes.go](./wakatimeTypes.go)

This function requires a `WAKATIME_API_KEY` and potentialy a `WAKATIME_URL`.

## Template Engine

markscribe uses Go's powerful template engine. You can find its documentation
here: https://golang.org/pkg/text/template/

## Template Helpers

markscribe comes with [sprout](https://docs.atom.codes/sprout) and a few more template helpers:

To format timestamps, call `humanize`:

```
{{humanize .Timestamp}}
```

To limit the length of an array 

```
{{.Array | chunk 5 | first}}
```

## GitHub Authentication

In order to access some of GitHub's API, markscribe requires you to provide a
valid GitHub token in an environment variable called `GITHUB_TOKEN`. You can
create a new token by going to your profile settings:

`Developer settings` > `Personal access tokens` > `Generate new token`

## GoodReads API key

In order to access some of GoodReads' API, markscribe requires you to provide a
valid GoodReads key in an environment variable called `GOODREADS_TOKEN`. You can
create a new token by going [here](https://www.goodreads.com/api/keys).
Then you need to go to your repository and add it, `Settings -> Secrets -> New secret`.
You also need to set your GoodReads user ID in your secrets as `GOODREADS_USER_ID`.

## Wakatime Authentication

In order to access any wakatime data you need to provide your api key as `WAKATIME_API_KEY`. If you use an 
alternative wakatime server such as [wakapi](https://github.com/muety/wakapi)
or [hackatime](https://github.com/hackclub/hackatime) then export the base wakatime compatible route as `WAKATIME_URL` e.g. `https://waka.hackclub.com/api/compat/wakatime/v1/`.

## FAQ

Q: That's awesome, but can you expose more APIs and data?  
A: Of course, just open a new issue and let me know what you'd like to do with markscribe!

Q: That's awesome, but I don't have my own server to run this on. Can you help?  
A: Check out [readme-scribe](https://github.com/charmbracelet/readme-scribe/), a GitHub Action that runs markscribe for you!
