#!/usr/bin/python
import sys
""" Input needs to be changed to sys.argv[1::] if we are reading arguments. """ 
input = sys.stdin.read() 
def msplice(num): return ''.join([ d.replace(d, "%s%s" % (d, "," if pos != 0 and (pos % 3) == 0 else "")) for pos, d in enumerate(num[::-1]) ][::-1])
print msplice(input.strip())
