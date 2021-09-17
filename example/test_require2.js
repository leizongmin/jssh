const xss = require("https://unpkg.com/xss@1.0.9/lib/index.js");

while (true) {
  println("请输入一段HTML（直接回车表示结束）：")
  const line = readline().trimRight();
  if (!line) {
    break
  }
  console.log("输出：");
  console.log(xss(line));
}
