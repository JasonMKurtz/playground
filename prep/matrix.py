#!/usr/bin/python
matrix = [ [1, 2, 3, 3], [0, 0, 0, 3], [0, 1, 3, 3], [0, 3, 3, 0] ]
def neighbors(x, y): 
    global total
    if x+1 >= cols or y+1 >= rows: 
        return total+1
    if matrix[x][y] == matrix[x+1][y]: 
        total += 1
        return neighbors(x+1, y)
    if matrix[x][y] == matrix[x][y+1]:
        total += 1
        return neighbors(x, y+1)
    if matrix[x][y] == matrix[x+1][y+1]: 
        total += 1
        return neighbors(x+1, y+1)
    if matrix[x][y] == matrix[x-1][y]: 
        total += 1
        return neighbors(x-1, y)
    if matrix[x][y] == matrix[x-1][y+1]: 
        total += 1
        return neighbors(x-1, y+1)
    if matrix[x][y] == matrix[x-1][y-1]: 
        total += 1
        return neighbors(x-1, y-1)
    


rows = 4
cols = 4
total = 0
answer = 0
winners = dict()

for row in xrange(rows): 
    for col in xrange(cols): 
        total = 0
        cur = neighbors(row, col)
        if cur > answer: 
            answer = cur
        
print answer
