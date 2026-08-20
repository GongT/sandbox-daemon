package instance

import "gitlab.com/tozd/go/errors"

func (d *D) Destroy() error {
	d.mu.Lock()
	defer d.mu.Unlock()

	errs := []error{}
	for _, part := range d.parts {
		err := part.Close()
		if err != nil {
			errs = append(errs, err)
		}
	}

	close(d.quitChan)

	return errors.Join(errs...)
}

func (d *D) Join() <-chan struct{} {
	return d.quitChan
}
