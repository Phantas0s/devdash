package platform

import (
	"context"
	"strconv"

	"github.com/shuheiktgw/go-travis"
)

const noToken = "none"

type TravisCI struct {
	client *travis.Client
}

func NewTravisCI(token string) *TravisCI {
	if token == noToken {
		token = ""
	}

	return &TravisCI{
		client: travis.NewClient(travis.ApiOrgUrl, token),
	}
}

func (tc TravisCI) Builds(repository string, owner string, limit int64) ([][]string, error) {
	include := []string{
		"build.repository",
		"build.state",
		"build.duration",
		"build.started_at",
		"build.finished_at",
	}

	var builds []*travis.Build
	var err error
	if repository == "" || owner == "" {
		builds, _, err = tc.client.Builds.List(
			context.Background(),
			&travis.BuildsOption{
				Include: include,
				Limit:   int(limit),
			},
		)
		if err != nil {
			return nil, err
		}
	} else {
		builds, _, err = tc.client.Builds.ListByRepoSlug(
			context.Background(),
			createRepoName(repository, owner),
			&travis.BuildsByRepoOption{
				Include: include,
				Limit:   int(limit),
				SortBy:  "build.finished_at",
			},
		)
		if err != nil {
			return nil, err
		}
	}

	table := formatBuilds(builds, limit)

	return table, nil
}

func formatBuilds(builds []*travis.Build, limit int64) [][]string {
	table := make([][]string, limit+1)

	table[0] = []string{
		"Repository",
		"State",
		"Duration",
		"Finished At",
	}

	if len(builds) > 0 {
		for k, v := range builds {
			table[k+1] = append(table[k+1], *v.Repository.Name)
			table[k+1] = append(table[k+1], *v.State)

			// Can be nil if
			if v.Duration != nil {
				table[k+1] = append(table[k+1], strconv.FormatUint(uint64(*v.Duration), 10))
			} else {
				table[k+1] = append(table[k+1], "Running")
			}
			if v.FinishedAt != nil {
				table[k+1] = append(table[k+1], *v.FinishedAt)
			} else {
				table[k+1] = append(table[k+1], "Running")
			}
		}
	}

	return table
}

func createRepoName(repository string, owner string) string {
	return owner + "/" + repository
}
