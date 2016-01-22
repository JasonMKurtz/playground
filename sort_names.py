#!/usr/bin/python
"""
3
Mike Thomson 20 M
Robert Bustle 32 M
Andria Bustle 30 F
"""
import operator
def format(outer): 
   def inner(*args, **kwargs): 
        names = outer(*args, **kwargs).split("\n")
        ret = []
        for name in names: 
            title = "Mr." if name.split()[3].lower() == "m" else "Ms." 
            ret.append("%s %s" % (title, ' '.join(name.split()[0:2])))

        return '\n'.join(ret)

   return inner

@format
def dosort(str): 
    s = sorted([ i.split() for i in str.split("\n") ], key=operator.itemgetter(2))
    return '\n'.join([ ' '.join(n) for n in s ])

print dosort("Mike Thomson 20 M\nRobert Bustle 32 M\nAndria Bustle 30 F")
