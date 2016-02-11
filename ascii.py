#!/usr/bin/python 
""" Some nice ascii art. """
print "\n".join([ "{a}{b}{a}".format(a="*"*(30-i), b="||"*(i)) for i in range(1, 30) ]), "\n", "\n".join([ "{a}{b}{a}".format(a="*"*i, b="||"*(30-i)) for i in range(2, 30) ])
