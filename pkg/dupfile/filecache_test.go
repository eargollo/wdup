package dupfile

import (
	"reflect"
	"testing"
)

func TestNewMD5Cache(t *testing.T) {
	type args struct {
		file string
	}
	tests := []struct {
		name    string
		args    args
		wantNil bool
		wantErr bool
	}{
		{name: "Invalid file", args: args{""}, wantNil: true, wantErr: true},
		{name: "Valid path", args: args{t.TempDir()}, wantNil: false, wantErr: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewFileCache(tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMD5Cache() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (got == nil) != tt.wantNil {
				t.Errorf("NewMD5Cache() = %v, wantNil %v", got, tt.wantNil)
			}
		})
	}
}

func TestCache(t *testing.T) {
	c, err := NewFileCache(t.TempDir())
	if err != nil {
		t.Fatalf("Could not create cache: %v", err)
	}
	defer c.Close()

	aFile := &File{Path: "/my/file", Name: "file", Size: 100, Hash: []byte("abc")}
	bFile := &File{Path: "/my/other/file", Name: "file", Size: 200, Hash: []byte("efg")}
	c.Put(aFile)
	c.Put(bFile)

	res := c.Get("/not/in/there", 50)
	if res != nil {
		t.Errorf("should return nil when file does not exist")
	}

	res = c.Get("/my/file", 100)
	if !reflect.DeepEqual(res, aFile.Hash) {
		t.Errorf("Get() = %v, want %v", res, aFile)
	}

	res = c.Get("/my/other/file", 200)
	if !reflect.DeepEqual(res, bFile.Hash) {
		t.Errorf("Get() = %v, want %v", res, bFile)
	}

	res = c.Get("/my/file", 101)
	if res != nil {
		t.Error("should get nil if size doesn't match")
	}

	res = c.Get("/my/file", 100)
	if res != nil {
		t.Error("should get nil after a mismatch")
	}
}

func TestMD5Cache_List(t *testing.T) {
	seed := []string{
		"/a/one",
		"/a/two",
		"/a/three",
		"/b/one",
		"/b/two/one",
		"/b/two/two",
	}
	c, err := NewFileCache(t.TempDir())
	if err != nil {
		t.Fatalf("Could not create cache: %v", err)
	}
	defer c.Close()

	// Add records
	for _, name := range seed {
		c.Put(&File{Path: name})
	}

	tests := []struct {
		name  string
		paths []string
		want  []string
	}{
		{name: "All no parameter", paths: []string{}, want: seed},
		{name: "All explicit", paths: []string{"/"}, want: seed},
		{name: "Filter single path", paths: []string{"/b"}, want: seed[3:]},
		{name: "No match", paths: []string{"b"}, want: []string{}},
		{name: "Or match", paths: []string{"/a/t", "/b/o"}, want: seed[1:4]},
		{name: "Or match no duplicate", paths: []string{"/a/t", "/a/t"}, want: seed[1:3]},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := c.List(tt.paths); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("MD5Cache.List() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMD5Cache_Delete(t *testing.T) {
	seed := []string{
		"/a/one",
		"/a/two",
		"/a/three",
		"/b/one",
		"/b/two/one",
		"/b/two/two",
	}
	c, err := NewFileCache(t.TempDir())
	if err != nil {
		t.Fatalf("Could not create cache: %v", err)
	}
	defer c.Close()

	// Add records
	for _, name := range seed {
		c.Put(&File{Path: name})
	}

	c.Delete([]string{})
	res := c.List([]string{})
	if !reflect.DeepEqual(res, seed) {
		t.Errorf("MD5Cache.Delete() = %v, want %v", res, seed)
	}

	c.Delete([]string{"/a", "/a/one", "/b/two/two"})
	res = c.List([]string{})
	if !reflect.DeepEqual(res, seed[1:5]) {
		t.Errorf("MD5Cache.Delete() = %v, want %v", res, seed[1:5])
	}
}
