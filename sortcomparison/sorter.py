#!/usr/local/bin/python

class Person:
    def __init__(self, name, age): 
        self.name = name
        self.age = age
    
    def __repr__(self): 
        return "{} is {} years old".format(self.name, self.age)
    

def main():
    people = [
        Person("Jason", 28),
        Person("Sabah", 27)
    ]

    print(sorted(people, key=lambda person: person.age))

main()