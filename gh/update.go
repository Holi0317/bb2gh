package gh

import (
	_ "embed"
	"errors"
	"strings"
	"text/template"
	"time"

	"github.com/google/go-github/v48/github"
	"github.com/holi0317/bb2gh/bb"
	"github.com/sirupsen/logrus"
)

//go:embed issue.md.tmpl
var issueTemplateRaw string

var issueTemplate *template.Template

func init() {
	tmpl, err := template.New("issue").Parse(issueTemplateRaw)
	if err != nil {
		panic(err)
	}

	issueTemplate = tmpl
}

type templateInput struct {
	PR     bb.PullRequest
	Now    string
	Author string
}

func (c *Client) UpdateIssue(number int, source *bb.PullRequest) error {
	input := templateInput{
		PR:  *source,
		Now: time.Now().UTC().Format("2006-01-02T15:04:05.999Z"),
	}

	sb := &strings.Builder{}
	err := issueTemplate.Execute(sb, input)
	if err != nil {
		return err
	}

	issueReq := &github.IssueRequest{
		Title: github.String(source.Rendered.Title.Raw),
		Body:  github.String(sb.String()),
	}

	for i := 0; i < retry; i++ {
		log := logrus.WithFields(logrus.Fields{
			"retry":  i,
			"number": number,
		})

		_, _, err = c.ghc.Issues.Edit(c.ctx, c.owner, c.repo, number, issueReq)
		if sleep, ok := isRateLimit(err); ok {
			log.WithField("sleep", sleep).Debug("Hit rate limit. Sleeping before retry")
			time.Sleep(sleep)
			continue
		}

		if err != nil {
			return err
		}

		return nil
	}

	return errors.New("Rate limit")
}
