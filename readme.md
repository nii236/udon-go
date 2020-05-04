```
$ go run *.go
.data_start

    ret_addr: %SystemUInt32, 0xFFFFFFFF
    this_trans: %UnityEngineTransform, this
    this_gameObj: %UnityEngineGameObject, this
    func1__SystemInt32_SystemInt32_x1: %SystemInt32, null
    func1__SystemInt32_SystemInt32_y1: %SystemInt32, null
    func2__SystemInt32_SystemInt32_x2: %SystemInt32, null
    func2__SystemInt32_SystemInt32_y2: %SystemInt32, null

.data_end

.code_start

    .export _start
        PUSH, ret_addr
        COPY
        JUMP_INDIRECT, ret_addr
        PUSH, func1__SystemInt32_SystemInt32_y1
        COPY
        PUSH, func1__SystemInt32_SystemInt32_x1
        COPY
        PUSH, ret_addr
        COPY
        JUMP_INDIRECT, ret_addr
        PUSH, func2__SystemInt32_SystemInt32_y2
        COPY
        PUSH, func2__SystemInt32_SystemInt32_x2
        COPY
        PUSH, ret_addr
        COPY
        JUMP_INDIRECT, ret_addr

.code_end
```
