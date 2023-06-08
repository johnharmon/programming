#include <stdio.h>
#include <stdlib.h>
#include <math.h>


int main() {

	int x = 5;

	switch ( x ) {

		case 2:
			printf("The value is 2\n");
			break;

		case 3:
			printf("The value is 3\n");
			break;

		default:
			printf("The value was not 2 or 3\n");
			break;
	}
}
