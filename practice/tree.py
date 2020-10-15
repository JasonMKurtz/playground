#!/usr/local/bin/python3

class Node:
    def __init__(self, val):
        self.left = None
        self.right = None
        self.val = val
    
    def __repr__(self): 
        return f"({self.val})"


class Tree:
    def __init__(self, root, limit):
        self.root = Node(root)
        self.limit = limit
    
    def __repr__(self):
        return f"[{self.root.left} <-- {self.root} --> {self.root.right}]"
    
def make(): 
    r = Node(0)
    r.left = Node(1)
    r.right = Node(2)
    r.left.left = Node(3)
    r.left.right = Node(4)
    r.right.left = Node(5)
    r.right.right = Node(6)
    return r

"""
make() gives us the following tree:
                0
        1               2
    3       4       5       6
"""

def inOrder(node):
    if node:
        inOrder(node.left)
        print(node.val)
        inOrder(node.right)

def preOrder(node):
    if node:
        print(node.val)
        preOrder(node.left)
        preOrder(node.right)

def postOrder(node):
    if node:
        postOrder(node.left)
        postOrder(node.right)
        print(node.val)

"""
>>> from tree import * 
>>> t = make()
>>> inOrder(t)
3
1
4
0
5
2
6
>>> preOrder(t)
0
1
3
4
2
5
6
>>> postOrder(t) 
3
4
1
5
6
2
0
"""

from queue import Queue
from typing import Dict

def bfs(node):
    # for each node in our tree, visit its children, and their children, and so on
    visited = Queue()
    visited.put(node)
    seen: Dict[Node, bool] = {node: True}
    while not visited.empty():
        n = visited.get()
        print(n)
        for node in [n.left, n.right]:
            if node is None:
                continue
            if not seen.get(node):
                visited.put(node)


"""
>>> bfs(t)
(0)
(1)
(2)
(3)
(4)
(5)
(6)
"""