
#include <stdio.h>
#include <math.h>
#include <stdlib.h>
#include <time.h>
#include <stdbool.h>
#include <unistd.h>
#include <string.h>

int main(){
	char string1[]="abcdefghijklmnopqrstuvwxyz";
	char string2[]="123456789";
	int result=0;

	result=(strcmp(string1, string2));
	printf("strcmp(string1, string2) returns a value of %d\n", result);
	result=(strcmp(string2, string1));
	printf("strcmp(string2, string1) returns a value of %d\n", result);


}
