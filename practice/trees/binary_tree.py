#!/usr/local/bin/python3 

from typing import Dict, List, Set, Optional
from queue import Queue

class Graph:
    pairs: Dict[str, List[str]] = {}
    def __init__(self):
        self.pairs = {
            10: [5, 12],
            5: [2, 7],
            12: [11, 17],
        }
    
    def inorder(self, item):
        node = self.pairs.get(item)

        if node:
            self.inorder(node[0])
        print(item)
        if node:
            self.inorder(node[1])
    
    def is_bst(self, item, l = None, r = None) -> bool:
        node = self.pairs.get(item)
        if node is None:
            return True
        
        if l is not None and item < l:
            return False
        
        if r is not None and item > r:
            return False
        
        return self.is_bst(node[0], l, item) and \
            self.is_bst(node[1], item, r)

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
