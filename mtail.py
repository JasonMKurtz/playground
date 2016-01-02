#!/usr/bin/python
def tail(lines2find): 
    file = "abc"
    block = 1024
    with open(file, 'r') as f: 
        f.seek(0, 2)
        total_bytes = f.tell()
        scanned, lines = 0, 0
        while (total_bytes > scanned and lines2find + 1 > lines): 
            byte_block = min(block, total_bytes - scanned)
            f.seek(-byte_block, 2)
            scanned += byte_block
            lines += f.read().count("\n")
        f.seek(-scanned, 2)
        return '\n'.join([ line.strip() for line in f.readlines()[-lines2find:] ])
print tail(10)        
