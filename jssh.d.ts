/** JSSH版本号 */
declare const __version: string;

/** JSSH二进制文件路径 */
declare const __bin: string;

/** 当前进程PID */
declare const __pid: number;

/** 临时文件目录 */
declare const __tmpdir: string;

/** 当前用户HOME目录 */
declare const __homedir: string;

/** 当前主机名 */
declare const __hostname: string;

/** 当前脚本目录名 */
declare var __dirname: string;

/** 当前脚本文件名 */
declare var __filename: string;

/** 当前命令行参数 */
declare const __args: string[];

/** 当前环境变量 */
declare const __env: Record<string, string>;

/** 最近一次执行命令输出的内容 */
declare const __output: string;

/** 最近一次执行命令输出内容的字节数 */
declare const __outputbytes: number;

/** 最近一次执行命令进程退出code */
declare const __code: number;


/**
 * 设置全局变量
 * @param name 变量名
 * @param value 值
 * @return 是否成功
 */
declare function set(name: string, value: any): boolean;

/**
 * 获取全局变量
 * @param name 变量名
 * @return 值
 */
declare function get(name: string): any;

/**
 * 格式化文本内容
 * @param format 模板，支持%s等格式
 * @param args 参数列表
 * @return 格式化后的文本内容
 */
declare function format(format: any, ...args: any[]): string;

/**
 * 输出到控制台
 * @param format 模板，支持%s等格式
 * @param args 参数列表
 * @return 是否成功
 */
declare function print(format: any, ...args: any[]): boolean;

/**
 * 输出一行到控制台
 * @param format 模板，支持%s等格式
 * @param args 参数列表
 * @return 是否成功
 */
declare function println(format: any, ...args: any[]): boolean;


/**
 * 睡眠
 * @param milliseconds 毫秒
 * @return 毫秒
 */
declare function sleep(milliseconds: number): number;

/**
 * 退出进程
 * @param code 进程退出code
 */
declare function exit(code?: number): void;

/**
 * Shell相关操作模块
 */
declare const sh: ShModule

/**
 * 文件相关操作模块
 */
declare const fs: FsModule;

/**
 * 文件路径相关操作模块
 */
declare const path: PathModule;

/**
 * 命令行参数相关操作模块
 */
declare const cli: CliModule;

/**
 * HTTP相关操作模块
 */
declare const http: HttpModule;

/**
 * 日志相关操作模块
 */
declare const log: LogModule;

interface FsModule {
    /**
     * 读取目录
     * @param path 路径
     * @return 文件信息列表
     */
    readdir(path: string): FileStat[];

    /**
     * 读取文件内容
     * @param path 路径
     * @return 文件内容
     */
    readfile(path: string): string;

    /**
     * 读取文件信息
     * @param path 路径
     * @return 文件信息
     */
    stat(path: string): FileStat;

    /**
     * 写入文件
     * @param path 路径
     * @param data 内容
     * @return 是否成功
     */
    writefile(path: string, data: string): boolean;

    /**
     * 追加内容到文件末尾
     * @param path 路径
     * @param data 内容
     * @return 是否成功
     */
    appendfile(path: string, data: string): boolean;
}

interface FileStat {
    /** 文件名 */
    name: string;
    /** 是否为目录 */
    isdir: boolean;
    /** 文件mode，如0644 */
    mode: number;
    /** 文件最后修改秒时间戳 */
    modtime: number;
    /** 文件大小 */
    size: number;
}

interface PathModule {
    /**
     * 拼接文件路径
     * @param args 子路径列表
     * @return 拼接后的路径
     */
    join(...args: string[]): string;

    /**
     * 获取绝对路径
     * @param path 原始路径
     * @return 绝对路径
     */
    abs(path: string): string;

    /**
     * 获取文件名
     * @param path 路径
     * @return 文件名
     */
    base(path: string): string;

    /**
     * 获取文件扩展名
     * @param path 路径
     * @return 文件扩展名
     */
    ext(path: string): string;

    /**
     * 获取路径的上级目录
     * @param path 路径
     * @return 上级目录
     */
    dir(path: string): string;
}

interface CliModule {
    /**
     * 获取指定名称的参数
     * @param flag 参数名称
     * @return 参数值
     */
    get(flag: string): string;

    /**
     * 获取指定索引的参数
     * @param index 参数索引
     * @return 参数值
     */
    get(index: number): string;

    /**
     * 获取指定名称的参数的布尔值，f,false,0表示假，其余为真
     * @param flag 参数名称
     * @return 参数值
     */
    bool(flag: string): boolean;

    /**
     * 获取args参数列表
     * @return 参数列表
     */
    args(): string[];

    /**
     * 获取opts参数
     * @return 参数Map
     */
    opts(): Record<string, string>;
}

interface HttpModule {
    /**
     * 设置HTTP请求的超时时间
     * @param milliseconds 毫秒
     * @return 毫秒
     */
    timeout(milliseconds: number): number;

    /**
     * 发送HTTP请求
     * @param method 请求方法
     * @param url 请求URL
     * @param headers 请求头
     * @param body 请求体
     * @return 响应结果
     */
    request(method: string, url: String, headers?: Record<string, string>, body?: string): HttpResponse;
}

interface HttpResponse {
    /** 状态码 */
    status: number;
    /** 响应头 */
    headers: Record<string, string | string[]>;
    /** 响应体 */
    body: string;
}

interface LogModule {
    /**
     * 日志输出到控制台
     * @param format 模板，支持%s等格式
     * @param args 参数列表
     * @return 是否成功
     */
    info(format: any, ...args: any[]): boolean;

    /**
     * 日志输出到控制台
     * @param format 模板，支持%s等格式
     * @param args 参数列表
     * @return 是否成功
     */
    error(format: any, ...args: any[]): boolean;
}

interface ShModule {
    /**
     * 设置环境变量
     * @param name 环境变量名称
     * @param value 环境变量值
     * @return 是否成功
     */
    setenv(name: string, value: string): boolean;

    /**
     * 执行命令
     * @param cmd 命令
     * @param env 额外的环境变量
     * @param combineOutput 是否合并输出，当为true时不直接输出命令执行结果，而存储到__output变量中
     * @return 进程信息
     */
    exec(cmd: string, env?: Record<string, string>, combineOutput?: boolean): ExecResult;

    /**
     * 后台执行命令
     * @param cmd 命令
     * @param env 额外的环境变量
     * @param combineOutput 是否合并输出，当为true时不直接输出命令执行结果，而存储到__output变量中
     * @return 进程信息
     */
    bgexec(cmd: string, env?: Record<string, string>, combineOutput?: boolean): ExecResult;

    /**
     * 改变当前工作目录
     * @param dir 目录路径
     * @return 是否成功
     */
    chdir(dir: string): boolean;

    /**
     * 改变当前工作目录
     * @param dir 目录路径
     * @return 是否成功
     */
    cd(dir: string): boolean;

    /**
     * 获取当前工作目录
     * @return 当前工作目录路径
     */
    cwd(): string;

    /**
     * 获取当前工作目录
     * @return 当前工作目录路径
     */
    pwd(): string;
}

interface ExecResult {
    /**
     * 进程PID
     */
    pid: number;
    /**
     * 进程退出code
     */
    code?: number;
    /**
     * 进程输出内容，仅当combineOutput=true时有效
     */
    output?: string;
    /**
     * 进出输出内容字节数，仅当combineOutput=true时有效
     */
    outputbytes?: number;
}