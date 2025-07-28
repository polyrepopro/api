package repositories

import (
	"fmt"

	"github.com/polyrepopro/api/git"

	g "github.com/go-git/go-git/v5"
)

type StatusReason struct {
	Type    StatusReasonType
	Error   error
	Message string
}

type StatusReasonType string

const (
	StatusReasonLocal    StatusReasonType = "local"
	StatusReasonRemote   StatusReasonType = "remote"
	StatusReasonUnknown  StatusReasonType = "unknown"
	StatusReasonUnpushed StatusReasonType = "unpushed"
	StatusReasonUnpulled StatusReasonType = "unpulled"
)

var StatusReasonMapping = map[StatusReasonType]StatusReason{
	StatusReasonLocal:    {Type: StatusReasonLocal, Error: nil},
	StatusReasonRemote:   {Type: StatusReasonRemote, Error: nil},
	StatusReasonUnknown:  {Type: StatusReasonUnknown, Error: nil},
	StatusReasonUnpushed: {Type: StatusReasonUnpushed, Error: nil},
	StatusReasonUnpulled: {Type: StatusReasonUnpulled, Error: nil},
}

type StatusCode string

const (
	StatusUnknown  StatusCode = "unknown"
	StatusMissing  StatusCode = "missing"
	StatusError    StatusCode = "error"
	StatusClean    StatusCode = "clean"
	StatusDirty    StatusCode = "dirty"
	StatusUnpushed StatusCode = "unpushed"
	StatusUnpulled StatusCode = "unpulled"
)

var StatusMessage = map[StatusCode]string{
	StatusUnknown:  "unknown",
	StatusMissing:  "missing",
	StatusError:    "error",
	StatusClean:    "clean",
	StatusDirty:    "dirty",
	StatusUnpushed: "unpushed",
	StatusUnpulled: "unpulled",
}

// GitStatusCodes maps git errors to their corresponding status codes.
var GitStatusCodes = map[error]StatusCode{
	g.ErrRepositoryNotExists:         StatusMissing,
	g.ErrRepositoryIncomplete:        StatusMissing,
	g.ErrRepositoryAlreadyExists:     StatusMissing,
	g.ErrRemoteNotFound:              StatusMissing,
	g.ErrRemoteExists:                StatusMissing,
	g.ErrBranchExists:                StatusError,
	g.ErrBranchNotFound:              StatusError,
	g.ErrTagExists:                   StatusError,
	g.ErrTagNotFound:                 StatusError,
	g.ErrFetching:                    StatusError,
	g.ErrInvalidReference:            StatusError,
	g.ErrAnonymousRemoteName:         StatusError,
	g.ErrWorktreeNotProvided:         StatusError,
	g.ErrIsBareRepository:            StatusError,
	g.ErrUnableToResolveCommit:       StatusError,
	g.ErrPackedObjectsNotSupported:   StatusError,
	g.ErrSHA256NotSupported:          StatusError,
	g.ErrAlternatePathNotSupported:   StatusError,
	g.ErrUnsupportedMergeStrategy:    StatusError,
	g.ErrFastForwardMergeNotPossible: StatusError,
}

type StatusResult struct {
	g.Status
	Code         StatusCode
	Message      string
	Enhanced     git.EnhancedStatus
	NeedsPush    bool
	NeedsPull    bool
	AheadCount   int
	BehindCount  int
}

func (s StatusCode) String() string {
	return StatusMessage[s]
}

// Status retrieves and processes the git status for a repository at the specified path.
//
// Arguments:
// - path: the file system path to the git repository
//
// Returns:
// - StatusResult: contains the status code and message for the repository
func Status(path string) StatusResult {
	gitStatus, err := git.Status(path)
	if err != nil {
		return StatusResult{
			Code:    StatusError,
			Message: err.Error(),
		}
	}

	var code StatusCode
	statusString := gitStatus.String()

	switch statusString {
	case "":
		code = StatusClean
	default:
		code = StatusDirty
	}

	return StatusResult{
		Status:  gitStatus,
		Code:    code,
		Message: statusString,
	}
}

// StatusWithRemote retrieves comprehensive status including remote branch comparison.
//
// Arguments:
// - path: the file system path to the git repository
// - remoteName: the name of the remote to compare against (e.g., "origin", "upstream")
//
// Returns:
// - StatusResult: comprehensive status including working tree and branch comparison
func StatusWithRemote(path, remoteName string) StatusResult {
	// Get enhanced status with remote comparison
	enhancedStatus, err := git.EnhancedStatusWithRemote(path, remoteName)
	if err != nil {
		return StatusResult{
			Code:    StatusError,
			Message: err.Error(),
		}
	}

	// Determine status code based on working tree and branch status
	var code StatusCode
	var message string

	statusString := enhancedStatus.WorkingTree.String()
	
	// Priority: working tree changes > remote sync status
	if enhancedStatus.HasChanges {
		code = StatusDirty
		message = statusString
	} else if enhancedStatus.Branch.NeedsPush && enhancedStatus.Branch.NeedsPull {
		code = StatusDirty // Diverged branches - need both push and pull
		if enhancedStatus.Branch.Ahead > 0 && enhancedStatus.Branch.Behind > 0 {
			message = fmt.Sprintf("diverged: %d ahead, %d behind", enhancedStatus.Branch.Ahead, enhancedStatus.Branch.Behind)
		} else {
			message = "diverged from remote"
		}
	} else if enhancedStatus.Branch.NeedsPush {
		code = StatusUnpushed
		if enhancedStatus.Branch.Ahead > 0 {
			message = fmt.Sprintf("%d commit(s) ahead", enhancedStatus.Branch.Ahead)
		} else {
			message = "needs push"
		}
	} else if enhancedStatus.Branch.NeedsPull {
		code = StatusUnpulled
		if enhancedStatus.Branch.Behind > 0 {
			message = fmt.Sprintf("%d commit(s) behind", enhancedStatus.Branch.Behind)
		} else {
			message = "needs pull"
		}
	} else {
		code = StatusClean
		message = "clean and up to date"
	}

	return StatusResult{
		Status:      enhancedStatus.WorkingTree,
		Code:        code,
		Message:     message,
		Enhanced:    enhancedStatus,
		NeedsPush:   enhancedStatus.Branch.NeedsPush,
		NeedsPull:   enhancedStatus.Branch.NeedsPull,
		AheadCount:  enhancedStatus.Branch.Ahead,
		BehindCount: enhancedStatus.Branch.Behind,
	}
}
