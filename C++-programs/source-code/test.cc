#include <vector>
#include <iostream>
#include <algorithm>
using namespace std;

int main(){
	vector<int> nums = {6,5,3,3,5,1};

	sort(nums.begin(), nums.end());
	int i;
	for (i=0; i<nums.size(); i++){
	cout << nums[i] << endl;
	}
}
