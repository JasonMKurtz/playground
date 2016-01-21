#!/usr/bin/python
def tail(f, lines): 
    blocksize = 10
    lines_found, bytes_read = 0, 0
    with open(f, 'r') as fh: 
        fh.seek(0, 2)
        size = fh.tell() # how big is the file?
        step = min(blocksize, size)
        while (lines_found <= lines and bytes_read < size): 
            fh.seek(-step, 1) # from end, go back x bytes
            bytes_read += step
            lines_found += fh.read().count("\n")
        
        fh.seek(-bytes_read, 2)
        print '\n'.join(fh.read().split("\n"))
tail("dafile", 3)
