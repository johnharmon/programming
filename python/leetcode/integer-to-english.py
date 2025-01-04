#!/bin/python3
class Solution:
    def parse_tens(self, num_to_word, number):
        if number == 0:
            return ''
        if number < 20:
            prefix_number = num_to_word[str(number)]
            prefix_number += ' '
            return prefix_number
        prefix_number = ''
        tens_place = number//10
        prefix_number = num_to_word[f'{tens_place}0']
        prefix_number += ' '
        if number % 10 != 0:
            prefix_number += f'{num_to_word[str(number % 10)]} '
        return prefix_number

    def parse_hundreds(self, num_to_word, number):
        prefix_number = ''
        hundreds_place = number//100
        tens_place = number//10
        prefix_number = f'{num_to_word[str(hundreds_place)]} Hundred '
        if tens_place > 0 :
            prefix_number += self.parse_tens(num_to_word, number % 100)
        return prefix_number
    
    def numberToWords(self, num: int) -> str:
        if num == 0:
            return 'Zero'
        word_to_num = {
        'One': 1,
        'Two': 2,
        'Three': 3,
        'Four': 4,
        'Five': 5,
        'Six': 6,
        'Seven': 7,
        'Eight': 8,
        'Nine': 9,
        'Ten': 10,
        'Eleven': 11,
        'Twelve': 12,
        'Thirteen': 13,
        'Fourteen': 14,
        'Fifteen': 15,
        'Sixteen': 16,
        'Seventeen': 17,
        'Eighteen': 18,
        'Nineteen': 19,
        'Twenty': 20,
        'Thirty': 30,
        'Forty': 40,
        'Fifty': 50,
        'Sixty': 60,
        'Seventy': 70,
        'Eighty': 80,
        'Ninety': 90,
        'Hundred': 100,
        'Thousand': 1000,
        'Million': 1000000,
        'Billion': 1000000000
        }
        word_names = list(word_to_num.keys())
        word_names.reverse()
        num_to_word = dict()
        for key in word_names:
            num_to_word[str(word_to_num[key])] = key 
        number_counts = dict()
        for word in word_names:
            number_counts[word] = 0 
        previous_index = 0
        while num > 0:
            for index in range(previous_index,len(word_names)):
                word = word_names[index]
                if num >= word_to_num[word]:
                    number_counts[word] = num//word_to_num[word]
                    num = num % word_to_num[word]
                    previous_index = index
                    break
        number_string = ''
        for word in word_names:
            if number_counts[word] > 0:
                if word_to_num[word] >= 100:
                    if number_counts[word] < 20 and number_counts[word] > 0:
                        prefix_number = str(num_to_word[str(number_counts[word])])
                        number_string += f'{prefix_number} '
                    elif number_counts[word] < 100:
                        prefix_number = self.parse_tens(num_to_word, number_counts[word])
                        number_string += f'{prefix_number}'
                    else:
                        prefix_number = self.parse_hundreds(num_to_word, number_counts[word])
                        number_string += f'{prefix_number}'
                number_string += f'{word} '
        number_string = number_string.strip()
        return number_string

