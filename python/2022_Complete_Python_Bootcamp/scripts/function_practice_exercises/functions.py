#!/usr/bin/python3

import math

##Warmup Functions---{{{
#def lesser_of_two_evens(a,b):
#    if a%2 == 0 and b%2 == 0:
#        if a < b:
#            return a
#    else:
#        if a > b:
#            return a
#    return b
#print(lesser_of_two_evens(2,4))
#print(lesser_of_two_evens(2,5))
#
#def animal_crackers(text):
#    letters=[word[0] for word in text.split(' ')]
#    if letters[0] == letters[1]:
#        return True
#    else:
#        return False
#print(animal_crackers('Levelheaded Llama'))
#print(animal_crackers('Crazy Kangaroo'))
#
#def makes_twenty(a,b):
#    if a == 20 or b == 20:
#        return 20
#    else:
#        if a+b == 20:
#            return True
#    return False
#print(makes_twenty(20,10))
#print(makes_twenty(2,3))
##}}}
#
##Level 1 Functions---{{{
#def old_macdonald(name):
#    part1=name[0:3].capitalize()
#    part2=name[3:].capitalize()
#    answer=part1+part2
#    return answer
#print(old_macdonald('macdonald'))
#
#def master_yoda(sentence):
#    words=sentence.split()
#    words.reverse()
#    return " ".join(words)
#print(master_yoda("here i am"))
#
#def almost_there(num):
#    ans1=abs(100-num)
#    ans2=abs(200-num)
#    if ans1 <=10 or ans2 <=10:
#        return True
#    return False
#print(almost_there(104))
#print(almost_there(150))
#print(almost_there(200))
##}}}
#
##Level 2 Functions---{{{
#
#def has_33(numbers):
#    previous_3=False
#    for number in numbers:
#        if number == 3:
#            if previous_3:
#                return True
#            previous_3 = True
#        else:
#            previous_3 = False
#    return False
#print(has_33([1, 3, 3]))
#print(has_33([1, 3, 1, 3]))
#print(has_33([3, 1, 3]))
#
#def paper_doll(word):
#    letters=[letter*3 for letter in word]
#    answer=''.join(letters)
#    return answer
#print(paper_doll('Hello'))
#print(paper_doll('Mississippi'))
#
#
#def blackjack(a,b,c):
#    sum = a+b+c
#    if sum <= 21:
#        return sum
#    else:
#        if a == 11 or b == 11 or c == 11:
#            sum-=10
#        if sum > 21:
#            return 'Bust'
#        return sum
#print(blackjack(5,6,7))
#print(blackjack(9,9,9))
#print(blackjack(9,9,11))
#
#def summer_69(nums):
#    ignore_nums=False
#    answer=0
#    for number in nums:
#        if number == 6:
#            ignore_nums=True
#        if number == 9:
#            ignore_nums=False
#        if not ignore_nums:
#            answer+=number
#    return answer
#print(summer_69([1, 3, 5]))
#print(summer_69([4, 5, 6, 7, 8, 9]))
#print(summer_69([2, 1, 6, 9, 11]))
##}}}
#
##Challenge Functions---{{{
#
#def spy_game(numbers):
#    for index in range(1,len(numbers)-1):
#        if numbers[index] != 0:
#            continue
#        else:
#            if numbers[index-1] == 0 and numbers[index+1] == 7:
#                return True
#    return False
#print(spy_game([1,2,4,0,0,7,5]))
#print(spy_game([1,0,2,4,0,5,7]))
#print(spy_game([1,7,2,0,4,5,0]))
#print(spy_game([1,7,2,0,4,5,0]))
#print(spy_game([1,7,2,0,0,0,7]))
#print(spy_game([0,0,7,0,4,5,0]))
#
#def count_primes(max_number):
#    num_primes=0
#    for number in range(1,max_number+1):
#        if number == 1:
#            continue
#        elif number == 2 or number == 3:
#            num_primes+=1
#        else:
#            if (math.factorial(number-1)%number == number-1):
#                num_primes+=1
#    return num_primes
#print(count_primes(100))
        

def print_big(letter):
    letter_patterns = {
            'a' : "00100 01010 01110 01010 01010", 
            'b' : "10100 10001 10010 10001 10100",
            'c' : "00001 00100 10000 00100 00001",
            'd' : "10100 10001 10001 10001 10100",
            'e' : "11111 10000 11100 10000 11111",
            'f' : "11111 10000 11100 10000 10000",
    }
    control_string=letter_patterns[letter]
    #print(control_string)
    printed_row=''
    for row in control_string.split():
        #print(row)
        for character in row:
            if character == '1':
                printed_row+='*'
            else:
                printed_row+=' '
        print(printed_row)
        printed_row=''


print_big('f')

 










