#!/usr/bin/python
import operator
str = "we don't need no stinkin' editors"
count = { } 
for c in str: 
    if c != " ": 
        if str.count(c) in count: 
            if c not in count[str.count(c)]: 
                count[str.count(c)].append(c)
        else: 
            count[str.count(c)] = [c]
print str
print count
