#!/usr/bin/python
class File: 
    def __init__(self, filename, blocksize=1024): 
        self.file = filename
        self.block = blocksize
        self.reads = 0

    def numreads(self): 
        return self.reads

    def tail(self, lines2find): 
        with open(self.file, 'r') as f: 
            f.seek(0, 2)
            total_bytes = f.tell()
            scanned, lines = 0, 0
            while (total_bytes > scanned and lines2find + 2 >= lines): 
                byte_block = min(self.block, total_bytes - scanned)
                f.seek(-byte_block, 2)
                scanned += byte_block
                lines += f.read().count("\n")
                self.reads += 1
            f.seek(-scanned, 2)
            return '\n'.join([ line.strip() for line in f.readlines()[-lines2find:] ])

    def unique_tail(self, lines): 
        lines = self.tail(lines).split("\n")
        b = [ ] 
        for i in lines: 
            if i not in b: 
                b.append(i)

        return '\n'.join(b)


f = File("abc", 10)
#print f.tail(10)        
#print '\n\n'
print f.unique_tail(10)
