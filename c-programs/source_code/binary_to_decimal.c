#include <stdio.h>
#include <math.h>
#include <string.h>
#include <stdlib.h>


char *toString(char c){

	char string[2] = { c, '\0' };
	char *returnString = string;
	printf("return string is: %s\n", returnString);
	return returnString;
}



int main(){

	//char digit[2];
	char binary[64];
	int result;
	printf("Enter a binary number to convert (max 64 bits)\n");
	scanf("%s", binary);
	int length = strlen(binary);

	for (int i=length-1; i >=0; i--){

		//printf("Index number: %d\n", i);
		int exponent=(length-1)-i;
		//printf("exponent for this iteration is: %d\n", exponent);
		//digit[0]=binary[i];
		//digit[1]='\0';
	//	char *digit = toString(binary[i]);
	//	printf("Value after return is: %s\n", digit);
		char digit[2] = { binary[i], '\0' };
		int converted_digit = atoi(digit);
		//printf("Value after integer conversion is: %d\n", converted_digit);
		if (converted_digit == 1){
			result += (int) pow(2.0, (double) exponent);
		}
		
		//free(digit);
	}

	printf("Total is : %d\n", result);
	
}
