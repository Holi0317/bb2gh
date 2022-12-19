package gh

import (
	"context"
	_ "embed"
	"strings"
	"text/template"
	"time"

	"github.com/google/go-github/v48/github"
	"github.com/holi0317/bb2gh/bb"
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

func (c *Client) UpdateIssue(ctx context.Context, number int, source *bb.PullRequest) error {
	client := c.get(ctx)

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

	_, _, err = client.Issues.Edit(ctx, c.owner, c.repo, number, issueReq)
	if err != nil {
		return err
	}

	return nil
}
