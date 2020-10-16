#!/usr/local/bin/python3 

from typing import Dict, List, Set
from queue import Queue

class Graph:
    pairs: Dict[str, List[str]] = {}
    def __init__(self):
        self.pairs = {
            'a': set(['b', 'c']),
            'b': set(['a', 'd']),
            'c': set(['a', 'd']),
            'd': set(['b', 'c', 'e']),
            'e': set(['f', 'g']),
            'f': set(['e']),
            'g': set(['e'])
        }
    
    def shortest_path(self, src, dest) -> List[str]:
        if not self.pairs.get(src):
            raise Exception(f"No node {src} exists.")
        elif not self.pairs.get(dest):
            raise Exception(f"No node {dest} exists.")

        if src == dest:
            return [src]
        
        seen: Set[str] = set()
        q = Queue() 
        q.put((src, [src]))
        while not q.empty():
            item, path = q.get()

            if item == dest:
                return path
            
            pairs = self.pairs.get(item)
            if not pairs:
                raise Exception(f"No node {item} exists.")

            for node in self.pairs.get(item):
                if node not in seen:
                    q.put((node, path + [node]))
                    seen.add(node)
        return []