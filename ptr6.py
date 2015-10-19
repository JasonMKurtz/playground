#!/usr/bin/python
"""
Usage: ./ptr6.py <ipv6 address>
Example: 

$ ./ptr6.py 2001:470:b995:a019::1
1.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.9.1.0.a.5.9.9.b.0.7.4.0.1.0.0.2.ip6.arpa.

"""
import sys
<<<<<<< HEAD
def ptr6(address): return "%s.ip6.arpa." % ('.'.join([ '.'.join(oct.zfill(4)[::-1]) for oct in address.replace("::", ''.join([ ":" for i in range(9 - address.count(":")) ])).split(":") ][::-1]))
=======
def ptr6(str): return "%s.ip6.arpa." % (".".join([ ".".join(list(oct.zfill(4)[::-1])) for oct in str.replace("::", "".join([ ":" for i in range(9 - len([ oct for oct in str.split(":") if oct != '' ])) ])).split(":") ][::-1]))
>>>>>>> a5ae27dc7f53a75d62fa8276cb4be681f32646c6

print ptr6("".join(sys.argv[1::]))
