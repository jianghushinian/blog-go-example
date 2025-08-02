import json
import yaml

yaml_example = """
object:
  a: 1
  1: 2
  "1": 3
  key: value
  array:
  - null_value:
  - boolean: true
  - integer: 1
"""


def main():
    y = yaml.load(yaml_example, Loader=yaml.SafeLoader)
    print(type(y), y)

    print("--------------------------")

    j = json.dumps(y)
    print(type(j), j)

    print("--------------------------")

    d = json.loads(j)
    print(type(d), d)

if __name__ == '__main__':
    main()
