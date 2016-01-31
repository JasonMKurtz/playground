#!/usr/bin/python
n = int(raw_input().strip())
a = 0
a_primary = []
a_secondary = []
for oPos, a_i in enumerate(xrange(n)):
   line = raw_input().strip().split(' ') 
   for iPos, i in enumerate(line): 
        if iPos == oPos: 
            a_primary.append(int(i)) 
        if iPos+oPos == n-1: 
            a_secondary.append(int(i))
a = abs(sum(a_primary) - sum(a_secondary))
print a
