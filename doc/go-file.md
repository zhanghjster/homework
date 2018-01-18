总结一下golang常用的文件操作函数，这些函数定义在os包，

#### OpenFile(name string, flag int, perm FileMode) (*File, error)

通用的打开文件的函数，打开一个以’name'为名称， 标示为flag，权限为perm的文件

```go
package main

import (
	"os"
	"log"
)

func main() {
	file, err = os.OpenFile("test.log", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()
}
```



#### func Create(name string) (*File, error)

封装了OpenFile操作，创建一个文件，如果存在则清空文件里的内容，文件的权限是0666, 

```go
package main

import (
	"os"
	"log"
)

func main() {
    // 等同于
    // file, err = os.OpenFile("test.log", os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
	file, err := os.Create("test.log")
	if err != nil {
		log.Fatal(err)
	}
	file.Close()
}

```

#### func Open(name string) (*File, error)

封装OpenFile操作，以只读方式打开一个文件

```go
package main

import (
	"os"
	"log"
)

func main() {
    // 等同于
    // file, err = os.OpenFile("test.log", O_RDONLY, 0)
	file, err = os.Open("test.log")
	if err != nil {
		log.Fatal(err)
	}
	file.Close()
}
```

#### func Remove(name string) error 

删除文件

#### func Rename(org, dest string) error

重命名

#### func Stat(name string) (*FileInfo, error)

获取文件信息 FileInfo

```go
package main

import (
	"os"
	"log"
	"fmt"
)

func main() {
	fileInfo, err := os.Stat("test.log")
	fmt.Printf("Name, %s \n", fileInfo.Name())
	fmt.Printf("IsDir, %t \n", fileInfo.IsDir())
	fmt.Printf("Mode, %v \n", fileInfo.Mode())
	fmt.Printf("ModTime, %v \n", fileInfo.ModTime())
	fmt.Printf("Size, %v \n", fileInfo.Size())
}

```

#### func Link(old, new string) error 

创建硬链接

#### func Symlink(old, new string) error

创建软连接

#### func Chmod(name string, mod FileMode)

修改文件权限

#### 检查是否为链接

#### 检查是否为Device

https://www.devdungeon.com/content/working-files-go#everything_is_a_file