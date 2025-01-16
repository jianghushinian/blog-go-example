package main

import (
	"log"
	"os/exec"
)

// NOTE: 运行一个命令

func main() {
	// 创建一个命令
	cmd := exec.Command("echo", "Hello, World!")

	// 执行命令并等待命令完成
	err := cmd.Run() // 执行后控制台不会有任何输出
	if err != nil {
		log.Fatalf("Command failed: %v", err)
	}

	// fmt.Println(cmd.String())
}

// NOTE: 带 context 的 CommandContext

// func main() {
// 	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
// 	defer cancel()
//
// 	cmd := exec.CommandContext(ctx, "sleep", "5")
//
// 	if err := cmd.Run(); err != nil {
// 		log.Fatalf("Command failed: %v", err) // signal: killed
// 	}
// }

// NOTE: 后台运行命令

// func main() {
// 	cmd := exec.Command("sleep", "3")
//
// 	// 执行命令（非阻塞，不会等待命令执行完成）
// 	if err := cmd.Start(); err != nil {
// 		log.Fatalf("Command start failed: %v", err)
// 		return
// 	}
//
// 	fmt.Println("Command running in the background...")
//
// 	// 阻塞等待命令完成
// 	if err := cmd.Wait(); err != nil {
// 		log.Fatalf("Command wait failed: %v", err)
// 		return
// 	}
//
// 	log.Println("Command finished")
// }

// NOTE: 获取命令的输出

// func main() {
// 	// 创建一个命令
// 	cmd := exec.Command("echo", "Hello, World!")
//
// 	// 执行命令，并获取命令的输出，Output 内部会调用 Run 方法
// 	output, err := cmd.Output()
// 	if err != nil {
// 		log.Fatalf("Command failed: %v", err)
// 	}
//
// 	fmt.Println(string(output)) // Hello, World!
// }

// NOTE: 获取组合的标准输出和错误输出

// func main() {
// 	cmd := exec.Command("echox", "Hello, World!") // 执行一个不存在的命令 echox
//
// 	// 获取 标准输出 + 标准错误输出 组合内容
// 	output, err := cmd.CombinedOutput()
// 	if err != nil {
// 		log.Fatalf("Command failed: %v", err)
// 	}
//
// 	fmt.Println(string(output)) // 命令执行报错，会返回 err，所以这里不会输出内容
// }

// func main() {
// 	// 使用一个命令，既产生标准输出，也产生标准错误输出
// 	cmd := exec.Command("sh", "-c", "echo 'This is stdout'; echo 'This is stderr' >&2")
//
// 	// 获取 标准输出 + 标准错误输出 组合内容
// 	output, err := cmd.CombinedOutput()
// 	if err != nil {
// 		log.Fatalf("Command execution failed: %v", err)
// 	}
//
// 	// 打印组合输出
// 	fmt.Printf("Combined Output:\n%s", string(output))
// }

// NOTE: 设置标准输出和错误输出

// func main() {
// 	cmd := exec.Command("ls", "-l")
//
// 	// 设置标准输出和标准错误输出到当前进程，执行后可以在控制台看到命令执行的输出
// 	cmd.Stdout = os.Stdout
// 	cmd.Stderr = os.Stderr
//
// 	if err := cmd.Run(); err != nil {
// 		log.Fatalf("Command failed: %v", err)
// 	}
// }

// NOTE: 使用标准输入传递数据

// func main() {
// 	cmd := exec.Command("grep", "hello")
//
// 	// 通过标准输入传递数据给命令
// 	cmd.Stdin = bytes.NewBufferString("hello world!\nhi there\n")
//
// 	// 获取标准输出
// 	output, err := cmd.Output()
// 	if err != nil {
// 		log.Fatalf("Command failed: %v", err)
// 		return
// 	}
//
// 	fmt.Println(string(output)) // hello world!
// }

// func main() {
// 	file, err := os.Open("demo.log") // 打开一个文件
// 	if err != nil {
// 		log.Fatalf("Open file failed: %v\n", err)
// 		return
// 	}
// 	defer file.Close()
//
// 	cmd := exec.Command("cat")
// 	cmd.Stdin = file       // 将文件作为 cat 的标准输入
// 	cmd.Stdout = os.Stdout // 获取标准输出
//
// 	if err := cmd.Run(); err != nil {
// 		log.Fatalf("Command failed: %v", err)
// 	}
// }

// NOTE: 设置和使用环境变量

// func main() {
// 	cmd := exec.Command("printenv", "ENV_VAR")
//
// 	log.Printf("ENV: %+v\n", cmd.Environ())
//
// 	// 设置环境变量
// 	cmd.Env = append(cmd.Environ(), "ENV_VAR=HelloWorld")
// 	// cmd.Env = append(cmd.Env, "ENV_VAR=HelloWorld")
// 	// cmd.Env = []string{"ENV_VAR=HelloWorld"} // 全新的 ENV，原来内置的都丢弃了
//
// 	log.Printf("ENV: %+v\n", cmd.Environ())
//
// 	// 获取输出
// 	output, err := cmd.Output()
// 	if err != nil {
// 		log.Fatalf("Command failed: %v", err)
// 	}
//
// 	fmt.Println(string(output)) // HelloWorld
// }

// NOTE: Pipe

// func main() {
// 	// 命令中使用了管道
// 	cmdEcho := exec.Command("echo", "hello world\nhi there")
//
// 	outPipe, err := cmdEcho.StdoutPipe()
// 	if err != nil {
// 		log.Fatalf("Command failed: %v", err)
// 	}
//
// 	// 注意，这里不能使用 Run 方法阻塞等待，应该使用非阻塞的 Start 方法
// 	if err := cmdEcho.Start(); err != nil {
// 		log.Fatalf("Command failed: %v", err)
// 	}
//
// 	cmdGrep := exec.Command("grep", "hello")
// 	cmdGrep.Stdin = outPipe
// 	output, err := cmdGrep.Output()
// 	if err != nil {
// 		log.Fatalf("Command failed: %v", err)
// 	}
//
// 	fmt.Println(string(output)) // hello world
// }

// NOTE: 使用 `bash -c` 执行复杂命令

// func main() {
// 	// 命令中使用了管道
// 	cmd := exec.Command("bash", "-c", "echo 'hello world\nhi there' | grep hello")
//
// 	output, err := cmd.Output()
// 	if err != nil {
// 		log.Fatalf("Command failed: %v", err)
// 	}
//
// 	fmt.Println(string(output)) // hello world
// }

// NOTE: 指定工作目录

// func main() {
// 	cmd := exec.Command("cat", "demo.log")
// 	cmd.Stdout = os.Stdout // 获取标准输出
// 	cmd.Stderr = os.Stderr // 获取错误输出
//
// 	// cmd.Dir = "/tmp" // 指定绝对目录
// 	cmd.Dir = "." // 指定相对目录
//
// 	if err := cmd.Run(); err != nil {
// 		log.Fatalf("Command failed: %v", err)
// 	}
// }

// NOTE: 捕获退出状态

// func main() {
// 	// 查看一个不存在的目录
// 	cmd := exec.Command("ls", "/nonexistent")
//
// 	// 运行命令
// 	err := cmd.Run()
//
// 	// 检查退出状态
// 	var exitError *exec.ExitError
// 	if errors.As(err, &exitError) {
// 		log.Fatalf("Process PID: %d exit code: %d", exitError.Pid(), exitError.ExitCode()) // 打印 pid 和退出码
// 	}
// }

// NOTE: 搜索可执行文件

// func main() {
// 	path, err := exec.LookPath("ls")
// 	if err != nil {
// 		log.Fatal("installing ls is in your future")
// 	}
// 	fmt.Printf("ls is available at %s\n", path)
// }

// func main() {
// 	path, err := exec.LookPath("lsx")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Printf("ls is available at %s\n", path)
// }
