#include <stdio.h>
#include <errno.h>
#include <string.h>  // 包含 strerror 函数的头文件

int main() {
    // 清除 errno 初始值，这是一个好的编程习惯
    errno = 0;

    FILE *file;
    char *filename = "example.txt";

    // 尝试以读取模式打开文件
    file = fopen(filename, "r");
    if (file == NULL) {
        // 打开文件出错
        printf("Failed to open file: %s\n", filename);
        // 查看错误码 errno
        printf("ErrNo: %d\n", errno);
        // perror 函数显示传入的字符串，后跟一个冒号、一个空格和当前 errno 值的文本表示形式
        perror("Error");
        // strerror 函数返回一个指针，指向当前 errno 值的文本表示形式
        printf("Error: %s\n", strerror(errno));
        // 在发生错误时，大多数的 C 或 UNIX 函数调用返回 1 或 NULL
        return 1;
    }

    // 文件操作
    printf("open file success\n");

    // process file...

    // 完成操作后，关闭文件
    fclose(file);

    return 0;
}
