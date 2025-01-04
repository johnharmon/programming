#!/bin/python3
import sys

board=[['_','_','_'],
       ['_','_','_'],
       ['_','_','_']
       ]

player_markers=['X','O']
player_turn=0
num_moves=0

def clear_board():
    global board
    for row in range(0,3):
        for column in range(0,3):
            board[row][column]='_'

def print_board():
    global board
    for row in board:
        print(' '.join(row))

def check_win():
    global board
    #win_sequences=[[(0,0),(1,0),(2,0),
    if (board[1][1] != '_') and (board[0][0] == board[1][1] and board[1][1] == board[2][2]) or (board[0][2] == board[1][1] and board [1][1] == board[2][0]):
        return True
    for starting_index in range(0,3):
        if (board[starting_index][0] != '_' and board[starting_index][0] == board[starting_index][1] and board[starting_index][1] == board[starting_index][2]):
            return True
        elif (board [0][starting_index] != '_' and board[0][starting_index] == board[1][starting_index] and board[0][starting_index] == board[2][starting_index]):
            return True
        else:
            return False

def make_move(row,column):
    if num_moves > 9:
        sys.exit("")
    global player_turn
    global play_markers
    global num_moves
    if board[row][column] == '_':
        board[row][column] = player_markers[player_turn%2]
        num_moves+=1
        player_turn+=1
        if num_moves >= 5:
            if check_win():
                print("Player "+ str((player_turn-1)%2+1) +" wins!")
                print_board()
                sys.exit("")
        return True
    else:
        print("row: "+row+" column: "+column+" is not a valid choice. Please pick a empty location")
        get_move()
        
def get_move():
    print_board()
    position = input("Please choose a move to make in the form of row,column: ")
    row = int(position[0])-1
    column = int(position[-1])-1
    if len(position) > 3 or row > 3 or column > 3:
        print(position+" is not a valid choice!")
        get_move()
    else:
        make_move(row,column)

while True:
    get_move()




