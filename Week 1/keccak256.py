from base64 import encode
from attr import has
import sha3
import os

print("Implementasi Keccak 256 di python\n")
print("+++++++++++++++++++++++++++++++++++++")
nama_ibu = input("Masukan Nama Ibu kandung: ")
os.system('CLS')
print("Nama Ibu kandung sebelum di hash: \n", nama_ibu)
encoded = nama_ibu.encode()
obj_encoded = sha3.keccak_256(encoded)
print("Nama Ibu kandung setelah di hash: \n", obj_encoded.hexdigest())


