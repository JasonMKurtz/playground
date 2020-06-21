#!/usr/bin/python
matrix = [ [1, 2, 3, 3], [0, 0, 0, 3], [0, 1, 3, 3], [0, 3, 3, 0] ]
positions = list()
rows = 4
cols = 4
last = 0
def neighbors(x, y): 
    global total
    # (x, y), right (x+1,y), down (x, y+1), right and down (x+1,y+1), left (x-1, y), up (x, y-1), left and up (x-1, y-1), right and up (x+1,y-1), left and down (x-1,y+1) 
    #if x+1 >= cols or y+1 >= rows or [x,y] in positions: 
    if x >= cols or y >= rows or [x,y] in positions or x < 0 or y < 0: 
        return 0
    last = matrix[x][y]        
    if matrix[x][y] == last: 
        positions.append(matrix[x][y])
    return 1 + neighbors(x, y) + neighbors(x+1,y) + neighbors(x,y+1) + neighbors(x-1,y) + neighbors(x,y-1) + neighbors(x+1,y+1) + neighbors(x-1,y-1) + neighbors(x+1,y-1) + neighbors(x-1,y+1)
    
total = 0
answer = 0
winners = dict() # keep track of which items we have checked so far

for row in xrange(rows): 
    for col in xrange(cols): 
        if [row,col] not in positions: 
            total = 0
            neigh = neighbors(row, col)
            print "%d (%d, %d) %d -> %s" % (matrix[row][col], row, col, neigh, positions)
            positions = list()
