#!/bin/python3
def convert(s, numRows):
    multi_row = True
    single_rows = numRows//2
    rows = list()
    columns = list()
    column_length = numRows
    column_offsets = list()
    if numRows > 3:
        for index in range(1,numRows-1):
            column_offsets.append(index)
        column_offsets.reverse()
        
        #print(column_offsets)
    elif numRows == 3:
        column_offsets.append(1)
    elif numRows == 2:
        row1 = ''
        row2 = ''
        for index in range(0,len(s)):
            if index % 2 == 0:
                row1 += s[index]
            else:
                row2 += s[index]
        return row1+row2
    while(len(s)>0):
        columns.append(s[0:column_length])
        s = s[column_length:]
        for offset in column_offsets:
            next_column = [' '] * numRows
            for index in range(0,column_length):
                if len(s)>0:
                    if index == offset:
                        next_column[index]=s[0]
                        s = s[1:]
                        columns.append(''.join(next_column))
                        break
    answer = ''
   # for row in zip(*columns):
   #    print(row)
    for index in range(0,numRows):
        for column in columns:
            if index < len(column):
                if column[index] != ' ':
                    answer += column[index]
    return answer


print(convert("PAYPALISHIRING", 5))
