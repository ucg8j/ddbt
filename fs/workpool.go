package fs

import (
	"fmt"
	"sync"

	"ddbt/config"
	"ddbt/utils"
)

// Process the given file list through `f`. If a progress bar is given, then it will show the stauts line as we go
func ProcessFiles(files []*File, f func(file *File) error, pb *utils.ProgressBar) error {
	var wait sync.WaitGroup
	var errMutex sync.RWMutex
	var firstError error

	c := make(chan *File, len(files))

	worker := func() {
		var statusRow *utils.StatusRow

		if pb != nil {
			statusRow = pb.NewStatusRow()
		}

		for file := range c {
			errMutex.RLock()
			if firstError != nil {
				// If we're already had an error, skip through the rest of the items
				errMutex.RUnlock()
				wait.Done()
				continue
			}
			errMutex.RUnlock()

			if statusRow != nil {
				statusRow.Update(fmt.Sprintf("Running %s", file.Name))
			}

			err := f(file)
			wait.Done()

			if statusRow != nil {
				statusRow.SetIdle()
			}
			if err != nil {
				errMutex.Lock()
				if firstError == nil {
					firstError = err
				}
				errMutex.Unlock()

				return
			}

		}
	}

	for i := 0; i < config.NumberThreads(); i++ {
		go worker()
	}

	wait.Add(len(files))
	for _, file := range files {
		c <- file
	}

	wait.Wait()

	return firstError
}
