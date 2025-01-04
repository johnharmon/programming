package main

import (
	"os"
	"syscall"
	"time"
)

type GitTree struct {
	entryCount int
	header     *[]byte
	entries    []GitTreeEntry
	hash       []byte
	name       string
	size       int
	content    *[]byte
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

type GitObject interface {
	GetType() string
	GetHash() []byte
	GenHash() []byte
	//SetContent(*[]byte) error
}

type GitBlob struct {
	hash    []byte
	size    int64
	name    string
	details *FileDetails
	content *[]byte
}

type GitCommit struct {
	header        string
	hash          []byte
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

type GitTreeEntry struct {
	mode      int
	entryType string
	hash      []byte
	name      string
	object    GitObject
}

type IpHeader struct {
	version        [4]byte
	headerLength   [3]byte
	typeOfService  [8]byte
	totalSize      [16]byte
	identification [16]byte
	flags          [3]byte
	fragmentOffset [13]byte
	timeToLive     [8]byte
	protocol       [8]byte
	headerChecksum [16]byte
	sourceIp       [32]byte
	destinationIp  [32]byte
	options        [32]byte
}
