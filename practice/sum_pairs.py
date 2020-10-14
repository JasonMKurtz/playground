#!/usr/local/bin/python3 

"""
Determine whether pairs exist in a given list such that they add to a given sum.
If given [1, 2, 2, 4] and a sum of 5, we want to return True, because 4 + 1 = 5.
But, if the sum if 7, return False, because no two items in the list sum to 7.
"""

from typing import List, Set

def pairs_exist(nums: List[int], target: int) -> bool:
    complements: Set[int] = set()
    for num in nums:
        complement = target - num
        if num in complements:
            return True
        complements.add(complement)
    return False

"""
>>> from sum_pairs import * 
>>> pairs_exist([1, 2, 4], 6) 
True
>>> pairs_exist([1, 2, 4, 5], 9) 
True
>>> pairs_exist([1, 2, 4, 5], 4) 
False
>>> pairs_exist([-1, 2, 4, 5], 4) 
True
"""



        


