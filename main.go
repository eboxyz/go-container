package main

import(
    "fmt"
    "os"
    "os/exec"
    "syscall"
)

func main(){
  //switch case function
  //'run' will run the parent, which contains an in-memory image of the current executable file
  switch os.Args[1]{
  case "run":
      parent()
  case "child":
      child()
  default:
        panic("what should I do")
  }
}

func parent() {
  //runs the executable on an in-memory image
  cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
  //setting namespaces for when the child process is running
  cmd.SysProcAttr = &syscall.SysProcAttr{
    Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
  }
  //standard input/output/error file descriptors
  cmd.Stdin = os.Stdin
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr

  if err := cmd.Run(); err != nil {
    fmt.Println("ERROR", err)
    os.Exit(1)
  }
}

func child(){
  //swap into a root filesystem
  //the new directory is '/' as opposed to rootfs/
  //pivot root is used to swap two filesystems that are not part of the same tree
  must(syscall.Mount("rootfs", "rootfs", "", syscall.MS_BIND, ""))
  must(os.MkdirAll("rootfs/oldrootfs", 0700))
  must(syscall.PivotRoot("rootfs", "rootfs/oldrootfs"))
  must(os.Chdir("/"))

  cmd := exec.Command(os.Args[2], os.Args[3:]...)
  cmd.Stdin = os.Stdin
  cmd.Stdout = os.Stdout
  cmd.Stderr = os.Stderr

  if err := cmd.Run(); err != nil {
    fmt.Println("ERROR", err)
    os.Exit(1)
  }
}

//error handling for must
func must(err error){
  if err != nil {
    panic(err)
  }
}









