from pywebio.input import input, FLOAT
from pywebio.output import put_text
from pywebio import start_server

def user_info():
    name = input("Enter your name:", type=FLOAT)
    age = input("Enter your age:", type=FLOAT)
    put_text("Your name is:", name)
    put_text("Your age is:", age)

if __name__ == '__main__':
    start_server(user_info, port=8080)
  
