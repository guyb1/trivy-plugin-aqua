package metadata

import (
	"fmt"
	"strings"

	"github.com/aquasecurity/trivy-plugin-aqua/pkg/git"
	"github.com/aquasecurity/trivy-plugin-aqua/pkg/log"
	"github.com/argonsecurity/go-environments/enums"
	"github.com/argonsecurity/go-environments/models"
	"github.com/pkg/errors"
	"github.com/samber/lo"
)

func GetBaseRef(envconfig *models.Configuration) (r string) {
	branch := lo.Ternary(envconfig.PullRequest.TargetRef.Branch != "", envconfig.PullRequest.TargetRef.Branch, "master")

	if envconfig.Repository.Source == enums.Azure {
		return GetFullBranchName(strings.TrimPrefix(branch, "refs/heads/"), "origin")
	}

	if strings.ToLower(envconfig.Builder) == string(enums.Jenkins) {
		return GetFullBranchName(branch, "upstream")
	}

	return GetFullBranchName(branch, "origin")
}

func GetFullBranchName(branchName, remoteFallback string) string {
	branchPattern := fmt.Sprintf("*/%s", branchName)
	branch := fmt.Sprintf("%s/%s", remoteFallback, branchName)
	out, err := git.GitExec("branch", "-a", "--list", branchPattern, "--format=%(refname)", "--sort=-refname")
	if err != nil {
		log.Logger.Info("err is not nil")
		if strings.Contains(err.Error(), "error: unknown option `format") { // git version < 2.7.0
			log.Logger.Warnf("Git version is too old. Please upgrade to 2.7.0 or newer")
			log.Logger.Debug("Trying to get branch name with old git version")
			out, err = git.GitExec("branch", "-a", "--list", branchPattern)
			if err != nil {
				log.Logger.Error(errors.Wrap(err, "failed git branch -a"))
			}
			if out != "" {
				branchs := lo.Reverse(lo.FilterMap(strings.Split(out, "\n"), func(b string, _ int) (string, bool) {
					trimmed := trimBranchName(b)
					if trimmed == "" {
						return "", false
					}
					return trimmed, true
				}))
				if len(branchs) > 0 {
					return branchs[0]
				}
			}
		} else {
			log.Logger.Error(errors.Wrap(err, "failed git branch -a"))
		}
	}

	if out != "" {
		branchs := strings.Split(out, "\n")
		if len(branchs) > 0 {
			branch = trimBranchName(branchs[0])
		}
	}

	log.Logger.Debugf("Using branch full name: %s", branch)
	return branch
}

func trimBranchName(branch string) string {
	branch = strings.TrimSpace(branch)
	branch = strings.TrimPrefix(strings.TrimSpace(branch), "refs/")
	return strings.TrimPrefix(strings.TrimSpace(branch), "remotes/")
}
