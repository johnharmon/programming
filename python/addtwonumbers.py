#!/bin/python3
from typing import List


class ListNode:
     def __init__(self, val=0, next=None):
         self.val = val
         self.next = next
class Solution:
    def addTwoNumbers(self, l1: Optional[ListNode], l2: Optional[ListNode]) -> Optional[ListNode]:
        string1=''
        string2=''
        node1 = l1
        node2 = l2
        while node1.next:
            string1+=str(node1.val)
            if node1.next:
                node1=node1.next
            else:
                break
        while node2.next:
            string1+=str(node2.val)
            if node2.next:
                node2=node2.next
            else:
                break
        string1 = string1[::-1]
        string2 = string2[::-1]
        sum = int(string1) + int(string2)
        sum = string(sum)
        result_list = []
        for digit in range(0,len(sum)):
            my_node = ListNode(int(sum[digit], next = None))
            result_list.append(my_node)
            if digit < len(sum) and digit > 0:
                result_list[digit-1].next = my_node
        result_list.reverse()
        return result_list




list1 = []
list2 = []
for number in range(0,9):
    list1.append(ListNode(val=number))
    list2.append(ListNode(val=number))
    # Definition for singly-linked list.
# class ListNode:
#     def __init__(self, val=0, next=None):
#         self.val = val
#         self.next = next
class Solution:
    def addTwoNumbers(self, l1: Optional[ListNode], l2: Optional[ListNode]) -> Optional[ListNode]:
        string1=''
        string2=''
        node1 = l1
        node2 = l2
        while node1.next:
            string1+=str(node1.val)
            if node1.next:
                node1=node1.next
            else:
                break
        while node2.next:
            string2+=str(node2.val)
            if node2.next:
                node2=node2.next
            else:
                break
        string1 = string1[::-1]
        string2 = string2[::-1]
        sum = int(string1) + int(string2)
        sum = str(sum)
        result_list = []
        for digit in range(-1, -len(sum), -1):
            my_node = ListNode(int(sum[digit]), next = None)
            result_list.append(my_node)
            if digit > -len(sum) and digit < -1:
                result_list[digit+1].next = my_node
            
        return result_list