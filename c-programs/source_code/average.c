#include <stdio.h>
#include <math.h>
#include <stdlib.h>
#include <time.h>
#include <stdbool.h>
#include <unistd.h>
#include <string.h>

int findGreater (int a, int b);

int main(){
	int total1 = 0;
	int total2 = 0;
	int test=0;
	int count=10000000;
	float divisor=(float)count;
	//int rand_input=system("date +%N");
	srand(time(NULL));
	rand();
	for (int i = 0; i<count; i++){
		total1+=(rand() % 10 + 1);
		total2+=( findGreater( (rand() % 10 + 1), (rand() % 10 + 1) ) );
	}
	float average1=total1/divisor;
	float average2=total2/divisor;

	printf("The average of single dice rolls was %f\n", average1);
	printf("The average of rolling twice and choosing the highest was %f\n", average2);
}


int findGreater (int a, int b){
	if (a > b){
		return a;
	} else
		return b;
}
