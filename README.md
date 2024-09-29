# Snapshot Autovoter

This is a simple autovoter for the [Snapshot](https://snapshot.org/) platform. It allows you to vote on proposals automatically based on a set of rules.

The autovoter will watch proposals creation on Thursdays morning and vote on them if the title matches

## Configuration

```bash
VOTER_ADDRESS=
VOTER_PRIVATE_KEY=
# Exact proposal choice label
PROPOSAL_CHOICE=
# Space ENS
PROPOSAL_SPACE=
# Title of the proposal, must be the start of the title
# using [HasPrefix](https://pkg.go.dev/strings#HasPrefix)
PROPOSAL_TITLE=
```
