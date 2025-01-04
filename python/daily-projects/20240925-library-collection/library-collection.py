#!/bin/python3

from collections import namedtuple, defaultdict, deque, Counter


Book = namedtuple('Book', ['title', 'author', 'year_published', 'borrowed' ])

def add_book(books, book):
    books.append(book)
    return


def borrow_book(book, borrow_counter, recent_que):
    borrow_counter.update((book.title, book.author)) 
    recent_que.append(book)
    return

def return_book(book, recent_que):
    if book in recent_que:
        recent_que.remove(book)
    return

def most_common(book_counter):
     return book_counter.most_common(3)

def main():
    books = defaultdict(lambda: Book('N/A', 'N/A', 'N/A', False))
    borrow_counter = Counter()
    recently_borrowed = deque()
    recently_borrowed.remove
    book = Book('Atomic Habits', 'James Clear', '2018')

if __name__ == '__main__':
	main()
