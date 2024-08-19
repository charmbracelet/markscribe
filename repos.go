package main

import (
	"context"
	"fmt"
	"slices"
	"sort"

	"github.com/shurcooL/githubv4"
)

var recentContributionsQuery struct {
	User struct {
		Login                   githubv4.String
		ContributionsCollection struct {
			CommitContributionsByRepository []struct {
				Contributions struct {
					Edges []struct {
						Cursor githubv4.String
						Node   struct {
							OccurredAt githubv4.DateTime
						}
					}
				} `graphql:"contributions(first: 1)"`
				Repository qlRepository
			} `graphql:"commitContributionsByRepository(maxRepositories: 100)"`
		}
	} `graphql:"user(login:$username)"`
}

var recentPullRequestsQuery struct {
	User struct {
		Login        githubv4.String
		PullRequests struct {
			TotalCount githubv4.Int
			Edges      []struct {
				Cursor githubv4.String
				Node   qlPullRequest
			}
		} `graphql:"pullRequests(first: $count, orderBy: {field: CREATED_AT, direction: DESC})"`
	} `graphql:"user(login:$username)"`
}

var recentReposQuery struct {
	User struct {
		Login        githubv4.String
		Repositories struct {
			TotalCount githubv4.Int
			Edges      []struct {
				Cursor githubv4.String
				Node   qlRepository
			}
		} `graphql:"repositories(first: $count, privacy: PUBLIC, isFork: $isFork, ownerAffiliations: OWNER, orderBy: {field: CREATED_AT, direction: DESC})"`
	} `graphql:"repositoryOwner(login: $owner)"`
}

var recentReleasesQuery struct {
	User struct {
		Login                     githubv4.String
		RepositoriesContributedTo struct {
			TotalCount githubv4.Int
			Edges      []struct {
				Cursor githubv4.String
				Node   struct {
					qlRepository
					Releases qlReleases `graphql:"releases(first: 10, orderBy: {field: CREATED_AT, direction: DESC})"`
				}
			}
		} `graphql:"repositoriesContributedTo(first: 100, after:$after includeUserRepositories: true, contributionTypes: COMMIT, privacy: PUBLIC)"`
	} `graphql:"user(login:$username)"`
}

/*
Order by stars

	{
	  repositoryOwner(login: "charmbracelet") {
	    id
	    login
	    repositories(
	      first: 5
	      privacy: PUBLIC
	      orderBy: {field: STARGAZERS, direction: DESC}
	    ) {
	      edges {
	        node {
	          name
	          description
	          url
	        }
	      }
	    }
	  }
	}
*/
func popularRepos(owner string, count int) []Repo {
	var query struct {
		Owner struct {
			Repositories struct {
				Edges []struct {
					Node qlRepository
				}
			} `graphql:"repositories(first: $count, privacy: PUBLIC, orderBy: {field: STARGAZERS, direction: DESC})"`
		} `graphql:"repositoryOwner(login: $owner)"`
	}

	fmt.Println("Finding popular repos...")

	var repos []Repo
	variables := map[string]interface{}{
		"owner": githubv4.String(owner),
		"count": githubv4.Int(count + 1), // +1 in case we encounter the meta-repo itself
	}
	err := gitHubClient.Query(context.Background(), &query, variables)
	if err != nil {
		panic(err)
	}

	for _, v := range query.Owner.Repositories.Edges {
		// ignore meta-repo
		if string(v.Node.NameWithOwner) == fmt.Sprintf("%s/%s", owner, username) {
			continue
		}
		if len(repos) == count {
			break
		}

		repos = append(repos, repoFromQL(v.Node))
	}

	fmt.Printf("Found %d repos!\n", len(repos))
	return repos
}

var repoQuery struct {
	Repository struct {
		Description githubv4.String
		Owner       struct {
			Login githubv4.String
		}
		Name          githubv4.String
		NameWithOwner githubv4.String
		IsPrivate     githubv4.Boolean
		URL           githubv4.String
		Stargazers    struct {
			TotalCount githubv4.Int
		}
		Releases qlReleases `graphql:"releases(last: 1)"`
	} `graphql:"repository(owner:$owner, name:$name)"`
}

var repoRecentReleasesQuery struct {
	Repository struct {
		Releases qlReleases `graphql:"releases(first: $count, orderBy: {field: CREATED_AT, direction: DESC})"`
	} `graphql:"repository(name: $name, owner: $owner)"`
}

func recentContributions(count int) []Contribution {
	var contributions []Contribution
	variables := map[string]interface{}{
		"username": githubv4.String(username),
	}
	err := gitHubClient.Query(context.Background(), &recentContributionsQuery, variables)
	if err != nil {
		panic(err)
	}

	for _, v := range recentContributionsQuery.User.ContributionsCollection.CommitContributionsByRepository {
		// ignore meta-repo
		if string(v.Repository.NameWithOwner) == fmt.Sprintf("%s/%s", username, username) {
			continue
		}
		if v.Repository.IsPrivate {
			continue
		}

		c := Contribution{
			Repo:       repoFromQL(v.Repository),
			OccurredAt: v.Contributions.Edges[0].Node.OccurredAt.Time,
		}

		contributions = append(contributions, c)
	}

	sort.Slice(contributions, func(i, j int) bool {
		return contributions[i].OccurredAt.After(contributions[j].OccurredAt)
	})

	if len(contributions) > count {
		return contributions[:count]
	}
	return contributions
}

func recentPullRequests(count int) []PullRequest {
	var pullRequests []PullRequest
	variables := map[string]interface{}{
		"username": githubv4.String(username),
		"count":    githubv4.Int(count + 1), // +1 in case we encounter the meta-repo itself
	}
	err := gitHubClient.Query(context.Background(), &recentPullRequestsQuery, variables)
	if err != nil {
		panic(err)
	}

	for _, v := range recentPullRequestsQuery.User.PullRequests.Edges {
		// ignore meta-repo
		if string(v.Node.Repository.NameWithOwner) == fmt.Sprintf("%s/%s", username, username) {
			continue
		}
		if v.Node.Repository.IsPrivate {
			continue
		}

		pullRequests = append(pullRequests, pullRequestFromQL(v.Node))
		if len(pullRequests) == count {
			break
		}
	}

	return pullRequests
}

func recentCreatedRepos(owner string, count int) []Repo {
	var repos []Repo
	variables := map[string]interface{}{
		"owner":  githubv4.String(owner),
		"count":  githubv4.Int(count + 1), // +1 in case we encounter the meta-repo itself
		"isFork": githubv4.Boolean(false),
	}
	err := gitHubClient.Query(context.Background(), &recentReposQuery, variables)
	if err != nil {
		panic(err)
	}

	for _, v := range recentReposQuery.User.Repositories.Edges {
		// ignore meta-repo
		if string(v.Node.NameWithOwner) == fmt.Sprintf("%s/%s", owner, owner) {
			continue
		}

		repos = append(repos, repoFromQL(v.Node))
		if len(repos) == count {
			break
		}
	}

	return repos
}

func recentForkedRepos(owner string, count int) []Repo {
	var repos []Repo
	variables := map[string]interface{}{
		"owner":  githubv4.String(owner),
		"count":  githubv4.Int(count + 1), // +1 in case we encounter the meta-repo itself
		"isFork": githubv4.Boolean(true),
	}
	err := gitHubClient.Query(context.Background(), &recentReposQuery, variables)
	if err != nil {
		panic(err)
	}

	for _, v := range recentReposQuery.User.Repositories.Edges {
		// ignore meta-repo
		if string(v.Node.NameWithOwner) == fmt.Sprintf("%s/%s", owner, owner) {
			continue
		}

		repos = append(repos, repoFromQL(v.Node))
		if len(repos) == count {
			break
		}
	}
	return repos
}

func latestReleasedRepos(owner string, count int) []Repo {
	var query struct {
		Owner struct {
			Repositories struct {
				Edges []struct {
					Cursor githubv4.String
					Node   struct {
						qlRepository
						Release qlRelease `graphql:"latestRelease"`
					}
				}
			} `graphql:"repositories(first: 100, privacy: PUBLIC, orderBy: {field: UPDATED_AT, direction: DESC})"`
		} `graphql:"repositoryOwner(login: $owner)"`
	}

	var repos []Repo
	variables := map[string]interface{}{
		"owner": githubv4.String(owner),
	}
	err := gitHubClient.Query(context.Background(), &query, variables)
	if err != nil {
		panic(err)
	}

	for _, v := range query.Owner.Repositories.Edges {
		repo := repoFromQL(v.Node.qlRepository)
		release := releaseFromQL(v.Node.Release)
		repo.LastRelease = release
		if repo.LastRelease.Name != "" {
			repos = append(repos, repo)
		}
	}

	slices.SortFunc(repos, func(a, b Repo) int {
		return a.LastRelease.PublishedAt.Compare(b.LastRelease.PublishedAt)
	})
	slices.Reverse(repos)
	return repos[:count]
}

func recentReleases(count int) []Repo {
	var after *githubv4.String
	var repos []Repo

	for {
		variables := map[string]interface{}{
			"username": githubv4.String(username),
			"after":    after,
		}
		err := gitHubClient.Query(context.Background(), &recentReleasesQuery, variables)
		if err != nil {
			panic(err)
		}

		if len(recentReleasesQuery.User.RepositoriesContributedTo.Edges) == 0 {
			break
		}

		for _, v := range recentReleasesQuery.User.RepositoriesContributedTo.Edges {
			r := repoFromQL(v.Node.qlRepository)

			for _, rel := range v.Node.Releases.Nodes {
				if rel.IsPrerelease || rel.IsDraft {
					continue
				}
				if v.Node.Releases.Nodes[0].TagName == "" ||
					v.Node.Releases.Nodes[0].PublishedAt.Time.IsZero() {
					continue
				}
				r.LastRelease = releasesFromQL(v.Node.Releases)
				break
			}

			if !r.LastRelease.PublishedAt.IsZero() {
				repos = append(repos, r)
			}

			after = githubv4.NewString(v.Cursor)
		}
	}

	sort.Slice(repos, func(i, j int) bool {
		if repos[i].LastRelease.PublishedAt.Equal(repos[j].LastRelease.PublishedAt) {
			return repos[i].Stargazers > repos[j].Stargazers
		}
		return repos[i].LastRelease.PublishedAt.After(repos[j].LastRelease.PublishedAt)
	})

	if len(repos) > count {
		return repos[:count]
	}
	return repos
}

/*
	{
		repositoryOwner(login: "charmbracelet") {
		id
		login
		repositories(
			first: 5
			privacy: PUBLIC
			orderBy: {field: PUSHED_AT, direction: DESC}
		) {
			edges {
			node {
				name
				description
				url
			}
			}
		}
		}
	}
*/
func recentPushedRepos(owner string, count int) []Repo {
	var query struct {
		Owner struct {
			Repositories struct {
				Edges []struct {
					Node qlRepository
				}
			} `graphql:"repositories(first: $count, privacy: PUBLIC, orderBy: {field: PUSHED_AT, direction: DESC})"`
		} `graphql:"repositoryOwner(login: $owner)"`
	}
	var repos []Repo
	variables := map[string]interface{}{
		"count": githubv4.Int(count),
		"owner": githubv4.String(owner),
	}
	err := gitHubClient.Query(context.Background(), &query, variables)
	if err != nil {
		panic(err)
	}

	for _, v := range query.Owner.Repositories.Edges {
		repos = append(repos, repoFromQL(v.Node))
		if len(repos) == count {
			break
		}
	}
	return repos
}

func repo(owner, name string) Repo {
	variables := map[string]interface{}{
		"owner": githubv4.String(owner),
		"name":  githubv4.String(name),
	}
	err := gitHubClient.Query(context.Background(), &repoQuery, variables)
	if err != nil {
		panic(err)
	}
	repo := repoQuery.Repository
	return Repo{
		Owner:         string(repo.Owner.Login),
		Name:          string(repo.Name),
		NameWithOwner: string(repo.NameWithOwner),
		URL:           string(repo.URL),
		Description:   string(repo.Description),
		Stargazers:    int(repo.Stargazers.TotalCount),
		IsPrivate:     bool(repo.IsPrivate),
		LastRelease:   releasesFromQL(repo.Releases),
	}
}

func repoRecentReleases(owner, name string, count int) []Release {
	var releases []Release

	variables := map[string]interface{}{
		"owner": githubv4.String(owner),
		"name":  githubv4.String(name),
		"count": githubv4.Int(count),
	}
	err := gitHubClient.Query(context.Background(), &repoRecentReleasesQuery, variables)
	if err != nil {
		panic(err)
	}

	for _, rel := range repoRecentReleasesQuery.Repository.Releases.Nodes {
		if bool(rel.IsPrerelease) {
			continue
		}
		releases = append(releases, Release{
			Name:         string(rel.Name),
			TagName:      string(rel.TagName),
			PublishedAt:  rel.PublishedAt.Time,
			CreatedAt:    rel.CreatedAt.Time,
			URL:          string(rel.URL),
			IsLatest:     bool(rel.IsLatest),
			IsPreRelease: bool(rel.IsPrerelease),
			IsDraft:      bool(rel.IsDraft),
		})
	}

	return releases
}

/*
{
  user(login: "muesli") {
    login
    repositoriesContributedTo(first: 100, includeUserRepositories: true, contributionTypes: COMMIT) {
      totalCount
      edges {
        cursor
        node {
          id
          nameWithOwner
        }
      }
    }
  }
}

{
  user(login: "muesli") {
    login
    repositoriesContributedTo(first: 100, includeUserRepositories: true, contributionTypes: COMMIT) {
      totalCount
      edges {
        cursor
        node {
          id
          nameWithOwner
		  releases(first: 3, orderBy: {field: CREATED_AT, direction: DESC}) {
          	nodes {
          	  name
              PublishedAt
			  url
			  isPrerelease
			  isDraft
            }
          }
        }
      }
    }
  }
}

{
  user(login: "muesli") {
    login
    repositories(first: 10, privacy: PUBLIC, isFork: false, ownerAffiliations: OWNER, orderBy: {field: CREATED_AT, direction: DESC}) {
      totalCount
      edges {
        cursor
        node {
          id
          nameWithOwner
        }
      }
    }
  }
}

{
  user(login: "muesli") {
    login
    contributionsCollection {
      commitContributionsByRepository {
        contributions(first: 1) {
          edges {
            cursor
            node {
              occurredAt
            }
          }
        }
        repository {
          id
		  nameWithOwner
		  url
		  description
        }
      }
    }
  }
}
*/
