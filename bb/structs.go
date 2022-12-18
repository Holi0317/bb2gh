package bb

import (
	"time"

	"github.com/sirupsen/logrus"
)

type Links struct {
	HTML struct {
		Href string `json:"href"`
	} `json:"html"`
}

type Markup struct {
	// The text as it was typed by a user.
	Raw string `json:"raw"`

	// The type of markup language the raw content is to be interpreted in.
	// Valid values: markdown, creole, plaintext
	Markup string `json:"markup"`

	// The user's content rendered as HTML.
	HTML string `json:"html"`
}

func (m *Markup) check(prid int) {
	if m == nil {
		return
	}

	if m.Markup != "creole" {
		return
	}

	logrus.WithFields(logrus.Fields{
		"prid":   prid,
		"markup": m.Markup,
	}).Warn("Got unsupported creole markup")
}

type Account struct {
	AccountID   string `json:"account_id"`
	DisplayName string `json:"display_name"`
	Nickname    string `json:"nickname"`
}

type PREndpoint struct {
	Repository struct {
		Links Links  `json:"links"`
		UUID  string `json:"uuid"`

		// The concatenation of the repository owner's username and the slugified name, e.g. "evzijst/interruptingcow".
		// This is the same string used in Bitbucket URLs.
		FullName string `json:"full_name"`
	} `json:"repository"`

	Branch struct {
		Name string `json:"name"`
	} `json:"branch"`

	Commit struct {
		Hash string `json:"hash"`
	} `json:"commit"`
}

type PullRequest struct {
	Links struct {
		HTML struct {
			Href string `json:"href"`
		} `json:"html"`
	} `json:"links"`
	ID    int    `json:"id"`
	Title string `json:"title"`

	Rendered struct {
		Description *Markup `json:"description"`
		Title       *Markup `json:"title"`
		Reason      *Markup `json:"reason"`
	} `json:"rendered"`

	Summary Markup `json:"summary"`

	// The pull request's current status.
	//
	// Valid values: OPEN, MERGED, DECLINED, SUPERSEDED
	State string `json:"state"`

	Author Account `json:"author"`

	Source      PREndpoint  `json:"source"`
	Destination PREndpoint  `json:"destination"`
	MergeCommit *PREndpoint `json:"merge_commit"`

	CommentCount int `json:"comment_count"`

	ClosedBy *Account `json:"closed_by"`

	// Explains why a pull request was declined. This field is only applicable to pull requests in rejected state.
	Reason string `json:"reason"`

	CreatedOn time.Time `json:"created_on"`
	UpdatedOn time.Time `json:"updated_on"`

	Reviewers []Account `json:"reviewers"`
}
