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
        Node(4, [Node(5, [Node(6), Node(7), Node(9)])]),
    ])

from queue import Queue
from typing import Dict

def bfs(node):
    # for each node in our tree, visit its children, and their children, and so on
    visited = Queue()
    visited.put(node)
    seen: Dict[Node, bool] = {node: True}
    levels: List[int] = []
    while not visited.empty():
        n = visited.get()
        print(n)
        for node in n.neighbors:
            if node is None:
                continue
            if not seen.get(node):
                visited.put(node)