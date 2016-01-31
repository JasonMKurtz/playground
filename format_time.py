#!/usr/bin/python
import time

def military(f): 
    def inner(*args, **kwargs): 
        now = f(*args, **kwargs)
        postnoon = now[-2:].lower() == "pm"
        hour, minute, seconds = now.split()[0].split(":")
        seconds = seconds[0:2]
        hour = int(hour)
        if not postnoon and hour == 12: 
            hour = "00"
        elif postnoon and hour < 12: 
            hour += 12

        return "%s:%s:%s" % (str(hour).zfill(2), minute, seconds)
    return inner

def words(f): 
    def inner(*args, **kwargs): 
        now = f(*args, **kwargs)
        hour, minute, seconds = now.split()[0].split(":")
        postnoon = -1
        if len(seconds) > 2: 
            seconds = seconds[0:2]
            postnoon = now[-2:].lower() == "pm"

        hr  = num2words(hour)
        if int(minute) == 0: 
            min = "o'clock"
        elif int(minute) < 10: 
            min = "oh " + num2words(minute)
        else: 
            min = num2words(minute)

        extra = ""
        if postnoon != -1: 
            if postnoon: 
                if hour >= 5: 
                    extra = "in the evening"
                else: 
                    extra = "in the afternoon"
            else: 
                extra = "in the morning" 
        
        return "%s %s %s" % (hr, min, extra)

    return inner


def num2words(num): 
    words = { 
        0: 'oh', 1: 'one', 2: 'two', 3: 'three', 4: 'four', 5: 'five', 6: 'six', 7: 'seven', 8: 'eight', 9: 'nine', 
        10: 'ten', 11: 'eleven', 12: 'twelve', 13: 'thirteen', 20: 'twenty', 30: 'thirty', 40: 'forty', 50: 'fifty'
    }
    
    num = int(num)

    if num <= 13: 
        ret = words[num]
    elif num >= 14 and num < 20: 
        if int(num) == 18: 
            ret = "eighteen"
        else:
            ret = words[int(num) - 10] + "teen"
    else: 
        mod = (num % 10) 
        if (mod == 0): 
            ret = words[num]
        else: 
            ret = words[(num / 10) * 10] + " " + words[mod]
                 
    return ret
            
@words
def what_time(str=""):
    return time.strftime("%I:%M:%S%p") if str == "" else str

print what_time()
"""
07:32:40PM
@military 19:32:40
@words seven thirty two in the evening
@words, @military nineteen thirty two
"""
