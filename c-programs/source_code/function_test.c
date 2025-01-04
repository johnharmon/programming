#include <stdio.h>
#include <math.h>
#include <stdlib.h>
#include <time.h>
#include <stdbool.h>

int GCD(int a, int b); // Find the Greatest Common Denominator for two ints

float calcAbs(float x); // Find the absolute (positive) value for given float

float calcSqrt(int x); // Find the square root of a number

int findGreater (int a, int b);

int largestWholeSquare(int x);

int main(){ 

	int sqrtN;
	int gcdA; 
	int gcdB;
	int gcdResult;
	float absF;
	float sqrtF;
	int closestSqrt;
	sqrtN = 823809;
//	printf("Enter two numbers to find the greatest common demonator for: ");
//	scanf("%d %d", &gcdA, &gcdB);
//	getchar();
//	gcdResult = GCD(gcdA, gcdB);
//	printf("The GCD is %d\n", gcdResult);
//	absF = calcAbs(-233.58);
//	printf("the abs of -233.58 is %f\n", absF);
//

	//closestSqrt = largestWholeSquare(55);
	//printf("The largest whole square root under 55 is %d\n", closestSqrt);
	sqrtF = calcSqrt(sqrtN);
	printf("The sqrt of %d is %f\n", sqrtN, sqrtF);

}


int findGreater(int a, int b){
	if (a > b){
		return a, b;
	}
	else 
		return b, a;
}

int GCD(int a, int b){

	int big;
	int small;

	if (a > b){
		big=a;
		small=b;
		}
	else if (b > a){
		big=b;
		small=a;
	}
	else {
		printf("both integers are equal, their GCD will always be themselves\n");
		return 0;
	}
	while ( 0 == 0){
		int tempBig = big;
		if ( big % small == 0 ) {
			return small;
		}
		else { 
			big = small;
			small = tempBig % small;
		}
	}
}

float calcAbs(float a){
	if (a < 0){
		return a*(-1.0);
	}
	else
		return a;
}


float calcSqrt(int x){
	int precision = 5; // how many decimal places we will calculate the square root to
	float result;
	int remainder = x;
	float addToResult;
	int wholeSqrt;
	int outer;
	int inner;
	wholeSqrt = largestWholeSquare(x);
	result += wholeSqrt;

	for (outer=0; outer<=precision; outer++){
	
		addToResult = 1/pow(10, outer+1);
		for (inner = 0; inner < 10; inner++){
			if (pow((result + addToResult), 2) <=x){
				result+=addToResult;
			}
			else {
				break;
			}
		}
	}
	return result;
// ######################################Instructor Solution (couldn't find this algorithm online)#######################################################
// const float epsilon = .00001;
// float guess = 1.0;
// if ( x < 0 ){
// 	printf("negative number blah blah\n");
// 	return -1.0;
// }
//
// while ( calcAbs(guess * guess - x) >= epsilon){
// 	guess = (x / guess + guess ) / 2.0;
// }
//
// return guess;

}

//	addToResult = wholeSqrt /  (pow(10, i));
//	printf("%d\n", wholeSqrt);
//	printf("the result of the largest square divisoin is %f\n", addToResult);
//	result += addToResult;
//	remainder = remainder % (wholeSqrt * wholeSqrt);
//	x = (2 * wholeSqrt + remainder);
//	printf("remainder on iteration %d is %d\n", i, remainder);



int largestWholeSquare(int x){
	
	int min = 0;
	int max = x;
	int mid;
	unsigned long int squared;
	int result;
	while (min <= max){
		mid = (max + min)/2;
//		printf("mid: %d\n", mid);
		squared = mid * mid;
		if (squared == x){
			return squared;
		}
		else {
			if (squared > x){
				max = mid-1;
			}
			else{
				result = mid;
				min = mid + 1;
		}
	
	}
//	printf("max: %d\n", max);
	//printf("min: %d\n", min);
}
	return result;
/*				
	int squared;
	int i;
	for (i=0; i*i < x; i++){
		squared = i*i;
	}
	return i-1;
*/
}



////	printf("big: %d\nsmall:%d\n", big, small);
//	if ( big % small == 0){
////		printf("The GCD is %d\n", small);
//		return small;
//	}
//	else {
////		printf("Reached for loop\n");
//		int i;
//		for( i = small /2; i>=1; i--){
//		//	printf("%d\n", i);
//			if (small % i == 0){
////				printf("found divisor for the smaller int\n");
//				if (big % i == 0){
//					return i;
		//		}
		//	}
		//}
	//}


	

