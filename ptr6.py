#!/usr/bin/python
import sys
def ptr6(str): 
    return "%s.in-addr.arpa." % (".".join([ ".".join(list(oct.zfill(4)[::-1])) for oct in str.replace("::", "".join([ ":" for i in range(9 - len([ oct for oct in str.split(":") if oct != '' ])) ])).split(":") ][::-1]))

print ptr6("".join(sys.argv[1::]))
