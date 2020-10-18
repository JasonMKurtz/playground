#!/usr/local/bin/python3 

from typing import Dict, List, Set, Optional
from queue import Queue

class Graph:
    pairs: Dict[str, List[str]] = {}
    def __init__(self):
        self.pairs = {
            10: [5, 12],
            5: [2, 7],
            12: [9, 17],
        }
    
    def has_cycle(self, item, visited={}) -> bool:
        if visited.get(item):
            return True
        
        node = self.pairs.get(item)
        if not node:
            return False
        
        visited.update({item: True})
        for item in node:
            if self.has_cycle(item, visited):
                return True
        return False
    
    def preorder(self, item):
        node = self.pairs.get(item)
        print(item)
        if node:
            self.preorder(node[0])
            self.preorder(node[1])
    
    def inorder(self, item):
        node = self.pairs.get(item)

        if node:
            self.inorder(node[0])
        print(item)
        if node:
            self.inorder(node[1])
    
    def postorder(self, item):
        node = self.pairs.get(item)
        if node:
            self.postorder(node[0])
            self.postorder(node[1])
        print(item)
    
    def is_bst(self, item, minval = None, maxval = None) -> bool:
        if minval is not None and item < minval:
            return False
        
        if maxval is not None and item > maxval:
            return False

        node = self.pairs.get(item)
        if node is None:
            return True
        
        left = self.is_bst(node[0], minval, item)
        right = self.is_bst(node[1], item, maxval)

        return left and right

    def is_binary_tree(self, start) -> bool:
        # Does our node exist?
        # Does our node have more than two neighbors?
        # Do any of its neighbors have more than two neighbors?

        node = self.pairs.get(start)
        if node is None or len(node) > 2:
            return False

        visited = Queue()
        visited.put(start)
        seen = {}
        while not visited.empty():
            i = visited.get()
            for neighbor in self.pairs.get(i):
                if seen.get(neighbor):
                    continue
                node = self.pairs.get(neighbor)
                if node is None:
                    continue
                if len(node) > 2:
                    return False
                seen.update({neighbor: True})
                visited.put(neighbor)
        return True
