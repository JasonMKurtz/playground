#!/usr/local/bin/python3

from typing import Dict, List

class Table:
    items: Dict[int, List[str]] = {}
    size: int

    def __init__(self, size: int):
        self.size = size
        pass

    def hash(self, value: str) -> str:
        return str(len(value) % self.size) 
    
    def add(self, value: str) -> None:
        if len(self.items) > self.size:
            return

        index: int = self.hash(value)
        if index in self.items:
            update: List[str] = self.items[index] + [value]
        else:
            update: List[str] = [value]
        
        self.items.update({index: update})
        print(self)

    def __repr__(self) -> str:
        return '\n'.join([
            f"{self.items[k]} -> {k}"
            for k in self.items
        ])

    def lookup(self, value: str) -> (int, str):
        hash: int = self.hash(value)
        return hash, self.items[hash] or ''
        

    

