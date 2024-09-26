#!/bin/python3

import os, sys, re, string
from collections import Counter 


def countWords(filepath):
    wordCounter = Counter() 
    with open(filepath, 'r') as fp:
        words = fp.read().lower()
        puncTranslator = str.maketrans('','',string.punctuation)
        cleanWordString = words.translate(puncTranslator)
        wordList = cleanWordString.split(' ')
        wordCounts = Counter(wordList)
        return wordCounts


def main(): 
    if len(sys.argv) < 2:
        filepath = './words.txt'
    else:
        filepath = sys.argv[1]

    try:
        fileExists = os.stat(filepath)
    except Exception as e:
        print(f'Issue getting file stat: {e}')

   

if __name__ == '__main__': 
    main()