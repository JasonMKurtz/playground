#!/usr/bin/python 
""" Some nice ascii art. """
import sys
size = int(sys.argv[1]) if len(sys.argv) > 1 else 15
print "\n".join([ "{a}{b}{a}".format(a="*"*(size-i), b="||"*(i)) for i in range(1, size) ]), "\n", "\n".join([ "{a}{b}{a}".format(a="*"*i, b="||"*(size-i)) for i in range(2, size) ])
