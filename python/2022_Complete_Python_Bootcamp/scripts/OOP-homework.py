#!/bin/python3
import math

class Line():
    def __init__(self, coord1, coord2):
        self.x = []
        self.y = []
        self.x[0] = coord1[0]
        self.x[1] = coord2[0]
        self.y[0] = coord1[1]
        self.y[1] = coord2[1]
        self.coord1 = coord1
        self.coord2 = coord2


    def distance(self):
        a2 = math.exp(abs(self.x[0] - self.x[1]), 2)
        b2 = math.exp(abs(self.y[0] - self.y[1]), 2)
        c2 = a2 + b2
        distance = math.sqrt(c2)
        return distance

    def slope(self):
        rise = abs(self.y[0] - self.y[1])
        run = abs(self.x[0] - self.x[1])
        return float(rise)/float(run)

