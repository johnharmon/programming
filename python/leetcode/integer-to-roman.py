class Solution:
    def intToRoman(self, num: int) -> str: 
        numerals = {
            1000: 'M',
            900: 'CM',
            500: 'D',
            400: 'CD',
            100: 'C',
            90: 'XC',
            50: 'L',
            40: 'XL',
            10: 'X',
            9: 'IX',
            5: 'V',
            4: 'IV',
            3: 'III',
            2: 'II',
            1: 'I',
        }
        numeral = ''
        for key in numerals.keys():
            occurances = num//key
            letters = numerals[key]*occurances
            numeral += letters
            num -= occurances*key  
        return numeral
