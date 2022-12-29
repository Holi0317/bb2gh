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

> :warning: Tell everyone to turn off their notification on the migration repo on Github.

1. Create a new repository
2. Add `BITBUCKET_USER` and `BITBUCKET_PASSWORD` to the repository secret
3. On a new branch, create github workflow config (attached below)
4. Trigger `reserve` job multiple times until it succeed
5. Trigger `migrate`. Probably batch the jobs per 80-100 PR

We are doing this weird hack so that we can issue a new github token easily and avoid
github's rate limit. If we hit rate limit, we just start a new job with a new github
token.

(Remember to change the release number in the workflow files)

```yaml
# .github/workflows/reserve.yml
name: Reserve

on:
  workflow_dispatch:
    inputs:
      to:
        description: "Target issue number to get to"
        required: true
        type: string
      github-repository:
        description: "Github repository to work with"
        required: true
        type: string

permissions:
  issues: write
  pull-requests: read

jobs:
  run:
    runs-on: ubuntu-latest
    steps:
      - name: "Download binary"
        run: "curl -L -o bb2gh https://github.com/Holi0317/bb2gh/releases/download/v1.1.0/bb2gh && chmod +x bb2gh"

      - name: Run
        run: "./bb2gh reserve --github-token ${{ github.token }} --to ${{ inputs.to }} --github-repository ${{ inputs.github-repository }}"
```

```yaml
# .github/workflows/migrate.yml

name: Migrate

on:
  workflow_dispatch:
    inputs:
      from:
        description: "Target PR number to get from"
        required: true
        type: string
      to:
        description: "Target PR number to get to"
        required: true
        type: string
      github-repository:
        description: "Github repository to work with"
        required: true
        type: string
      bitbucket-repository:
        description: "Bitbucket repository to work with"
        required: true
        type: string

permissions:
  issues: write
  pull-requests: read

jobs:
  run:
    runs-on: ubuntu-latest
    steps:
      - name: "Download binary"
        # Plz update this link to the latest release
        run: "curl -L -o bb2gh https://github.com/Holi0317/bb2gh/releases/download/v1.1.0/bb2gh && chmod +x bb2gh"

      - name: Run
        shell: bash
        env:
          GITHUB_TOKEN: ${{ github.token }}
          BITBUCKET_USER: ${{ secrets.BITBUCKET_USER }}
          BITBUCKET_PASSWORD: ${{ secrets.BITBUCKET_PASSWORD }}
        run: "./bb2gh migrate --github-repository ${{ inputs.github-repository }} --bitbucket-repository ${{ inputs.bitbucket-repository }} {${{ inputs.from }}..${{ inputs.to }}}"
```

### Generate Bitbucket token

Generate a app token (I forgot where to get that one). Just make sure you are using your
bitbucket username (not your atlassian email!) when passing the argument.

### Driver shell script template

```bash
#!/bin/bash

bbrepo=bb-repo/repo
repo=octo-org/repo

batch_size=50
to=4000

# For reserve
for i in {1..20}; do
	gh workflow -R "$repo" run Reserve -f to=${to} -f "github-repository=${repo}"
	sleep 40
done

for i in {1..20}; do
	low=$(((i - 1) * batch_size + 1))
	up=$((i * batch_size))

	echo "Processing from ${low} to ${up}"

	gh workflow -R "$repo" run Migrate -f from=$low -f to=$up -f "github-repository=${repo}" -f "bitbucket-repository=${bbrepo}"
	sleep 60
done

# Show the date for prepare next run
date
```
