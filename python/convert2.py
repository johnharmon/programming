#!/bin/python3
def convert(s, numRows):
    if numRows == 1:
        return s
    else:
        rows = ['']*numRows
        decrementing = False
        previous_row_number = 0
        for index in range(0,len(s)):
            row_number = index % numRows
            if decrementing == True:
                row_number = previous_row_number-1
            elif decrementing == False and index !=0:
                row_number = previous_row_number + 1
            if (row_number == 0 and index != 0) or (row_number == numRows-1):
                decrementing = not decrementing
            rows[row_number] += s[index]
            previous_row_number = row_number
            print(f'row number: {row_number}, index: {index}')
        return ''.join(rows)
#    rows = ['']*numRows
#    decrementing = False
#    for index in range(0,len(s)):
#        row_number = index % numRows
#        if row_number == 0 and index != 0:
#            decrementing = not decrementing
#        if decrementing == True and previous_row_number != 0:
#            row_number = previous_row_number - 1
#        rows[row_number] += s[index]
#        previous_row_number = row_number
#        print(f'row number: {row_number}, index: {index}')
#    return ''.join(rows)
#
#

print(convert('PAYPALISHIRING', 3))
