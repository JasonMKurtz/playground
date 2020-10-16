#!/usr/local/bin/python3 


class Node:
    def __init__(self, val, neighbors = []):
        self.neighbors = neighbors
        self.val = val
    
    def addNeighbors(self, neighbors):
        self.neighbors = neighbors
    
    def __repr__(self): 
        neighbors = [n.val for n in self.neighbors]
        if not len(neighbors):
            return f"{self.val}"
        
        return f"{self.val} {neighbors}"
    
def make(): 
    return Node(0, [
        Node(1, [Node(2, [Node(3), Node(4), Node(8)])]),
        Node(10, [Node(15, [Node(16), Node(17), Node(9)])]),
    ])

from queue import Queue
from typing import Dict

"""
NameError: name 'bfs_shortest_path' is not defined
>>> from path_exists import * 
>>> bfs_shortest_path(make(), 15)
[0, 10, 15]
"""

def bfs_shortest_path(node, target):
    # for each node in our tree, visit its children, and their children, and so on
    if node.val == target:
        return [node.val, target]

    visited = Queue()
    visited.put((node, [node.val]))
    seen: Dict[Node, bool] = {node: True}
    path = [node.val]
    while not visited.empty():
        n, path = visited.get()
        if n.val == target:
            return path + [n.val]

        for node in n.neighbors:
            if node is None:
                continue
        
            if node.val == target:
                return path + [node.val]
            
            if not seen.get(node):
                visited.put((node, path + [node.val]))
            
    
