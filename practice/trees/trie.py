#!/usr/local/bin/python3

from typing import List

class TrieNode:
    def __init__(self):
        self._end = "*"
        self.trie = self.make_trie("hello", "hi", "hey")
        pass

    def make_trie(self, *words):
        trie = dict()
        for word in words:
            temp_dict = trie
            for letter in word:
                temp_dict = temp_dict.setdefault(letter, {})
            temp_dict[self._end] = self._end
        return trie

    def contains(self, prefix: str) -> bool:
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