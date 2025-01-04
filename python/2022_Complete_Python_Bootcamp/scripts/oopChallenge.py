#!/bin/python3

class Account:

    def __init__(self, owner, balance):
        self.owner = owner
        self.balance = balance 
    def deposit(self, amount):
        self.balance += amount
    def withdraw(self, amount):
        if (amount <= self.balance):
            self.balance -= amount
        else:
            print(f'\033[00;33mWithdraw amount of \033[01;31m${amount}\033[00;33m exceeds available funds of \033[01;31m${self.balance}!\033[00m')
    def __str__(self):
        return f'This account is owned by {self.owner} and has an available balance of ${self.balance}'
    def __len__(self):
        return f'This account has a balance of {self.balance}'
    def __del__(self):
        print(f'\033[01;33mWarning, deleting accoutn belonging to {self.owner} with a balance of {self.balance}\033[00m')
    