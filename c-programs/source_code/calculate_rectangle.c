#include <stdio.h>

int main (){

	double height;
	double width;
	double perimeter;
	double area;

	printf("Enter the height of the rectangle: ");
	scanf("%lf", &height);
	getchar();
	printf("Enter the width of the rectangle: ");
	scanf("%lf",&width);
	getchar();

	perimeter = (height * 2) + (width * 2);
	area = (height * width);

	printf("height %lf   width %lf\n", height, width);
	printf("The perimeter is: %lf\n", perimeter);
	printf("The area is : %lf\n", area);

	return 0;

}
