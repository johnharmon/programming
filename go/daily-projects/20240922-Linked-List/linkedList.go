package main 


type listNode struct {
	data int 
	next *listNode 
}

type linkedList struct {
	head *listNode
}


func findListEnd(list *linkedList) (previousNode *listNode, lastNode *listNode) {
	head := list.head
	var previousNode *listNode
	for {
		if !head.next {
			return previousNode, lastNode := head
		} else {
			previousNode = head
			head = head.next
		}
	}
}


func removeFromEnd(list *linkedList){
	newTailNode, oldTailNode := findListEnd(list)
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

func removeNode(head *listnode, previousNode *listNode) {
	if previousNode != nil { // We are not at the head of the list
		if head.next { // we are not at the end either
			previousNode.next = head.next // Update previous node pointer
			head = nil
		} else { // We are at the end of the list
			head = nil // Delt last node
			previousNode.next = nil // Zero out pointer for the next node on the new tail
		}
	} else { // We are at the current head of the list
		if head.next {
			head = head.next // Move head to the second element, GC will clean up the old head
		} else { // We are a one node list
			head = nil // Zero out only existing node
		}
	}
}

func removeNodeValue(list *linkedList, value int) bool {
	var previousNode *listNode 
	head := list.head
	index := 0
	for {
		if head.data == value {
			removeNode(head, previousNode)
			return true
		}
		if head.next {
			head = head.next 
		} else {
			return false
		}
	}
}

func insertAtBeginning (list *linkedList, value int) {
	newHead := &listNode{ 
		next: listNode.next 
		data: int  
	}
	list.next = newHead
}

func insertAtEnd(list *linkedList, value int) {
	_, listEnd := findListEnd(list)
	newEnd := &listNode{
		data: value 
		next: nil
	}
	listEnd.next = newEnd
}

func main () {

}