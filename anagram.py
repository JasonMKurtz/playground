#!/usr/bin/python

import string

def isAnagram(word):
    print("Checking: %s" % (word))
    word = word.replace(" ", "")
    l = len(word)
    pairs = l / 2
    for i in range(pairs):
        li = -(i + 1)
        print("Checking the pair (%s, %s)" % (word[i], word[li]))
        if string.lower(word[i]) != string.lower(word[li]):
            return False
    return True

print(isAnagram("ana"))
print(isAnagram("assa"))
print(isAnagram("asdf"))
print(isAnagram("A man a plan a canal panama"))

