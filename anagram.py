#!/usr/bin/python

import string

def isAnagram(word):
    print("Checking: %s" % (word))
    word = word.replace(" ", "")
    for i in range(len(word) / 2):
        li = -(i + 1)
        print("Checking the pair (%s, %s)" % (word[i], word[li]))
        if string.lower(word[i]) != string.lower(word[li]):
            return False
    return True

print(isAnagram("ana"))
print(isAnagram("assa"))
print(isAnagram("asdf"))
print(isAnagram("A man a plan a canal panama"))

