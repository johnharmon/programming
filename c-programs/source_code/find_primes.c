#include <stdio.h>
#include <math.h>
#include <stdlib.h>
#include <time.h>


int main(){

	int primes[30] = {2,3};
	int max=100;
	int start = 4;
	int i; 
	int c;
	int arrLen= sizeof(primes)/sizeof(primes[0]);
	int insertIndex=2;	
	for (i = start; i <= max; i++){
		for (c = 0; c < arrLen; c++){	
			if (i % primes[c] == 0){
				break;
			}
			else if (primes[c] >= i/2 || primes[c] == 0){
				primes[insertIndex] = i;
				insertIndex++;
				break;
			}
			else
				continue;
		}
	}
	for ( i=0; i<arrLen; i++){
		if (primes[i] != 0 ){
			printf("%d\n", primes[i]);
		}
	}
}


/*  Instructor Solution
 *
 *














