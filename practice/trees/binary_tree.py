#!/usr/local/bin/python3 

from typing import Dict, List, Set, Optional
from queue import Queue

class Graph:
    pairs: Dict[str, List[str]] = {}
    def __init__(self):
        self.pairs = {
            10: [5, 12],
            5: [2, 11],
            12: [11, 17],
        }
    
    def is_bst(self, start) -> bool:
        # For every subtree, [0] < root < [1]
        node = self.pairs.get(start)
        if node is None:
            return True

        if len(node) > 2:
            return False
    
        left = node[0]
        right = node[1]

        if left is not None and left > start:
            return False

        if right is not None and right < start:
            return False

        if not self.is_bst(left) or not self.is_bst(right):
            return False

        return True
    
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
