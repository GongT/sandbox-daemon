package instance

import "gitlab.com/tozd/go/errors"

func (pm *D) Destroy() error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	errs := []error{}
	for _, part := range pm.parts {
		err := part.Stop()
		if err != nil {
			errs = append(errs, err)
		}
	}

	close(pm.quitChan)

	return errors.Join(errs...)
}

func (d *D) Join() {
	<-d.quitChan
}
