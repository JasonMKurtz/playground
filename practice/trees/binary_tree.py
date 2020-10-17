#!/usr/local/bin/python3 

from typing import Dict, List, Set
from queue import Queue

class Graph:
    pairs: Dict[str, List[str]] = {}
    def __init__(self):
        self.pairs = {
            0: set([1, 2]),
            1: set([0, 3]),
            2: set([0, 3]),
            3: set([1, 4]),
        }
    
    def is_binary_tree(self, item, seen={}) -> bool:
        # Does our node exist?
        # Does our node have more than two neighbors?
        # Do any of its neighbors have more than two neighbors?

        node = self.pairs.get(item)
        if node is None or len(node) > 2:
            return False

        visited = Queue()
        visited.put(item)
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
