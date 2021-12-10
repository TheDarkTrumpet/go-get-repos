#Go-Git-Repos

# Introduction

This repository is a very simple go application that's intended to act as a backup for one's Github repos.  Given issues
in the CyberSecurity realm, it's important to have backups for data - even if it exists on the cloud.  Even for a company
as great as Github.

# Dependencies and Usage

This project is fairly straightforward, and intended to be easy to use.  The main dependency is you need to have Go
installed (unless you use a package), and a JSON configuration file stored in your `$HOME/.creds/gh_vars.json`.  This
file takes the format as follows:

```json
{
	"token": "<YOUR_GH_TOKEN>",
	"backup-dir": "/path/to/backup/dir/"
}
```

# Obtaining the Token
To obtain a token:
1. Click on your profile on the upper right and select "Settings"
2. Click on "Developer Settings"
3. Click on "Personal access tokens"
4. Click on "Generate new token"

For permissions, the only thing that's needed is the `Full control of private repositories`

I also recommend setting the expiration to never.
