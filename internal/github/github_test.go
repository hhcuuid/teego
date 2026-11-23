package github

import (
	"testing"
)

func TestGitWorkflow(t *testing.T) {
	if err := CloneMarketRepo(); err != nil {
		t.Fatalf("CloneMarketRepo() failed: %v", err)
	}
	defer CleanFiles()

	if err := PushDummy(); err != nil {
		t.Fatalf("PushDummy() failed: %v", err)
	}

	if err := AddCommitPush(); err != nil {
		t.Fatalf("AddCommitPush() failed: %v", err)
	}
}
