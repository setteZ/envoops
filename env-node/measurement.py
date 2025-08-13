from machine import I2C, Pin

import sht30

class Measurement:
    def __init__(self):
        i2c = I2C(1, sda=Pin(21), scl=Pin(22))
        self.sht = sht30.SHT30(i2c=i2c, i2c_address=68)

    def run(self):
        t_int, t_dec, h_int, h_dec = self.sht.measure_int()
        t = f"{t_int}.{t_dec:02d}"
        h = f"{h_int}.{h_dec:02d}"
        return t, h
