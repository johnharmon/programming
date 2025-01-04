#include <stdio.h>
#include <math.h>
#include <stdlib.h>
#include <time.h>
#include <stdbool.h>
float calcAbs(float value);

int main(){

	int x;

	printf("Enter an integer to find the square root of: ");
	scanf("%d", &x);
	getchar();
//	printf("\n");

	const float epsilon = .00001;
	float guess = 1.0;
	if ( x < 0 ){
		printf("negative number blah blah\n");
		return -1.0;
	}

	while ( calcAbs(guess * guess - x) >= epsilon){
		guess = (x / guess + guess ) / 2.0;
	}

	printf("%f\n", guess);
}



float calcAbs(float value){
	if (value < 0)
		return value*-1;
	else
		return value;
}
