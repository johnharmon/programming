#include <stdio.h>
#include <math.h>
#include <stdlib.h>
#include <time.h>
#include <stdbool.h>
#include <unistd.h>
#include <string.h>

int longestBalance(char parens[100]);

int main(){
int x;
char input[100];


while (true){

	scanf("%s", input);
	x=longestBalance(input);
	printf("%d\n", x);
	getchar();
	printf("%s\n", input);
	}
}

int longestBalance(char parens[100]){
	int unbalancedOpen=0;
	int balancedOpen=0;
	int longestBalancedOpen=0;
	int length=strlen(parens);
	int index;
	for (index=0; index<length; index++){
		if (parens[index]=='('){
			unbalancedOpen++;
		} else if (parens[index]==')'){
			if (unbalancedOpen > 0){
				balancedOpen++;
				if (balancedOpen > longestBalancedOpen){
					longestBalancedOpen=balancedOpen;
				}
				unbalancedOpen--;
				if (unbalancedOpen == 0){
					balancedOpen=0;
				}
			}
			else {
				balancedOpen=0;
			}

		}
	}
	return longestBalancedOpen*2;
}
