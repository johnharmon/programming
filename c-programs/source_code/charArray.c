#include <stdio.h>
#include <math.h>
#include <stdlib.h>
#include <time.h>
#include <stdbool.h>
#include <unistd.h>
#include <string.h>


int strLength(char str[]);
void cat(char str1[], char str2[]);
_Bool cmp(char str1[], char str2[]);

int main(){


int x;
char string1[]="abcdefghijklmnopqrstuvwxyz";

x=strLength(string1);

printf("%d\n", x);

cat("hello my name is", "john");



}


int strLength(char str[]){

	int length=0;

	for (int i=0; str[i]!='\0'; i++){
		length=i;
	}
	return length+1;
}


void cat(char str1[], char str2[]){
	int length=(strLength(str1)+strLength(str2))+1;
	char result[length];
	char resultIndex=0;

	for (int i=0; i<strLength(str1); i++){
		result[resultIndex]=str1[i];
		resultIndex++;
	}
	for (int i=0; i<strLength(str2); i++){
		result[resultIndex]=str2[i];
		resultIndex++;
	}
	result[length-1]='\0';
	printf("%s\n", result);
}


