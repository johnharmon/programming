#!/bin/python3
def maxArea(height):
    right_index = len(height)-1
    left_index = 0
    best_left = left_index
    best_right = right_index
    min_height = min(height[right_index], height[left_index])

    while left_index < right_index:
        right_index -= 1
        h = min(height[right_index], height[best_left])
        print(h)
        if (h - min_height) >= 1:
            best_right = right_index
            min_height = h
            #continue
        left_index += 1
        h = min(height[best_right], height[left_index])
        print(h)
        if (h - min_height) >= 1:
            best_left = left_index
            min_height = h
            #continue
        print(f'Best Right Idex: {best_right} Height: {height[best_right]}')
        print(f'Best left Idex: {best_left} Height: {height[best_left]}')
        print(f'Min height: {min_height}')
        print()
    return abs(best_right - best_left) * min_height



print(maxArea([2,3,10,5,7,8,9]))
