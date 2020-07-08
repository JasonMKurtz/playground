#!/usr/bin/python

import string

def isPalindrome(word):
    print("Checking: %s" % (word))
    word = word.replace(" ", "")
    for i in range(len(word) / 2):
        li = -(i + 1)
        print("Checking the pair (%s, %s)" % (word[i], word[li]))
        if string.lower(word[i]) != string.lower(word[li]):
            return False
    return True

print(isPalindrome("ana"))
print(isPalindrome("assa"))
print(isPalindrome("asdf"))
print(isPalindrome("A man a plan a canal panama"))

