#include <stdio.h>
#include <math.h>
#include <stdlib.h>
#include <time.h>
#include <stdbool.h>

int placeMarker(int player, int x, int y);
void drawBoard(void);
_Bool checkForWin(void);
_Bool boardIsFull(void);
char board[3][3]={
		{'_','_','_'},
		{'_','_','_'},
		{'_','_','_'}
		};

int main(){

int player1=1;
int player2=2;
int playerTurn=1;
int row,column;
while ( checkForWin() == false ){
	drawBoard();
	printf("It is player %d's turn!\n", playerTurn);
	printf("Please enter x and y coordinates to place your marker (in the form of [row column]:");
	scanf("%d %d", &row, &column);
	getchar();
	printf("\n");
	placeMarker(playerTurn, row, column);
	if (checkForWin()){
		drawBoard();
		printf("The game is over. player %d wins!\n", playerTurn);
		break;
	}
	switch (playerTurn){
		case 1:
			playerTurn++;
			break;
		case 2:
			playerTurn--;
			break;
	}
}
}

int placeMarker(int player, int row, int column){
	row--;
	column--;
	char marker;
	char locationMarker;
	if (row < 0 || row > 2 || column < 0 || column > 2){
		printf("Invalid location. row or column is out of range\n");
		return 1;
	} 
	else{ 
		locationMarker=board[row][column];
		//printf("location marker: %c\n", locationMarker);
		if (player==1){
			marker='X';
		}
		else{
			marker='O';
		}
		if (locationMarker=='_'){
			board[row][column]=marker;
			return 0;
		}
		else{
			printf("Invalid location. That spot has already been marked\n");
			return 1;
		}
	}
}

void drawBoard(void){
	int row;
	int column;
	for (row=0; row<3; row++){
		for (column=0; column<3; column++){
		printf("%c ", board[row][column]);
		}
		printf("\n");
	}
	return;
}

_Bool checkForWin(void){
	char diagChar1;
	char diagChar2;
	char markers[2]={'X','O'};
	_Bool won=false;
	_Bool possibleWin=false;
	int marker;
	for (marker=0; marker<2; marker++){
		int row;
		int column;
		for (row=0; row<3; row++){
			if (board[row][0]==markers[marker]){
				if (row == 0 || row == 2){
					diagChar1=board[abs(row-1)][1];
					diagChar2=board[abs(row-2)][2];
					if (diagChar1==markers[marker] && diagChar2==markers[marker]){
						return true;
					}
				}
				if (board[row][1]==markers[marker] && board[row][2]==markers[marker]){
					return true;
				}
			}
		}			
		
		for(column=0; column<3; column++){
			if (board[0][column]==markers[marker]){
				if (board[1][column]==markers[marker] && board[2][column]==markers[marker]){
					return true;
				}
			}
		}
	}
	if (boardIsFull()){
		printf("Board is full, the match is a draw\n");
		return false;
	} else{
	return false;
	}

}

_Bool boardIsFull(void){
	int row;
	int column;
	for (row=0; row<3; row++){
		for (column=0; column<3; column++){
			if (board[row][column]=='_'){
				return false;
			}
		}
	}
	return true;
}
//printf("%c\n", board[1][1]);



/*while ( true ){
	scanf("%d %d", &x, &y);
	//printf("You input %d and %d", x, y);
	placeMarker(1, x, y);
	getchar();
	drawBoard();
	*/

//}
/*	if (checkForWin()){
		printf("The function returned true\n");
	}
	else{
		printf("The function returned false\n");
	} */
