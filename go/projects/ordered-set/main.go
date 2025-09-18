package main

// inserting a comment for github testing
import (
	"fmt"
	"math/rand"
	"os"

	"harmonlab.io/set"
)

func main() {
	mySet := set.NewOrderedSet[int]()
	for i := 0; i < 5; i++ {
		res := mySet.Append(i)
		fmt.Printf("Change Applied: %t\n", res)
	}
	for i := 0; i < 5; i++ {
		res := mySet.Append(i)
		fmt.Printf("Change Applied: %t\n", res)
	}
	//	fmt.Println("Deleting items:")
	//	for i := 0; i < 5; i++ {
	//		res := mySet.Delete(i)
	//		fmt.Printf("Change Applied: %t\n", res)
	//	}
	//	fmt.Println("Deleting items:")
	//	for i := 0; i < 5; i++ {
	//		res := mySet.Delete(i)
	//		fmt.Printf("Change Applied: %t\n", res)
	//	}
	for i := 0; i < 5; i++ {
		res := mySet.Append(rand.Int())
		fmt.Printf("Change Applied: %t\n", res)
	}
	mySet.DumpBitMap(os.Stdout)
	mySet.DumpItems(os.Stdout)
	mySet.DumpSequenceMap(os.Stdout)
	mySet.DumpDataMap(os.Stdout)
	for i := 0; i < 5; i++ {
		setLen := mySet.Len()
		idx := rand.Int() % setLen
		if res, ok := mySet.GetIdx(idx); ok {
			fmt.Printf("Value: %d, index: %d\n", res.Value, idx)
		} else {
			fmt.Printf("Error: index(%d) not found in set\n", idx)
		}

	}
}
