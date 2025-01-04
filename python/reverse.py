#!/bin/python3
def reverse(x):
    sign = ''
    string_x = str(x)
    print(string_x)
    string_32 = '2147483647'
    if string_x[0] == '-':
        sign = '-'
        string_x = string_x.strip('-')
    rev_x = string_x[::-1]
    print(rev_x)
    if len(string_32) < len(rev_x):
        print('string too long')
        return 0
    elif len(string_32) > len(string_x):
        return int(f'{sign}{rev_x}')
    else:
        for index in range(0,len(rev_x)):
            if int(rev_x[index]) > int(string_32[index]):
                print(f'{rev_x[index]} > {string_32[index]}')
                return 0
        return int(f'{sign}{rev_x}')

print(reverse(-2143847412))
