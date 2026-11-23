package github

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

func init() {
	token := os.Getenv("GH_TOKEN")
	if token == "" {
		log.Panic("GH_TOKEN is unset")
	}
	os.Unsetenv("GH_TOKEN")

	repoURL = fmt.Sprintf("https://x-access-token:%s@github.com/%s/%s.git",
		token,
		OWNER,
		REPO,
	)
}

const (
	OWNER       = "hhcuuid"
	REPO        = "teego"
	WORKING_DIR = "_tmp/teego"
	DATA_DIR    = WORKING_DIR + "/docs/data"
	DUMMY_FILE  = DATA_DIR + "/DUMMY.json"
)

var repoURL string

func CloneMarketRepo() error {

	// 1. Shallow clone (single branch, depth 1)
	err := runGit("", "clone", "--depth=1", "--single-branch", repoURL, WORKING_DIR)
	if err != nil {
		return err
	}
	// defer os.RemoveAll(workDir)

	// // 2. Process 1M+ files concurrently with a bounded worker pool
	// if err := processFiles(ctx, workDir); err != nil {
	// 	panic(err)
	// }

	return nil
}

func AddCommitPush() error {
	// 3. Stage, commit, and push in one packfile
	if err := runGit(WORKING_DIR, "config", "user.name", "Automation Bot"); err != nil {
		return err
	}
	if err := runGit(WORKING_DIR, "config", "user.email", "market@teego.dev"); err != nil {
		return err
	}
	if err := runGit(WORKING_DIR, "add", "-A"); err != nil {
		return err
	}
	if err := runGit(WORKING_DIR, "commit", "-m", "Bulk update market data"); err != nil {
		return err
	}
	if err := runGit(WORKING_DIR, "push", "origin", "HEAD"); err != nil {
		return err
	}

	return nil
}

func CleanFiles() error {
	return os.RemoveAll(WORKING_DIR)
}

func PushDummy() error {
	// create all the files up to DUMMY_FILE

	return os.WriteFile(
		DUMMY_FILE,
		fmt.Appendf(nil, `{"testTime":%d}`, time.Now().Unix()),
		0440,
	)
}

func runGit(dir string, args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
