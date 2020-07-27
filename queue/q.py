#!/usr/local/bin/python3

class Q:
    def __init__(self): 
        self.items = []

    def Add(self, value):
        print(f"Adding {value} to the queue.")
        self.items.append(value)

    def __len__(self) -> int:
        return len(self.items)


class QL(Q):
    def Dequeue(self) -> str:
        if not len(self.items):
            return ""
        latest = self.items[-1]
        self.items.pop()
        return latest

class QF(Q):
    def Dequeue(self) -> str:
        if not len(self.items):
            return ""
        
        first = self.items[0]
        self.items = self.items[1:]
        return first
    
if __name__ == '__main__': 
    q = QF()
    q.Add("foo")
    q.Add("bar")
    print(q.Dequeue())
    print(q.Dequeue())
    print(len(q))
    print(q.Dequeue())



