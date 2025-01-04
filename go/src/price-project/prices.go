/*
package main

import (

	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

)

//type priceTax struct {
//	Price int
//	Tax   float32
//}

	type userEntry struct {
		userName string
		ID       int
		err      error
	}

	func processUserEntry() {
		return
	}

//func extractUser(userinfo string) {
//
//}

	func getUser(entry string) (userEntry, error) {
		info := strings.Split(entry, ":")
		if len(info) < 2 {
			return userEntry{}, fmt.Errorf("Error, user:ID not in correct format:\n %s", entry)
		}
		username := strings.TrimSpace(info[0])
		sid := strings.TrimSpace(info[1])
		id, err := strconv.Atoi(sid)
		if err != nil {
			return userEntry{}, fmt.Errorf("Error, non-numeric user ID given, ID was: %s\n", info[1])
		}

		user := userEntry{userName: username, ID: id}
		return user, nil
	}

	func decodeUsers(reader io.Reader) chan userEntry {
		ch := make(chan userEntry, 1)
		go func() {
			scanner := bufio.NewScanner(reader)
			for scanner.Scan() {
				uentry, err := getUser(scanner.Text())
				if err != nil {
					uentry.err = err
					ch <- uentry
					return
				}
				ch <- uentry
			}
		}()
		return ch

}

	func main() {
		//fmt.Printf("Placeholder\n")
		args := os.Args[1:]
		filepath := args[0]
		//fmt.Printf("Filepath: %s\n", filepath)
		//price_matrix := []priceTax{}
		//userEntries := []userEntry{}

		reader, err := os.Open(filepath)
		if err != nil {
			return
		}
		defer reader.Close()
		for user := range decodeUsers(reader) {
			if user.err != nil {
				fmt.Println("Error: ", user.err)
			} else {
				fmt.Printf("User: \n+%v", user)
			}
		}

		//	for _, entry := range userEntries {
		//		fmt.Printf("Username: %s\n", entry.userName)
		//		fmt.Printf("User ID: %d\n", entry.ID)
	}
*/
package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

// filePath is the path to file we are going to read and write to.
//const filePath = "users.txt"

// init is a neat feature in Go. Every package can have init() functions and it is
// the only type of function whose name can be duplicated. init() is executed when
// the package loads, usually to initialize some variables or do other setup. In this
// case we use it to write some content to a file representing our users.
//func init() {
//	content := []byte("jdoak:0\nsmurphy:1\ndjustice:2")
//
//	if err := os.WriteFile(filePath, content, 0644); err != nil {
//		panic(err)
//	}
//}

// User represents our user data.
type User struct {
	// Name is our user's username.
	Name string
	// ID is their unique numeric ID in the system.
	ID int
	// err indicates there was an error in stream reading.
	err error
}

func (u User) String() string {
	return fmt.Sprintf("User: %s\nID: %d\n", u.Name, u.ID)
}

// String implememnts fmt.Stringer. It will output the data as "user:id", such as "jdoak:0".
//}

// getUser takes a string that should be formatted as [user]:[id], such as "jdoak:0" and returns
// a User object.
func getUser(s string) (User, error) {
	sp := strings.Split(s, ":")
	if len(sp) != 2 {
		return User{}, fmt.Errorf("record: ('%s'); Record was not <user>:<ID> format", s)
	}
	id, err := strconv.Atoi(sp[1])
	if err != nil {
		return User{}, fmt.Errorf("record: ('%s'); Record had non-numeric ID", s)
	}
	return User{Name: strings.TrimSpace(sp[0]), ID: id}, nil
}

// decodeUsers reads from a io.Reader breaking the file entries by a carriage return(\n)
// and decodes the entries to User objects and returns them on a channel.
// If there was an error, the returned entry will have .err != nil.
func decodeUsers(ctx context.Context, r io.Reader) chan User {
	ch := make(chan User, 1)

	// Spin of goroutine off to feed the channel we will return.
	go func() {
		// Close our channel on exit, signaling we are done.
		defer close(ch)

		// Wrap a Scanner around our reader so we can read each line of content.
		scanner := bufio.NewScanner(r)
		for scanner.Scan() { // Scan until nothing to scan.
			if ctx.Err() != nil { // Context was cancelled, return error.
				ch <- User{err: ctx.Err()}
			}
			// Turn the line of text into a User object.
			//fmt.Println("about to scan user")
			u, err := getUser(scanner.Text())
			//fmt.Println("User scanned")
			if err != nil { // line was in incorrect format, return an error.
				u.err = err
				ch <- u
				continue
				//return
			}
			// Everything was fine, return a user record.
			ch <- u
		}
		//close(ch)
	}()
	// Returns the channel we will read off of.
	return ch
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Open the file we created with init().
	filePath := os.Args[1]
	f, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// Start decoding the file one line at a time.
	//fmt.Println("decoding users")
	ch := decodeUsers(ctx, f)
	//fmt.Println("Channel returned")

	// Read each line of output and write the record to the screen.
	for u := range ch {
		if u.err != nil {
			fmt.Printf("Error processing line with %s\n\n", u.err)
			continue
		}
		fmt.Println(u)
	}
}
