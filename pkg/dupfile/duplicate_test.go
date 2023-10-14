package dupfile

import (
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	tempDir := t.TempDir()
	tmpCache, _ := NewFileCache(tempDir)
	tmpCache.Close()

	tests := []struct {
		name    string
		opts    []DedupOption
		want    *Dedup
		wantErr bool
	}{
		{name: "Defaults", opts: []DedupOption{}, want: &Dedup{paths: []string{"."},
			fileBySize: map[int64][]*File{}, duplicates: map[string][]*File{},
		}, wantErr: false},
		{
			name: "With paths",
			opts: []DedupOption{WithPaths([]string{"/Home", "/Volumes"})},
			want: &Dedup{paths: []string{"/Home", "/Volumes"}, fileBySize: map[int64][]*File{},
				duplicates: map[string][]*File{}},
			wantErr: false,
		},
		{
			name: "With cache",
			opts: []DedupOption{
				WithPaths([]string{"/Home", "/Volumes"}),
				WithCache(tempDir),
			},
			want: &Dedup{paths: []string{"/Home", "/Volumes"}, cache: tmpCache,
				fileBySize: map[int64][]*File{}, duplicates: map[string][]*File{}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.opts...)
			defer got.Close()
			if (err != nil) != tt.wantErr {
				t.Errorf("err = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.want.cache != nil {
				if got.cache == nil {
					t.Errorf("New() = %v, want %v", got.cache, tt.want.cache)
				}
				tt.want.cache = got.cache
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
