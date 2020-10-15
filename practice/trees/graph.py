#!/usr/local/bin/python3

class Node:
    neighbors = []
    value: int = 0
    def __init__(self, value):
        self.value = value
    
    def addNeighbors(self, neighbors):
        if not len(self.neighbors):
            self.neighbors = neighbors
            return

        self.neighbors += neighbors
