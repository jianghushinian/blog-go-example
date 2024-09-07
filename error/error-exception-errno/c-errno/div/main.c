#include <stdio.h>

int div(int a, int b, int *result) {
    if (b == 0) {
        return -1; // 返回 -1 表示错误
    }
    *result = a / b; // 将结果存储在指针所指向的变量中
    return 0; // 返回 0 表示成功
}

int main() {
    int result;
    int err;

    err = div(1, 0, &result);

    // 错误处理
    if (err == -1) {
        printf("division by zero\n");
    } else {
        printf("result: %d\n", result);
    }

    return 0;
}
