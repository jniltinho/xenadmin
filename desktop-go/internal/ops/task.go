package ops

import (
	"errors"
	"fmt"
	"time"

	xenapi "go/xenapi"
)

// pollEvery is how often a task is polled (AsyncAction.PollToCompletion polls every second too); pollRetries is how
// many failed polls in a row (a dropped connection, a coordinator changing) are tolerated before giving up.
var (
	pollEvery   = time.Second
	pollRetries = 60
)

// ErrUncertain: the task could not be followed to its end, so it may still be running or may have succeeded; never
// undo or repeat the operation on this error (the caller tells the user to check the server).
var ErrUncertain = errors.New("the outcome is unknown: the server stopped answering while the operation ran; check it before trying again")

// wait is AsyncAction.PollToCompletion: it runs an XAPI call as a task and polls it until it ends, so a long
// operation (an evacuation, a migration) is never cut by the HTTP client's timeout; a failed task's error_info
// becomes an *xenapi.APIError, as a synchronous call's. The task is destroyed afterwards.
// Use: wait(s)(xenapi.VM.AsyncStart(s, vm, false, false)).
func wait(s *xenapi.Session) func(xenapi.TaskRef, error) error {
	return func(task xenapi.TaskRef, err error) error { return poll(s, task, err) }
}

func poll(s *xenapi.Session, task xenapi.TaskRef, err error) error {
	if err != nil {
		return err
	}
	for failed := 0; ; {
		status, err := xenapi.Task.GetStatus(s, task)
		if err != nil {
			if failed++; failed >= pollRetries {
				return fmt.Errorf("%w (%v)", ErrUncertain, err) // the task is left on the server: it may still run
			}
			time.Sleep(pollEvery)
			continue
		}
		failed = 0
		if status != xenapi.TaskStatusTypePending && status != xenapi.TaskStatusTypeCancelling {
			defer xenapi.Task.Destroy(s, task) //nolint:errcheck // best effort, as the C# app
		}
		switch status {
		case xenapi.TaskStatusTypeSuccess:
			return nil
		case xenapi.TaskStatusTypeFailure:
			info, err := xenapi.Task.GetErrorInfo(s, task)
			if err != nil {
				return err
			}
			if len(info) == 0 {
				return errors.New("the operation failed")
			}
			return &xenapi.APIError{Message: info[0], Data: info[1:]}
		case xenapi.TaskStatusTypeCancelled:
			return errors.New("the operation was cancelled")
		}
		time.Sleep(pollEvery)
	}
}
