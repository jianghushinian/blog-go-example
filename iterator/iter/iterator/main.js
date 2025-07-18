// NOTE: JavaScript 迭代器

function* generator(num) {
    for (let i = 0; i < num; i++) {
        yield i;
    }
}

for (const v of generator(5)) {
    console.log(v);
}
