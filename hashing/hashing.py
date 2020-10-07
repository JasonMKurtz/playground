#!/usr/local/bin/python3

from typing import Dict, List

class HashTableLimitException(Exception):
    pass


class Table:
    items: Dict[int, List[str]] = {}
    size: int
    elements: int

    def __init__(self, size: int):
        self.size = size
        self.elements = 0
        pass

    def hash(self, value: str) -> str:
        return hash(value) % self.size
    
    def add(self, value: str) -> None:
        if self.elements >= self.size:
            raise HashTableLimitException(f"{self.size} items already exist.")

        self.elements += 1

        index: int = self.hash(value)
        if index in self.items:
            update: List[str] = self.items[index] + [value]
        else:
            update: List[str] = [value]
        
        self.items.update({index: update})

    def __repr__(self) -> str:
        return '\n'.join([
            f"{self.items[k]} -> {k}"
            for k in self.items
        ])
    
    def __len__(self) -> int:
        return self.elements

    def lookup(self, value: str) -> (int, str):
        hash: int = self.hash(value)
        items: List[str] = self.items[hash] or []
        if value not in items:
            return -1, ""

        return hash, self.items[hash] or ''
        

    

