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
	int balancedCounter=0;
	int longestBalancedOpen=0;
	int length=strlen(parens);
	int openAfterBalanceStart=0;
	_Bool balancedStarted=false;
	int index;
	for (index=0; index<length; index++){
		if (parens[index]=='(' && balancedStarted==false){
			unbalancedOpen++;
		}
		else if (balancedStarted==true && parens[index]=='('){
			openAfterBalanceStart++;
			balanceCounter++;
		}
		else if (balancedStarted==true && parens[index]==')'){
		
			openAfterBalanceStart--;
			balanceCounter++;
			if (openAfterBalanceStart==0){
				longestBalancedOpen=balanceCounter;
			}
			else if (openAfterBalanceStart==-1){
				longestBalancedOpen=balanceCounter+2;
				balanceStarted=false;
				openAfterBalanceStart=0;
				balanceCounter=0;
				}
		}
		else if (balanceStarted==false && parens[index]==')' && unbalancedOpen!=0){
			
			balanceStarted=true;
			unbalancedOpen--;
			balanceCounter+=2;
			balancedOpen=balanceCounter;

		}
	}
	return longestBalancedOpen*2;
}

			//if (unbalancedOpen > 0 && balancedStarted==false){
			//	balancedStarted=true;
			//	balancedOpen++;
			//	if (balancedOpen > longestBalancedOpen){
			//		longestBalancedOpen=balancedOpen;
			//	}
			//	unbalancedOpen--;
			//	if (unbalancedOpen == 0){
			//		balancedOpen=0;
			//	}
			//
			//}
			//else {
			//	balancedOpen=0;
			//}
