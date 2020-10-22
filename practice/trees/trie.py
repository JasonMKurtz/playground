#!/usr/local/bin/python3

from typing import List

class TrieNode:
    def __init__(self, *words):
        self._end = "*"
        self.trie = self.make_trie(*words)
        self.words: List[str] = []

    def make_trie(self, *words):
        trie = dict()
        for word in words:
            temp_dict = trie
            for letter in word:
                temp_dict = temp_dict.setdefault(letter, {})
            temp_dict[self._end] = self._end
        return trie
    
    def traverse(self, tree = None, path: List[str] = []) -> str:
        if tree is None:
            tree = self.trie
            
        if tree == self._end:
            self.words.append(''.join(path[:-1]))
            return
        
        for char, subtree in tree.items():
            self.traverse(subtree, path + [char])
        
    def search(self, prefix: str) -> List[str]:
        if not len(self.trie):
            return []
        
        sub_tree = self.trie
        for char in prefix:
            sub_tree = sub_tree.get(char)
            if not sub_tree:
                return []
            continue
    
        self.traverse(sub_tree)
        return [prefix + word for word in self.words]
        
    def contains(self, prefix: str) -> bool:
        """Does our trie contain the prefix we've searched for?"""
        if not len(prefix) or not len(self.trie):
            return False

        if self.trie.get(prefix):
            return True
        
        sub_tree = self.trie
        for char in prefix:
            sub_tree = sub_tree.get(char)
            if sub_tree is not None:
                continue
            return False
            
        return True