import logging


def div(a, b):
    return a / b


try:
    result = div(1, 0)
    print(result)
except ZeroDivisionError as e:
    logging.error(e)
except Exception as e:
    logging.error(e)


class MyException(Exception):
    ...


# raise MyException("my custom exception")


def a():
    ...


def b():
    ...


def c():
    ...


def main():
    try:
        a()
        b()
        c()
    except Exception as e:
        logging.error(e)
    else:
        print("success")
    finally:
        print("release resources")


if __name__ == '__main__':
    main()
