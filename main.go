package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

var creds = "gh_vars.json" // Stored in ~/.creds/

func main() {
	// Get Token from .creds
	vars, err := loadVars()
	if err != nil {
		log.Fatal(err)
		return
	}

	// Get all existing repos in backup directory
	backupDirContents, err := readBackupDirectory(vars)
	if err != nil {
		log.Fatal(err)
		return
	}
	backupDirFiles := getBackupDirectoryNames(backupDirContents)
	fmt.Printf("Current Directory Contents: \n")
	PrintList(backupDirFiles)

	// Get all repos from Github
	availableRepos, err := readGithubAvailableRepos(vars)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Printf("Number of repositories to process: %v\n", len(availableRepos))

	// Xor operation to determine what to clone, and to clone
	reposToClone := getReposToClone(backupDirFiles, availableRepos)
	fmt.Printf("Number of repositories to clone, first: %v\n", len(reposToClone))

	// For all repos, do a git fetch
	numProcessed, err := cloneNonBackedupRepos(reposToClone, vars)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Number of repositories cloned, %v\n", numProcessed)

	// For folder glob, do a git fetch on each one
	numProcessed, err = updateAllCachedRepos(vars)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Number of repositories updated, %v\n", numProcessed)
}

type GHVars struct {
	Token       string   `json:"token"`
	Types       []string `json:"types"` // public, internal, private
	Affiliation string   `json:"org"`
	BackupDir   string   `json:"backup-dir"`
}

func loadVars() (GHVars, error) {
	user, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	varsFile := fmt.Sprintf("%s/.creds/%s", user, creds)
	PrintHeader(fmt.Sprintf("Loading Creds from %v", varsFile))

	var vars GHVars
	contents, err := ioutil.ReadFile(varsFile)
	err = json.Unmarshal(contents, &vars)
	if err != nil {
		return vars, err
	}
	return vars, err
}

func readBackupDirectory(vars GHVars) ([]os.FileInfo, error) {
	PrintHeader(fmt.Sprintf("Reading backup directory: %s", vars.BackupDir))
	files, err := ioutil.ReadDir(vars.BackupDir)
	return files, err
}

func getBackupDirectoryNames(files []os.FileInfo) []string {
	returnFiles := make([]string, len(files), len(files))
	for ix, file := range files {
		returnFiles[ix] = file.Name()
	}
	return returnFiles
}

func readGithubRepos(vars GHVars, page int) ([]*github.Repository, error) {
	PrintHeader("Reading available Github Repos off base username (based off token)")
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: vars.Token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)
	lopt := github.ListOptions{PerPage: 100, Page: page}
	//opt := &github.RepositoryListOptions{Affiliation: "owner", ListOptions: lopt}
	//repos, _, err := client.Repositories.List(ctx, "", opt)
	opt := &github.RepositoryListByOrgOptions{
		Type:        "Private", // "Private", Or Internal
		ListOptions: lopt,
	}
	repos, _, err := client.Repositories.ListByOrg(context.Background(), "UFGInsurance", opt)

	return repos, err
}

func readGithubAvailableRepos(vars GHVars) ([]github.Repository, error) {
	ghRepos := make([]github.Repository, 0, 0)
	page := 1
	for {
		repos, err := readGithubRepos(vars, page)

		if err != nil {
			return ghRepos, err
		}

		for x := 0; x < len(repos); x++ {
			ghRepos = append(ghRepos, *repos[x])
		}

		if len(repos) == 0 {
			break
		}
		page += 1
	}
	return ghRepos, nil
}

func getReposToClone(files []string, repos []github.Repository) []github.Repository {
	var ghToClone []github.Repository
	for _, repo := range repos {
		inCache := false
		for _, fhave := range files {
			if fhave == *repo.Name {
				inCache = true
				break
			}
		}
		if !inCache {
			ghToClone = append(ghToClone, repo)
		}
	}
	return ghToClone
}

func cloneNonBackedupRepos(repos []github.Repository, vars GHVars) (int, error) {
	PrintHeader(fmt.Sprintf("Cloning all non-cached repos, number to process: %v", len(repos)))
	numReposProcessed := 0
	err := error(nil)
	for _, repo := range repos {
		fmt.Printf("==> Processing: %v\n", *repo.Name)
		fullRepoURI := fmt.Sprintf("https://%v:%v@github.com/%v", *repo.Owner.Login, vars.Token, *repo.FullName)
		cmd := exec.Command("git", "clone", fullRepoURI)
		cmd.Dir = vars.BackupDir
		err := cmd.Run()
		if err != nil {
			return numReposProcessed, err
		}
		numReposProcessed += 1
	}
	return numReposProcessed, err
}

func updateAllCachedRepos(vars GHVars) (int, error) {
	backupDirectoryFiles, err := readBackupDirectory(vars)

	numReposUpdated := 0
	if err != nil {
		return numReposUpdated, err
	}

	PrintHeader(fmt.Sprintf("Updating all cached repos, number to process: %v", len(backupDirectoryFiles)))

	for _, repo := range backupDirectoryFiles {
		fmt.Printf("==> Processing: %v\n", repo.Name())
		cmd := exec.Command("git", "fetch")
		cmd.Dir = fmt.Sprintf("%s/%s", vars.BackupDir, repo.Name())
		err := cmd.Run()
		if err != nil {
			return numReposUpdated, err
		}
		numReposUpdated += 1
	}
	return numReposUpdated, err
}
