package dupfile

type DedupOption func(*Dedup) error

func WithPaths(paths []string) DedupOption {
	return func(d *Dedup) error {
		d.paths = paths
		return nil
	}
}

func WithCache(path string) DedupOption {
	return func(d *Dedup) error {
		if path != "" {
			var err error
			d.cache, err = NewFileCache(path)
			if err != nil {
				return err
			}
		}

		return nil
	}
}

type ObservableEvent struct {
	Type        string     `json:"type"`
	Description string     `json:"description"`
	Stats       DedupStats `json:"stats"`
}

type Observer func(event *ObservableEvent)

func WithObserver(observer Observer) DedupOption {
	return func(d *Dedup) error {
		d.observer = observer
		return nil
	}
}
