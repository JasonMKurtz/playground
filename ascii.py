#!/usr/bin/python 
""" Some nice ascii art. """
size = 15
print "\n".join([ "{a}{b}{a}".format(a="*"*(size-i), b="||"*(i)) for i in range(1, size) ]), "\n", "\n".join([ "{a}{b}{a}".format(a="*"*i, b="||"*(size-i)) for i in range(2, size) ])
