#include <stdio.h>
//This dumb thing just prints input

int main() {

	int x;
	char str[100];

	printf("Enter an integer followed by a string:");
	scanf("%d %s", &x, str);
	printf("\nthe integer was: %d and the string was: %s\n", x, str); 
	getchar();
	printf("enter another int followed by another string");
	scanf("%d %s", &x, str);
	printf("\nSecond int was %d and the second string was %s\n", x, str);

}
