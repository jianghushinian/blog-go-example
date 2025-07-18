# NOTE: Python 迭代器

def generator(num: int):
    for i in range(num):
        yield i


for value in generator(5):
    print(value)
