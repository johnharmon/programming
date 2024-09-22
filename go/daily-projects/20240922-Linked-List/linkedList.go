package main

import "fmt"

type listNode struct {
	data int
	next *listNode
}

type linkedList struct {
	head *listNode
}

func findListEnd(list *linkedList) (previousNode *listNode, lastNode *listNode) {
	previousNode = nil
	//	fmt.Printf("finding the end of the list\n")
	//	fmt.Printf("List is %+v\n", list)
	head := list.head
	//fmt.Printf("Head is: %+v\n", head)
	//fmt.Printf("head found\n")
	for {
		//fmt.Printf("Loop entered\n")
		if head.next == nil {
			lastNode = head
			return previousNode, lastNode
		} else {
			previousNode = head
			head = head.next
		}
	}
}

func findValue(list *linkedList, value int) *listNode {
	head := list.head
	for {
		if head.data == value {
			return head
		} else {
			if head.next != nil {
				head = head.next
			} else {
				return nil
			}
		}
	}
}

func removeFromEnd(list *linkedList) {
	newTailNode, _ := findListEnd(list)
	if newTailNode == nil { // Single node list
		list.head = nil
	} else {
		newTailNode.next = nil
	}
	/*
		var previousNode listNode
		index := 0
		for {
			if head.next {
				previousNode = head
				head = head.next
				index++
			} else {
				if index == 0 {
					head = nil
				} else {
					head = nil
					previousNode.next = nil
				}
				break
			}
		}
	*/
}

func removeNodePtr(head **listNode, previousNode *listNode) {
	if previousNode != nil { // We are not at the head of the list
		if (*head).next != nil { // we are not at the end either
			previousNode.next = (*head).next
			*head = nil
		} else { // We are at the end of the list
			*head = nil // Delt last node
		}
	} else { // We are at the current head of the list
		if (*head).next != nil {
			*head = (*head).next // Move head to the second element, GC will clean up the old head
		} else { // We are a one node list
			*head = nil // Zero out only existing node
		}
	}
}

func removeNode(head *listNode, previousNode *listNode) {
	if previousNode != nil { // We are not at the head of the list
		if head.next != nil { // we are not at the end either
			previousNode.next = head.next
			head = nil
		} else { // We are at the end of the list
			fmt.Println("removing last node")
			previousNode.next = nil
		}
	}
}

func removeNodeValue(list *linkedList, value int) bool {
	var previousNode *listNode
	head := list.head
	if list.head.data == value {
		if list.head.next != nil {
			list.head = list.head.next
			return true
		} else {
			list.head = nil
			return true
		}
	} else {
		//index := 0
		for {
			if head.data == value {
				removeNode(head, previousNode)
				return true
			}
			if head.next != nil {
				previousNode = head
				head = head.next
			} else {
				return false
			}
		}
	}
}

func insertAtBeginning(list *linkedList, value int) {
	newHead := &listNode{
		next: list.head,
		data: value,
	}
	list.head = newHead
}

func insertAtEnd(list *linkedList, value int) {
	//fmt.Printf("Adding %d to list\n", value)
	_, listEnd := findListEnd(list)
	newEnd := &listNode{
		data: value,
		next: nil,
	}
	listEnd.next = newEnd
}

func printList(list *linkedList) {
	head := list.head
	index := 0
	if head != nil {
		for {
			if head != nil {
				fmt.Printf("Linked List as position %d has the value of %+v\n", index, head)
				index++
				head = head.next
			} else {
				return
			}
		}
	}
}

func main() {

	// make first node
	firstNode := &listNode{
		data: 5,
		next: nil,
	}

	// make list struct
	myList := linkedList{
		head: firstNode,
	}

	// populate list
	for i := 0; i < 10; i++ {
		//fmt.Printf("Adding to linked list\n")
		insertAtEnd(&myList, i)
	}

	// display list after population
	//printList(&myList)

	// add to beginning of list
	insertAtBeginning(&myList, 100)

	// display list after insertion at the beginning
	//printList(&myList)

	// remove value from list
	removeNodeValue(&myList, 9)
	//fmt.Printf("Value was removed: %t\n", removed)

	// Print list after value removal
	printList(&myList)

}
