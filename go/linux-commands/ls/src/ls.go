package main

import (
	"fmt"
	"os"
	"os/user"
	"reflect"
	"strconv"
	"syscall"

	"golang.org/x/term"
)

func getFileStructs(file os.FileInfo) (*syscall.Stat_t, os.FileMode) {
	file_stat := file.Sys().(*syscall.Stat_t)
	file_mode := file.Mode()
	return file_stat, file_mode
}

func getNamesFromIds(uid int, gid int) (string, string) {
	userName, err := user.LookupId(strconv.Itoa(uid))
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(3)
	}
	groupName, err := user.LookupGroupId(strconv.Itoa(gid))
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(3)
	}
	return userName.Username, groupName.Name
}

func getFileIds(file os.FileInfo) (uint32, uint32) {
	file_stat := file.Sys().(*syscall.Stat_t)
	uid := file_stat.Uid
	gid := file_stat.Gid
	return uid, gid
}

func getFileFromDirEntry(file os.DirEntry) (os.FileInfo, *syscall.Stat_t, os.FileMode) {
	file_info, err := file.Info()
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		os.Exit(3)
	}
	file_stat, file_mode := getFileStructs(file_info)
	return file_info, file_stat, file_mode
}

func getFilePerm(mode os.FileMode) {
	//perm_typs := []string{"---", "--x", "-w-", "-wx", "r--", "r-x", "rw-", "rwx"}
	t := reflect.TypeOf(mode)
	fmt.Printf("Value: %v, Type: %v, Kind: %v\n", mode, t, t.Kind())
}

func processDir(targetdir string, sep_string string, stat_fileInfo os.FileInfo, stat *syscall.Stat_t, stat_mode os.FileMode) {
	files_in_dir, err := os.ReadDir(targetdir)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(3)
	}
	var longest_name int = 0
	var longest_uid int = 0
	var longest_gid int = 0
	var longest_size int = 0
	for _, file := range files_in_dir {
		if len(file.Name()) > longest_name {
			longest_name = len(file.Name())
		}
		file_info, err := file.Info()
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(4)
		} else {
			file_stat_sys := file_info.Sys().(*syscall.Stat_t)
			if len(fmt.Sprintf("%d", file_stat_sys.Uid)) > longest_uid {
				longest_uid = len(fmt.Sprintf("%d", file_stat_sys.Uid))
			}
			if len(fmt.Sprintf("%d", file_stat_sys.Gid)) > longest_gid {
				longest_gid = len(fmt.Sprintf("%d", file_stat_sys.Gid))
			}
			if len(fmt.Sprintf("%d", file_info.Size())) > longest_size {
				longest_size = len(fmt.Sprintf("%d", file_info.Size()))
			}
		}

	}
	format_string_header := fmt.Sprintf("%%-12s |  %%-%ds  |  %%-%ds  | %%-%ds  |  %%-%ds  |  %%-s  \n\n", longest_name, longest_uid, longest_gid, longest_size)
	format_string_entries := fmt.Sprintf("%%-12v |  %%-%dv  |  %%-%dv  | %%-%dv  |  %%-%dv  |  %%-v  \n", longest_name, longest_uid, longest_gid, longest_size)
	for idx, file := range files_in_dir {
		if idx == 0 {

			fmt.Printf(format_string_header, "Permissions", "Name", "UID", "GID", "Size", "Last Modified")
			fmt.Printf(format_string_entries, stat_mode.Perm(), stat_fileInfo.Name(), stat.Uid, stat.Gid, stat_fileInfo.Size(), stat_fileInfo.ModTime())
			fmt.Printf("%s\n", sep_string)
		}
		file_info, err := file.Info()
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(3)
		} else {
			file_stat_sys := file_info.Sys().(*syscall.Stat_t)
			file_mode := file_info.Mode()
			fmt.Printf(format_string_entries, file_mode.Perm(), file_info.Name(), file_stat_sys.Uid, file_stat_sys.Gid, file_info.Size(), file_info.ModTime())

		}
	}
}

func processDirEntries(files_in_dir []os.DirEntry, sep_string string, stat_fileInfo os.FileInfo, stat *syscall.Stat_t, stat_mode os.FileMode) { // Unsused optional arguments, to be used later?

	// Set variables for longest entries for header types of unknown length (name, username, groupname, size)
	var longest_name int = 0
	var longest_username int = 0
	var longest_groupname int = 0
	var longest_size int = 0

	// Iterate over files in directory to determine longest name, username, groupname, and size
	for _, file := range files_in_dir {
		if len(file.Name()) > longest_name {
			longest_name = len(file.Name())
		}
		file_info, err := file.Info()
		if err != nil {
			fmt.Printf("%v\n", err)
			os.Exit(4)
		}

		// Direct syscall.Stat_t type conversion from file_info.Sys() method, needed specifically for uid and gid
		file_stat_sys := file_info.Sys().(*syscall.Stat_t)
		uid := file_stat_sys.Uid
		gid := file_stat_sys.Gid

		// Set username, convert uid to string, user user.LookupId to get a username from an id, gives user.User type
		userName, err := user.LookupId(strconv.Itoa(int(uid)))
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(3)
		}

		// Of note in the following if blocks, fmt.Sprintf() returns the strting generated from a standard fmt.Printf(), instead of actually printing it, useful for conversion and analysis, particularly for output formatting
		// Set groupname, convert gid to string, use user.LookupGroupId to get a groupname from an id, gives user.Group type
		groupName, err := user.LookupGroupId(strconv.Itoa(int(gid)))
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(3)
		}

		// Now check length of username string against longest_username, if greater, update longest_username
		if len(fmt.Sprintf("%v", userName.Username)) > longest_username {
			longest_username = len(userName.Username)
		}

		// Now check length of groupname string against longest_groupname, if greater, update longest_groupname
		if len(fmt.Sprintf("%v", groupName.Name)) > longest_groupname {
			longest_groupname = len(groupName.Name)
		}

		// Now check length of size string against longest_size, if greater, update longest_size
		if len(fmt.Sprintf("%d", file_info.Size())) > longest_size {
			longest_size = len(fmt.Sprintf("%d", file_info.Size()))
		}

	}

	// Set format strings for header and entries, using longest_name, longest_username, longest_groupname, and longest_size, need to set these so that the %d values will be expanded in these format strings, but leave behind the %s or %v values for the actual printing, %% results in literal % in the output
	format_string_header := fmt.Sprintf("%%-12s |  %%-%ds  |  %%-%ds  | %%-%ds  |  %%-%ds  |  %%-s  \n", longest_name, longest_username, longest_groupname, longest_size)
	format_string_entries := fmt.Sprintf("%%-12v |  %%-%dv  |  %%-%dv  | %%-%dv  |  %%-%dv  |  %%-v  \n", longest_name, longest_username, longest_groupname, longest_size)

	// Begin iterating over files to process them for printing
	for idx, file := range files_in_dir {

		// Get file info, check for errors
		file_info, err := file.Info()
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(3)
		}

		// Direct syscall.Stat_t type conversion from file_info.Sys() method, needed specifically for uid and gid
		file_stat_sys := file_info.Sys().(*syscall.Stat_t)
		uid := file_stat_sys.Uid
		gid := file_stat_sys.Gid

		// Set username, convert uid to string, use user.LookupId to get a username from an id, gives user.User type, same as above but we just care about actual values now
		userName, err := user.LookupId(strconv.Itoa(int(uid)))
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(3)
		}

		groupName, err := user.LookupGroupId(strconv.Itoa(int(gid)))
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(3)
		}

		file_perm := fmt.Sprintf("%v", file_info.Mode())
		if file_info.IsDir() {
			file_perm = "d" + file_perm[1:]
		}

		// Print header if idx is 0, then print the file info
		if idx == 0 {
			fmt.Printf(format_string_header, "Permissions", "Name", "UID", "GID", "Size", "Last Modified")
			// fmt.Printf(format_string_entries, file_perm, stat_fileInfo.Name(), stat.Uid, stat.Gid, stat_fileInfo.Size(), stat_fileInfo.ModTime())
			fmt.Printf("%s\n", sep_string)
		}

		// Print each processed file info
		fmt.Printf(format_string_entries, file_perm, file_info.Name(), userName.Username, groupName.Name, file_info.Size(), file_info.ModTime())
	}
}

func main() {

	// Generate separator string to separate header(s) from rest of file output
	width, _, err := term.GetSize((int(os.Stdout.Fd())))
	sep_string := ""
	if err == nil {
		for i := 0; i < width; i++ {
			sep_string = sep_string + "-"
		}
	}
	// Set empty string for target directory
	targetdir := ""

	// Read cli args, if any, determine location to list (default to cwd)
	args := os.Args
	if arg_len := len(args); arg_len == 1 {
		target_dir, err := os.Getwd()
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(1)
		} else {
			targetdir = target_dir
		}
	} else {
		targetdir = args[1]
	}

	// Get file info
	stat_fileInfo, err := os.Stat(targetdir)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Get raw stat info from the syscall.Stat_t type passed to the .Sys() method
	stat := stat_fileInfo.Sys().(*syscall.Stat_t)

	// Get file mode
	stat_mode := stat_fileInfo.Mode()

	// Check if we are listing a directory or single file, different call stack for each
	if stat_fileInfo.IsDir() {
		files_in_dir, err := os.ReadDir(targetdir)
		if err != nil {
			fmt.Printf("%s\n", err)
			os.Exit(3)
		} else {
			// Will process a slice of os.DirEntry types
			processDirEntries(files_in_dir, sep_string, stat_fileInfo, stat, stat_mode)
		}

		// If we are not listing a directory, we are listing a single file
	} else {

		// Set similar values for formatting output nicely, caring about filename, size, username, groupname
		longest_size := 0
		longest_name := len(stat_fileInfo.Name())
		uid := stat.Uid
		gid := stat.Gid

		// Same uid and gid conversion as above, but we just care about the actual values now
		userName, err := user.LookupId(strconv.Itoa(int(uid)))
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(3)
		}
		groupName, err := user.LookupGroupId(strconv.Itoa(int(gid)))
		if err != nil {
			fmt.Printf("Error: %s\n", err)
			os.Exit(3)
		}

		// Set longest username/groupname, don't need comparisons since the longest ones will always be the single file we are listing
		longest_username := len(userName.Username)
		longest_groupname := len(groupName.Name)

		// update longest_size as well
		if len(fmt.Sprintf("%d", stat_fileInfo.Size())) > longest_size {
			longest_size = len(fmt.Sprintf("%d", stat_fileInfo.Size()))
		}

		// Same string fstring header and entry stuff as above
		format_string_header := fmt.Sprintf("%%-12s |  %%-%ds  |  %%-%ds  | %%-%ds  |  %%-%ds  |  %%-s  \n", longest_name, longest_username, longest_groupname, longest_size)
		format_string_entry := fmt.Sprintf("%%-12v |  %%-%dv  |  %%-%dv  | %%-%dv  |  %%-%dv  |  %%-v  \n", longest_name, longest_username, longest_groupname, longest_size)

		// Same fstring printing as above
		fmt.Printf(format_string_header, "Permissions", "Name", "UID", "GID", "Size", "Last Modified")
		fmt.Printf(format_string_entry, stat_mode.Perm(), stat_fileInfo.Name(), userName.Username, groupName.Name, stat_fileInfo.Size(), stat_fileInfo.ModTime())
	}
}
