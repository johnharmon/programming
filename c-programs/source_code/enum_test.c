#include <stdio.h>
int main(){
	enum month { January, Febuary, March, April, May, June, July, August, Septemter, October, November, December };
	enum month current_month=May;
	printf("The current month is %d\n", current_month);
}
