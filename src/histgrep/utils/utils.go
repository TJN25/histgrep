package utils

import (
    "runtime"
)

func Btoi(b bool) int {
    if b {
        return 1
    }
    return 0
 }

func CallerName(skip int) string {
        pc, _, _, ok := runtime.Caller(skip + 1)
        if !ok {
                return ""
        }
        f := runtime.FuncForPC(pc)
        if f == nil {
                return ""
        }
        return f.Name()
}

//     log3("File {file} had error {error}", "file", file, "error", err)

// func Fstring(format string, args ...interface{}) string {
//     args2 := make([]string, len(args))
//     for i, v := range args {
//         if i%2 == 0 {
//             args2[i] = fmt.Sprintf("{%v}", v)
//         } else {
//             args2[i] = fmt.Sprint(v)
//         }
//     }
//     r := strings.NewReplacer(args2...)
//     return r.Replace(format)
// }

