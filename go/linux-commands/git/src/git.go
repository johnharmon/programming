package main

import (
	"bytes"
	"compress/zlib"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"
	fp "path/filepath"
	"strconv"
	"syscall"
)

func OpenFile(filepath string) (*os.File, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return file, err
	} else {
		return file, nil
	}
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

func NewGitBlobFromContent(size int, content []byte) *GitBlob {
	blob := &GitBlob{}
	blob.content = &content
	blob.size = int64(size)
	return blob

}

func (*GitTree) ProcessTreeEntries() error {
	return nil
}

func NewGitTreeFromContent(content []byte, hash []byte, header []byte) *GitTree {
	tree := &GitTree{}
	tree.hash = hash
	tree.content = &content
	tree.entries = []GitTreeEntry{}
	return tree
}

func NewGitCommitFromContent(content []byte) *GitCommit {
	commit := &GitCommit{}
	return commit
}

func (c *GitCommit) GetType() string {
	return "commit"
}

func (c *GitCommit) GetHash() []byte {
	return c.hash
}

func (c *GitCommit) GenHash() []byte {
	hash := sha256.New()
	hash.Write([]byte(fmt.Sprintf("commit %d\000", len(c.header))))
	hash.Write([]byte(c.header))
	hash.Write([]byte(fmt.Sprintf("tree %s\000", c.treeHash)))
	for _, parent := range *c.parentHashes {
		hash.Write([]byte(fmt.Sprintf("parent %s\000", parent)))
	}
	hash.Write([]byte(fmt.Sprintf("author %s <%s> %s\000", c.authorName, c.authorEmail, c.authorTime)))
	hash.Write([]byte(fmt.Sprintf("committer %s <%s> %s\000", c.comitterName, c.comitterEmail, c.comitterTime)))
	hash.Write([]byte(c.message))
	hash_sum := hash.Sum(nil)
	c.hash = hash_sum
	return hash_sum
}

func CompressBytes(content []byte) ([]byte, error) {
	var compressed bytes.Buffer
	w := zlib.NewWriter(&compressed)
	w.Write(content)
	defer w.Close()
	err := w.Flush()
	if err != nil {
		return nil, err
	}
	fmt.Printf("Compressed bytes: %s\n", compressed.String())
	return compressed.Bytes(), nil
}

func (b *GitBlob) WriteBlobToFile() error {
	hash := b.GenHash()
	dirmode := os.FileMode(0755)
	hash_prefix := hex.EncodeToString(hash[0:2])
	hash_suffix := hex.EncodeToString(hash[2:])
	fmt.Printf("Hash: %s\n", hex.EncodeToString(hash))
	hash_directory := fp.Join(".gitg", "objects", hash_prefix)
	os.MkdirAll(hash_directory, dirmode)

	hash_path := fp.Join(".gitg", "objects", hash_prefix, hash_suffix)
	fmt.Printf("Hash path: %s\n", hash_path)
	object_file, err := os.Create(hash_path)
	if err != nil {
		return err
	}
	defer object_file.Close()
	header := fmt.Sprintf("blob %d\000", b.size)
	object_file.Write([]byte(header))
	compressed_content, err := CompressBytes(*b.content)
	if err != nil {
		return err
	}
	object_file.Write(compressed_content)
	return nil
}

func ReadObjectFromFile(hash []byte) (GitObject, error) {
	hash_prefix := hex.EncodeToString(hash[0:2])
	hash_file := hex.EncodeToString(hash[2:])
	hash_path := fp.Join(".gitg", "objects", hash_prefix, hash_file)
	object_file, err := os.Open(hash_path)
	if err != nil {
		return nil, err
	}
	defer object_file.Close()
	const chunkSize = 8192
	chunk := make([]byte, chunkSize)
	var (
		bytesRead int
		buffer    bytes.Buffer
	)
	for {
		bytesRead, err = object_file.Read(chunk)
		if err != nil && err != io.EOF {
			return nil, err
		}
		if bytesRead > 0 {
			buffer.Write(chunk[:bytesRead])
		} else {
			break
		}
	}
	content := buffer.Bytes()
	headersize := bytes.Index(content, []byte("\000")) + 1
	header := content[:headersize]
	type_location := bytes.Index(header, []byte(" "))
	objectType := string(header[:type_location])
	content = content[headersize:]
	if objectType == "tree" {
		tree := NewGitTreeFromContent(content, hash, header)
		return tree, nil
	} else if objectType == "blob" {
		sblobSize := string(header[type_location+1 : headersize-1])
		blobSize, err := strconv.Atoi(sblobSize)
		if err != nil {
			return nil, err
		}
		blob := NewGitBlobFromContent(blobSize, content)
		return blob, nil
	} else if objectType == "commit" {
		commit := NewGitCommitFromContent(content)
		return commit, nil

	}
	return nil, nil
}

// func FetchObjectFromHash(hash []byte) (GitObject, error) {
// 	var new_object GitObject
// 	hash_prefix := hash[0:2]
// 	hash_path := fp.Join("\.gitg", "objects", string(hash_prefix), string(hash[2:]))
// 	object_file, err := os.Open(hash_path)
// 	if err != nil {
// 		return nil, err
// 	}
// 	defer object_file.Close()
// 	object_content, err := io.ReadAll(object_file)
// 	if err != nil {
// 		return nil, err
// 	}
// 	headerSize := bytes.Index(object_content, []byte("\000")) + 1
// 	header := string(object_content[:headerSize])
// 	content := object_content[headerSize:]
// 	var objectType string
// 	var object_size int
// 	_, err = fmt.Sscanf(header, "%s %d", &objectType, object_size)
// 	if err != nil {
// 		return nil, err
// 	} else {
// 		switch objectType {
// 		case "blob":
// 			new_object = NewGitBlobFromContent(object_size, content)
// 		case "tree":
// 			new_object = NewGitTreeFromContent(content)
// 		case "commit":
// 			new_object = NewGitCommitFromContent(content)

// 		}
// 	}
// 	return new_object, nil
// }

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

func (t *GitTree) ProcessEntries() error {
	return nil
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
	return tree, nil
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

func ReadWriteFile(sourcepath string, targetpath string, chunkSize int) error {
	sourcefile, err := os.Open(sourcepath)
	if err != nil {
		return err
	}
	defer sourcefile.Close()
	sourceFileDetails, err := NewFileDetailsPtr(sourcepath)
	if err != nil {
		return err
	}
	fmt.Printf("Source file details: %v\n", *sourceFileDetails)
	targetfile, err := os.Create(targetpath)
	if err != nil {
		return err
	}
	defer targetfile.Close()
	chunk := make([]byte, chunkSize)
	for {
		bytesRead, err := sourcefile.Read(chunk)
		if err != nil && err != io.EOF {
			return err
		}
		fmt.Printf("Bytes read: %d\n", bytesRead)
		if bytesRead > 0 {
			fmt.Printf("Writing %d bytes\n", bytesRead)
			targetfile.Write(chunk[:bytesRead])
		} else {
			break
		}
	}
	return nil
}

func main() {
	//var nFlag = flag.Int("n", 10, "Number of lines to read")
	var fFlag = flag.String("f", "", "File to read")
	var wFlag = flag.String("w", "output.txt", "File to write to")
	flag.Parse()
	fmt.Printf("%v\n", *fFlag)
	fmt.Printf("%v\n", *wFlag)

	blob, err := NewGitBlob(*fFlag)
	if err != nil {
		logExit(err, 1)
	}
	blob.WriteBlobToFile()
	fmt.Printf("Blob: %v\n", *blob)

	// err := ReadWriteFile(*fFlag, *wFlag, 8192)
	// if err != nil {
	// 	logExit(err, 1)
	// } else {
	// 	fmt.Printf("File %s written to %s\n", *fFlag, *wFlag)
	// 	os.Exit(0)
	// }

	//fmt.Printf("%v\n", *nFlag)
	// filePath := "/home/jharmon/.bashrc"
	// Open the file
	// file, err := os.Open(filePath)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer file.Close()

	// // Create a new SHA256 hash
	// hash := sha256.New()

	// // Copy the file data into the hash
	// if _, err := io.Copy(hash, file); err != nil {
	// 	log.Fatal(err)
	// }

	// // Get the final hash sum
	// hashInBytes := hash.Sum(nil)                   // returns []byte
	// hashInString := fmt.Sprintf("%x", hashInBytes) // returns string

	// fmt.Println("SHA-256 hash:", hashInString)
}
