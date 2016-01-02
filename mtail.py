#!/usr/bin/python
class File: 
    def __init__(self, filename, blocksize=1024): 
        self.file = filename
        self.block = blocksize

    def tail(self, lines2find): 
        with open(self.file, 'r') as f: 
            f.seek(0, 2)
            total_bytes = f.tell()
            scanned, lines = 0, 0
            while (total_bytes > scanned and lines2find + 1 > lines): 
                byte_block = min(self.block, total_bytes - scanned)
                f.seek(-byte_block, 2)
                scanned += byte_block
                lines += f.read().count("\n")
            f.seek(-scanned, 2)
            return '\n'.join([ line.strip() for line in f.readlines()[-lines2find:] ])

print File("abc").tail(10)        
