#include <stdio.h>
#include <math.h>
#include <stdlib.h>
#include <time.h>

#define  MAX  100
#define  MIN  1
#define  NUM_TRIES  7

int main(){
	
//	int max=100;
//	int min=1;

	// Set rand() seed to number of epoch seconds (seconds since Jan 1, 1970)
	srand(time(NULL));
	rand();

	int magic_num = ( rand() % ( MAX - MIN + 1 )) + MIN;
	int guess;
	int tries;
	
	printf("I have generated a random number between %d and %d\n", MIN, MAX);
	printf("You have %d chances to guess the correct number\n", NUM_TRIES);
	
	for (tries = NUM_TRIES; tries > 0; tries--){
		printf("Please guess a number: ");
		scanf("%d", &guess);
		getchar();
		
		if ( guess > magic_num ){
			printf("Too high. You have %d tries left\n", tries-1);
		} else if ( guess < magic_num){
			printf("Too low. You have %d tries left\n", tries-1);
		} else {
			printf("Correct! The number was %d\n", magic_num);
			return 0;
		}
	}

	printf("Sorry loser, but you were not able to guess the correct number\n");
}

			






//	printf("Number is %d\n", magic_num);
/*
	magic_num = ( rand() % ( max - min + 1 )) + min;
	printf("Number is %d\n", magic_num);
	magic_num = ( rand() % ( max - min + 1 )) + min;
	printf("Number is %d\n", magic_num);                       			
	magic_num = ( rand() % ( max - min + 1 )) + min;
	printf("Number is %d\n", magic_num);             
	magic_num = ( rand() % ( max - min + 1 )) + min;
	printf("Number is %d\n", magic_num);       

*/



	
