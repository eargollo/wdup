package dupfile

import (
	"bytes"
	"container/heap"
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"sort"
	"sync"
	"time"
)

type DedupStats struct {
	Files     int           `json:"files"`
	Checked   int           `json:"checked"`
	CacheHits int           `json:"cacheHits"`
	Started   time.Time     `json:"started"`
	Finished  time.Time     `json:"finished"`
	Duration  time.Duration `json:"duration"`
	Listing   bool          `json:"listing"`
}

type Dedup struct {
	paths      []string
	cache      *FileCache
	mu         sync.Mutex
	pq         PriorityQueue
	fileBySize map[int64][]*File
	duplicates map[string][]*File
	current    *File
	observer   Observer
	stats      DedupStats
}

func New(opts ...DedupOption) (*Dedup, error) {
	// Default
	d := Dedup{paths: []string{"."}}
	d.fileBySize = make(map[int64][]*File)
	d.duplicates = make(map[string][]*File)

	for _, opt := range opts {
		err := opt(&d)
		if err != nil {
			return &d, err
		}
	}

	return &d, nil
}

func (d *Dedup) Close() {
	if d.cache != nil {
		d.cache.Close()
	}
}

func (d *Dedup) observe(tp string, description string) {
	if d.observer != nil {
		event := &ObservableEvent{
			Type:        tp,
			Description: description,
			Stats:       d.stats,
		}
		d.observer(event)
	}
}

func (d *Dedup) Run() [][]*File {

	var prod sync.WaitGroup
	//Producer
	d.stats.Started = time.Now()
	d.stats.Listing = true
	for _, path := range d.paths {
		walkPath := path
		prod.Add(1)
		go func() {
			defer prod.Done()
			d.scanFolder(walkPath)
		}()
	}

	// Consumer
	var cons sync.WaitGroup
	done := make(chan bool)
	cons.Add(1)
	go func() {
		prodFinished := false
		defer cons.Done()
		for {
			select {
			case <-done:
				d.observe("ListFiles", "Finished listing all files")
				prodFinished = true
			default:
				ok := d.validateNext()
				if !ok && prodFinished {
					d.observe("FindDuplicates", "Finished finding duplicates")
					return
				}
			}
		}
	}()

	// Tracker
	var track sync.WaitGroup
	finished := make(chan bool)
	track.Add(1)
	go func() {
		defer track.Done()
		start := time.Now()
		// last := time.Now()
		for {
			select {
			case <-finished:
				d.stats.Duration = time.Since(start)
				d.observe("Progress",
					fmt.Sprintf("Processed %d files in %s time",
						d.stats.Checked,
						time.Since(start)))
				return
			default:
				// if time.Since(last) > (10 * time.Second) {

				// last = time.Now()
				total := d.stats.Files
				count := d.stats.Checked

				perc := 0
				if total > 0 {
					perc = count * 10000 / total
				}
				percR := float32(perc) / 100
				file := d.current
				message := ""
				if d.stats.Listing {
					message = "*"
				}
				if file == nil {
					message = fmt.Sprintf("Done %d of %d%s files [%.2f%%]", count, total, message, percR)
				} else {
					message = fmt.Sprintf("Done %d of %d%s files [%.2f%%] (%s %s) ", count, total, message, percR, file.Name, ByteCountIEC(file.Size))
				}
				d.observe("Progress", message)
				time.Sleep(250 * time.Millisecond)
				// }
			}
		}
	}()

	prod.Wait()
	log.Printf("Finished listing files")
	// Signal producers are finished
	done <- true
	d.stats.Listing = false
	cons.Wait()
	log.Printf("Finished finding duplicates")
	// Signal tracker
	finished <- true
	track.Wait()
	d.stats.Finished = time.Now()
	d.stats.Duration = d.stats.Finished.Sub(d.stats.Started)
	d.observe("Progress", fmt.Sprintf("Finished! Files: %d Duration %s", d.stats.Checked, d.stats.Duration))
	// Sort result
	return d.duplicatedList()
}

func (d *Dedup) scanFolder(path string) {
	log.Printf("Starting walk '%s'...", path)
	err := filepath.Walk(path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			d.observe("Error", fmt.Sprintf("Error scanning path %s: %v", path, err))
			log.Printf("ERROR: Walk error: %v", err)
			return nil
		}
		if info.IsDir() {
			return nil
		}

		d.push(File{Path: path, Name: info.Name(), Size: info.Size()})

		return nil
	})

	if err != nil {
		log.Printf("Scan error: %v", err)
	}
	log.Printf("Done walking '%s'", path)
}

func (d *Dedup) push(file File) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.stats.Files++

	heap.Push(&d.pq, &file)
}

func (d *Dedup) pop() *File {
	d.mu.Lock()
	defer d.mu.Unlock()

	if len(d.pq) == 0 {
		return nil
	}
	d.stats.Checked++
	return heap.Pop(&d.pq).(*File)
}

func (d *Dedup) validateNext() bool {
	file := d.pop()
	d.current = file

	if file == nil {
		return false
	}

	// Verify if there is a file with the same size
	fByS, ok := d.fileBySize[file.Size]
	if !ok {
		// No file with the same size, no duplicate
		d.fileBySize[file.Size] = []*File{file}
		return true
	}

	// There is a file with the same size, compare MD5s
	fileHash, err := file.CryptoHash(d.cache)
	if err != nil {
		d.observe("Error", fmt.Sprintf("Error calculating MD5 for %s: %v", file.AbsPath(), err))
		log.Printf("Error: Could not calculate MD5 for '%s': %v", file.AbsPath(), err)
		return true
	}

	for _, comp := range fByS {
		// compare md5 with existing files with the same size
		compHash, err := comp.CryptoHash(d.cache)
		if err != nil {
			d.observe("Error", fmt.Sprintf("Error calculating MD5 for %s: %v", file.AbsPath(), err))
			log.Printf("Error: Could not calculate MD5 for '%s': %v", comp.AbsPath(), err)
			continue
		}

		if bytes.Equal(fileHash, compHash) {
			// Found a duplicate
			log.Printf("Duplicate (%d bytes): %s [%x] - %s [%x]", file.Size, file.AbsPath(), fileHash, comp.AbsPath(), compHash)
			// Add to duplicates
			key := string(fileHash[:])
			dup, ok := d.duplicates[key]
			if !ok {
				d.duplicates[key] = []*File{comp, file}
			} else {
				dup = append(dup, file)
				d.duplicates[key] = dup
			}
			break
		}
	}
	// add to file by size group
	fByS = append(fByS, file)
	d.fileBySize[file.Size] = fByS

	return true
}

func (d *Dedup) duplicatedList() [][]*File {
	list := [][]*File{}
	for _, dup := range d.duplicates {
		sort.Slice(dup, func(i int, j int) bool {
			return dup[i].Path < dup[j].Path
		})
		list = append(list, dup)
	}
	sort.Slice(list, func(i int, j int) bool {
		return len(list[i])*int(list[i][0].Size) > len(list[j])*int(list[j][0].Size)
	})
	return list
}

func ByteCountIEC(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB",
		float64(b)/float64(div), "KMGTPE"[exp])
}
