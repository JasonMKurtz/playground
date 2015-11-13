#!/usr/bin/python
from subprocess import Popen, PIPE
ls   = Popen("ls | xargs stat --format '%n %i'", stdout=PIPE, stderr=PIPE, shell=True)
lines = ls.communicate()[0].split("\n")
files = dict()
for line in lines: 
    if len(line.split()) > 1: 
        if line.split()[1] in files: 
            files[line.split()[1]] += [ line.split()[0] ] 
        else:
            files[line.split()[1]] = [ line.split()[0] ]

for inode, files in files.items(): 
    if len(files) > 1: 
        print "%s are hard-linked files (%s)" % (', '.join(files), inode)
