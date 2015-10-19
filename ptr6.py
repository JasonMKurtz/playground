#!/usr/bin/python
"""
Usage: ./ptr6.py <ipv6 address>
Example: 

$ ./ptr6.py 2001:470:b995:a019::1
1.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.9.1.0.a.5.9.9.b.0.7.4.0.1.0.0.2.ip6.arpa.

"""
import sys
def ptr6(address): return "%s.ip6.arpa." % ('.'.join([ '.'.join(oct.zfill(4)[::-1]) for oct in address.replace("::", ''.join([ ":" for i in range(9 - address.count(":")) ])).split(":") ][::-1]))

print ptr6("".join(sys.argv[1::]))
