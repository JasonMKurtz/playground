#!/usr/bin/python
import sys
def msplice(num): return ''.join([ d.replace(d, "%s%s" % (d, "," if pos != 0 and (pos % 3) == 0 else "")) for pos, d in enumerate(num[::-1]) ][::-1])
print msplice(''.join(sys.argv[1::]))
