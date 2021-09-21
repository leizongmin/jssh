function println(...args) {
  __builtin_op_stdout_write(args.map(v => v.toString()).join(" ") + "\n");
}
