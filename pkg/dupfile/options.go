package dupfile

type DedupOption func(*Dedup)

func WithPaths(paths []string) DedupOption {
	return func(d *Dedup) {
		d.paths = paths
	}
}

func WithCache(path string) DedupOption {
	return func(d *Dedup) {
		d.cachePath = path
	}
}
