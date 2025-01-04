#include <stdio.h>
#include <math.h>
#include <stdlib.h>
#include <string.h>

int testInt;

void testFunction();
float squareRoot(float x);
float absoluteValue(float x);

int main(){
	
	/*
	double result = pow(2.0, 8.0);
	
	//printf("%lf\n", result);
	//
	char digits[100] = "1110001001";
	//printf("%d\n", (int) digit);
	
	printf("%c\n", digits[2]);
	
	//printf("%d\n", atoi( *digits[2]);
	
	char str[2];
	str[0]=digits[2];
	str[1]='\0';
	int converted = atoi(str);

	//int converted = atoi(digits);

	printf("Digit is %d\n", converted);
	*/

	/* int a = 127;
	char numString[100] = (char) a;
	printf("The number is %s\n", numString);
	*/


	/*
	int x = 1; 
	int y = 2;
	int z = 3;
	int res = (x == y != z);
	printf("Result is %d\n", res);
	res = ( 1 == 2 );
	printf("Result is %d\n", res);
	*/
	 /*int x=7-45;

	int y=abs(x);
	printf("%d\n", y);*/
	/*const char test[]="henlo world";
	//scanf("%s", test);
	//strcpy(test, "asdfaiodjnini");
	//getchar();
	//printf("%s\n", test);
	int len1=strlen(test);
	int len2=sizeof(test)/sizeof(test[0])-1;
	printf("len1 = %d, and len2 = %d\n", len1, len2);
	printf("please listen to this: (%s) message from the program\n", test);*/

	char string1[10]="first";
	char string2[]="secondasdfasdfasfdasf";
	char string3[100];
	int x = strlen(string1);
	int y = sizeof(string1);
	int z=y;
	printf("x=%d\n", x);
	printf("z=%d\n", z);
	strncat(string1, string2, 100);
	//strncat(string1, string2, (sizeof(string1)-strlen(string1)));
	printf("%s\n", string1);
//	strcat(string1, string2);
//	printf("%s\n", string1);

//	int x = sizeof(int);
//	printf("%d\n", x);
//	
//	x=sizeof(long);
//	printf("%d\n", x);
//
//	x=sizeof(double);
//
//	printf("%d\n", x);
//
//	x=sizeof(char);
//
//	printf("%d\n", x);
//
//	x=sizeof(float);
//
//	printf("%d\n", x);
//
//	x=sizeof(_Bool);
//
//	printf("%d\n", x);
//
//
//
//	int x;
//	x = (int) pow(2, 31);
//	printf("x = %d\n", x);
//	x = (int) -pow(2,31);
//	//x += 500;
//
//	printf("x now = %d\n", x);
//	
//	unsigned int y =(unsigned int) pow(2,32);
//	printf("Y = %u\n", y);
//

	
//	testFunction();

//	char x[100] = "10";
//	testFunction(x);
//
//	printf("%d was converted to an int\n", testInt);
//	testInt = testInt*3;
//	printf("%d\n", testInt);
//
//

//int x = 65;

//printf("%c\n", x);
//char x[3][3]={
//		{'a','b','c'},
//		{'d','e','f'},
//		{'g','h','i'}
//	};
//char y=x[1][2];
//printf("%c\n", y);


}
//void testFunction(char input[100]){
//	testInt = atoi(input);
//	printf("This is your first function!\n");
//	return;
//}
//
//
//float squareRoot(float x){
//
//	const float epsilon = .00001;
//	float guess = 1.0;
//
//	while (absoluteVAlue (guess * guess -x) >= epsilon){
//		guess = (x /guess + guess
