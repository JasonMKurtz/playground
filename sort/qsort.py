#!/usr/bin/python
def qsort(arr):
    if len(arr) <= 1:
        return arr
    left = []
    right = []
    equal = []
    for i in range(len(arr)):
        if arr[i] < arr[0]:
            left.append(arr[i])
        elif arr[i] > arr[0]:
            right.append(arr[i])
        elif arr[i] == arr[0]:
            equal.append(arr[i])
    return qsort(left) + equal + qsort(right)

l = [1, 5, 3, 3, 4]
print(qsort(l))
