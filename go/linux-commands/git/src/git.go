package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	fp "path/filepath"
	"syscall"
	"time"
)

func OpenFile(filepath string) (*os.File, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return file, err
	} else {
		return file, nil
	}
}

type FileDetails struct {
	Name      string
	Path      string
	Size      int64
	Mode      os.FileMode
	ModTime   string
	IsDir     bool
	IsSymlink bool
	Stat      *syscall.Stat_t
}

func (fd *FileDetails) Init(filepath string) error {
	path, err := fp.Abs(filepath)
	if err != nil {
		return err
	}
	file, err := os.Stat(filepath)
	if err != nil {
		return err
	}
	fd.Name = file.Name()
	fd.Path = path
	fd.Size = file.Size()
	fd.Mode = file.Mode()
	fd.ModTime = file.ModTime().String()
	fd.IsDir = file.IsDir()
	fd.IsSymlink = file.Mode()&os.ModeSymlink != 0
	fd.Stat = file.Sys().(*syscall.Stat_t)
	return nil
}

type GitObject interface {
	GetType() string
	GetHash() []byte
	GenHash() []byte
}

type GitBlob struct {
	hash    []byte
	size    int64
	name    string
	details *FileDetails
	content *[]byte
}

func (b *GitBlob) GetType() string {
	return "blob"
}

func (b *GitBlob) GetHash() []byte {
	return b.hash
}

func (b *GitBlob) GenHash() []byte {
	hash := sha256.New()
	header := fmt.Sprintf("blob %d\000", b.size)
	hash.Write([]byte(header))
	hash.Write(*b.content)
	hash_val := hash.Sum(nil)
	b.hash = hash_val
	return hash_val
}

func (b *GitBlob) Init(filepath string) error {
	var Content []byte
	var err error
	file, ferr := os.Open(filepath)
	if ferr != nil {
		return ferr
	}
	defer file.Close()
	b.hash, err = HashFile(file, 8192)
	if err != nil {
		return err
	}
	b.details = &FileDetails{}
	b.details.Init(filepath)
	b.name = filepath
	b.size = b.details.Size
	Content, err = ReadFile(filepath)
	if err != nil {
		return err
	}
	b.content = &Content
	return nil
}

type GitTree struct {
	entryCount int
	entries    []GitTreeEntry
	hash       []byte
	name       string
	size       int
}

func (t *GitTree) GetType() string {
	return "tree"
}

func (t *GitTree) GetHash() []byte {
	return t.hash
}

func (t *GitTree) GenHash() []byte {
	hash := sha256.New()
	hash.Write([]byte(fmt.Sprintf("tree %d\000", t.size)))
	for _, entry := range t.entries {
		hash.Write([]byte(fmt.Sprintf("%d %s %s\000", entry.mode, entry.name, entry.hash)))
	}
	hash_sum := hash.Sum(nil)
	t.hash = hash_sum
	return hash_sum
}

func (t *GitTree) AddEntry(entry GitTreeEntry) {
	t.entries = append(t.entries, entry)
}

type GitTreeEntry struct {
	mode      int
	entryType string
	hash      []byte
	name      string
	object    GitObject
}

func NewGitBlob(filepath string) (*GitBlob, error) {
	blob := &GitBlob{}
	err := blob.Init(filepath)
	if err != nil {
		return nil, err
	}
	return blob, nil
}

func NewGitTreeEntry(filepath string) (*GitTreeEntry, error) {
	entry := &GitTreeEntry{}
	//var err error
	file_details, fderr := NewFileDetailsPtr(filepath)
	if fderr != nil {
		return entry, fderr
	}
	entry.mode = int(file_details.Mode)
	entry.name = file_details.Name
	entry.entryType = entry.object.GetType()
	if !file_details.IsDir { // Open file and hash if if not a directory
		bObj, berr := NewGitBlob(filepath)
		if berr != nil {
			return entry, berr
		} else {
			entry.object = bObj
			entry.hash = entry.object.GenHash()
		}
		file, ferr := os.Open(filepath)
		if ferr != nil {
			return entry, ferr
		}
		defer file.Close()
		_ = bObj.GenHash()

	} else {
		tree, terr := NewGitTree(fp.Join(filepath, file_details.Name))
		if terr != nil {
			return entry, terr
		} else {
			entry.object = tree
			entry.hash = entry.object.GenHash()
		}
	}
	entry.name = file_details.Name
	entry.mode = int(file_details.Mode)
	if file_details.IsDir {
		entry.entryType = "tree"
	} else {
		entry.entryType = "blob"
	}

	return entry, nil
}

func NewGitTree(filepath string) (*GitTree, error) {
	//var err error
	tree := &GitTree{}
	tree.entries = []GitTreeEntry{}
	dirEntries, derr := os.ReadDir(filepath)
	if derr != nil {
		return tree, derr
	}
	for _, entry := range dirEntries {
		entryInfo, err := NewFileDetailsPtr(fp.Join(filepath, entry.Name()))
		if err != nil {
			return tree, err
		}
		treeEntry, terr := NewGitTreeEntry(fp.Join(filepath, entryInfo.Name))
		if terr != nil {
			return tree, terr
		}
		entries := append(tree.entries, *treeEntry)
		tree.entries = entries
		}
	}
	return tree, nil
}

type GitCommit struct {
	treeHash      []byte
	parentCount   int
	parentHashes  *[][]byte
	authorName    string
	authorEmail   string
	authorTime    time.Time
	comitterName  string
	comitterEmail string
	comitterTime  time.Time
	message       string
}

func NewFileDetailsPtr(filepath string) (*FileDetails, error) {
	fd := &FileDetails{}
	err := fd.Init(filepath)
	if err != nil {
		return nil, err
	}
	return fd, nil
}

func NewFileDetails(filepath string) (FileDetails, error) {
	fd := &FileDetails{}
	err := fd.Init(filepath)
	if err != nil {
		return FileDetails{}, err
	}
	return *fd, nil
}

func logExit(err error, exit_code int) {
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(exit_code)
	}
}

func GetFileDetails(dirpath string) ([]FileDetails, error) {
	var files []FileDetails
	if _, err := os.Stat(dirpath); os.IsNotExist(err) {
		return nil, err
	}
	if stat, err := os.Stat(dirpath); err != nil {
		logExit(err, 2)
	} else if !stat.IsDir() {
		return nil, fmt.Errorf("'%s' is not a directory", dirpath)
	} else {

		files_in_dir, err := os.ReadDir(dirpath)
		if err != nil {
			logExit(err, 1)
		}
		for _, file := range files_in_dir {
			fd := &FileDetails{}
			fd.Init(fp.Join(dirpath, file.Name()))
			files = append(files, *fd)
		}

	}
	return files, nil
}

func HashFile(f *os.File, chunkSize int) ([]byte, error) {
	hash := sha256.New()
	buf := make([]byte, chunkSize)

	for {
		bytesRead, err := f.Read(buf)
		if err != nil && err != io.EOF {
			return buf, err
		}
		if bytesRead > 0 {
			hash.Write(buf[:bytesRead])
		} else {
			break
		}
	}
	hash_bytes := hash.Sum(nil)
	return hash_bytes, nil
}

func HashToString(h []byte) string {
	return hex.EncodeToString(h)
}

func HashFileToString(filepath string, chunkSize int) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close() // Close the file when the function returns
	hash, err := HashFile(file, chunkSize)
	if err != nil {
		return "", err
	}
	return HashToString(hash), nil
}

func ReadFile(filepath string) ([]byte, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	fileContent, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return fileContent, nil
}

func main() {
	filePath := "/home/jharmon/.bashrc"
	// Open the file
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Create a new SHA256 hash
	hash := sha256.New()

	// Copy the file data into the hash
	if _, err := io.Copy(hash, file); err != nil {
		log.Fatal(err)
	}

	// Get the final hash sum
	hashInBytes := hash.Sum(nil)                   // returns []byte
	hashInString := fmt.Sprintf("%x", hashInBytes) // returns string

	fmt.Println("SHA-256 hash:", hashInString)
}
