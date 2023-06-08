#!/bin/python3

import random

class card():

    def __init__(self, suite, value):
        self.suite = suite
        self.value = value
        if self.suite == 'Clubs' or self.suite == 'Spades':
            self.color = 'Black'
        else:
            self.color = 'Red'
    def __str__(self):
        return f'{self.suite} of {self.value}'

class card_deck():

    def __init__(self):
        suites = ['Hearts', 'Diamonds', 'Spades', 'Clubs']
        values = {'Two': 2, 'Three': 3, 'Four': 4, 'Five': 5, 'Six': 6, 'Seven': 7, 'Eight': 8, 'Nine': 9, 'Ten': 10, 'Jack': 10, 'Queen': 10, 'King': 10, 'Ace': 11}
        unshuffled_cards = []
        for suite in suites:
            for key in values:
                unshuffled_cards.append(card(suite = suite, value = key))
     
    @classmethod
    def get_value(self, card):
        return self.values[card.value]

    @classmethod    
    def shuffle(self):
        self.shuffled_cards = random.sample(self.unshuffled_cards, len(self.unshuffled_cards))
    
    @classmethod
    def get_card(self):
        pulled_card = self.shuffled_cards.pop([-1])
        return pulled_card 

class player():
    def __init__(self, money = 100):
        self.money = money
        self.hand = []
        self.hand_value = 0

    @classmethod
    def bet(self, bet = 10):
        if bet <= self.money:
            return bet 
        else:
            return -1
    
    @classmethod
    def hit(self, card_deck):
        self.hand.append(card_deck.get_card())
        #self.hand_value += card_deck.get_value(self.hand[-1])
        return self.hand_value
    
    def new_hand(self):
        self.hand.clear()
        self.hand_value = 0

    def calculate_hand(hand = []):
        has_ace = False
        num_aces = 0
        total = 0
        for card in hand:
            if card.value == 'Ace':
                has_ace = True
                num_aces += 1
            total += card_values[card.value]
        if total > 21:
            if has_ace:
                total -= num_aces*10
        return total
        
class dealer(player):

    @classmethod
    def show(self):
        return self.hand[0]
            
global card_suites
global card_values

card_suites = ['Hearts', 'Diamonds', 'Spades', 'Clubs']
card_values = {'Two': 2, 'Three': 3, 'Four': 4, 'Five': 5, 'Six': 6, 'Seven': 7, 'Eight': 8, 'Nine': 9, 'Ten': 10, 'Jack': 10, 'Queen': 10, 'King': 10, 'Ace': 11}

my_dealer = dealer(money = 0)
my_player = player(money = 100)

def new_hand():
    this_deck = card_deck()
    this_deck.shuffle()

    my_dealer.hit()
    my_dealer.hit()
    my_player.hit()
    my_player.hit()

    print(f'Dealer shows {my_dealer.show()}')
    print(f'You have {my_player.hand[0]} and {my_player.hand[1]} which totals {my_player.calculate_hand()}')

def game_loop():
    if my_dealer.calculate_hand() < 17:
        my_dealer.hit()
    print(f'Dealer shows {my_dealer.show()}')
    print(f'You have: ')
    for my_card in my_player.hand:
        print(f'{my_card}')
    print(f'Which totals {my_player.calculate_hand()}')

    response = input('Would you like to hit? [hit/stay]')
    if response == 'hit':
        my_player.hit()
    elif response == 'stay':
        return False
    else:
        print('Please type "hit" or "stay"!')
        return True

    
def check_winner():
    dealer_hand = my_dealer.calculate_hand()
    player_hand = my_player.calculate_hand()
    if player_hand > 21:
        return False
    elif dealer_hand > 21:
        return True 
    elif dealer_hand >= player_hand:
        return False
    else:
        return True




#class Game():
#    def __init__(self):








    
        



        


