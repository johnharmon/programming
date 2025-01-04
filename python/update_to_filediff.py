#!/bin/python3

attributes = os.stat(self.filepath):
    for attribute in attributes.__repr__().split('(', 1)[1].split(')', 1)[0].replace(',', '').split(' '):
        self._dict['stat'][attribute.split('=')[0] = attribute.split('=')[1]

