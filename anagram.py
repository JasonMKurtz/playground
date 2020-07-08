#!/usr/bin/python
import string

def isAnagram(letters, target):
    lettermap = dict()
    letters = string.replace(letters, " ", "")
    target = string.replace(target, " ", "")
    for letter in letters:
        letter = string.lower(letter)
        m = lettermap.get(letter) or 1
        lettermap[letter] = m + 1

    for letter in target:
        letter = string.lower(letter)
        lm = (lettermap.get(letter) or 0) - 1
        if lm < 0:
            return False

    return True

print(isAnagram('abcdef', 'gh'))
print(isAnagram('October Sky', 'Rocket Boys'))
        

    
    
