const yaml = require('js-yaml');

const yamlExample = `
? - Detroit Tigers
  - Chicago cubs
: - 2001-07-23

? [ New York Yankees,
    Atlanta Braves ]
: [ 2001-07-02, 2001-08-12,
    2001-08-14 ]
`;
try {
    const data = yaml.load(yamlExample);
    console.log(data);
    // output:
    // {
    //   'Detroit Tigers,Chicago cubs': [ 2001-07-23T00:00:00.000Z ],
    //   'New York Yankees,Atlanta Braves': [
    //     2001-07-02T00:00:00.000Z,
    //     2001-08-12T00:00:00.000Z,
    //     2001-08-14T00:00:00.000Z
    //   ]
    // }
} catch (e) {
    console.error(e);
}

/**
 * $ npm install js-yaml
 * $ node js-examples/main.js
 * */
