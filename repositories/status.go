package repositories

import (
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
	StatusReasonLocal   StatusReasonType = "local"
	StatusReasonRemote  StatusReasonType = "remote"
	StatusReasonUnknown StatusReasonType = "unknown"
)

var StatusReasonMapping = map[StatusReasonType]StatusReason{
	StatusReasonLocal:   {Type: StatusReasonLocal, Error: nil},
	StatusReasonRemote:  {Type: StatusReasonRemote, Error: nil},
	StatusReasonUnknown: {Type: StatusReasonUnknown, Error: nil},
}

type StatusCode string

const (
	StatusUnknown StatusCode = "unknown"
	StatusMissing StatusCode = "missing"
	StatusError   StatusCode = "error"
	StatusClean   StatusCode = "clean"
	StatusDirty   StatusCode = "dirty"
)

var StatusMessage = map[StatusCode]string{
	StatusUnknown: "unknown",
	StatusMissing: "missing",
	StatusError:   "error",
	StatusClean:   "clean",
	StatusDirty:   "dirty",
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
	Code    StatusCode
	Message string
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
