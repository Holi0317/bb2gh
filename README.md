# bb2gh

Migrate PR from Bitbucket to Github as issue.

## Why

We want to migrate from bitbucket to github. However, the PR number will start from 1 if
we do the migration.

This raise the question: What does PR#123 stands for? Is that a github PR or a bitbucket
PR?

To solve (workaround, aka a hack) this, we are going to create a bunch of issues on github
side. This will make PR number smaller than 2000 -> bitbucket, while all github PR starts
from 2001.

## Feature

- Reserve github issue/PR number by creating placeholder issues
- Migrate Bitbucket PR to github as issue

Note that the code quality is a bit, meh. This tool is mostly for personal use. If you
intend to use this for your migration, please open a test repository and check if this
meet your requirement.

## Usage

Download go compiler. I didn't setup github action to compile binary on push.

> :warning: Tell everyone to turn off their notification on the migration repo on Github.

Run following commands:

```bash
export GITHUB_TOKEN=<GITHUB_TOKEN>
export BITBUCKET_USER=<BB_USER>
export BITBUCKET_PASSWORD=<BB_PASSWORD>
# Let's say max bitbucket PR number is #1600. So we want to reserve github issue up to 2000
# Throw this to a VM and let it run overnight. Github got some serious rate limit in place (which make sense).
go run ./main.go reserve --to 2000 --github-repository octo-org/bb2gh-test

# After reserving, technically developers can start making PR as the issue number has been bumped.
# Run this command to start migrate the pr content to the issue
# Note on the `{1..2000}`, which means expand 1-2000 numbers on bash. You can run a single PR migration base on the argument.
go run ./main.go migrate --github-repository octo-org/bb2gh-test --bitbucket-repository octo-old-org/bb2gh-test {1..2000}
```

You can replace `go run` with compiled binary, which will be useful if you throw this
thing to a cloud vm.

There will be a log file written to `bb2gh.log`. Check for warning messages to see which
migration failed and rerun it if needed.

### Generate Github token

Github Settings -> Developers settings -> Personal access token -> Generate new token

You will need to grant permission on all repo scope.

### Generate Bitbucket token

Generate a app token (I forgot where to get that one). Just make sure you are using your
bitbucket username (not your atlassian email!) when passing the argument.
