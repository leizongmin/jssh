#!/usr/bin/env -S go run github.com/leizongmin/jssh

// https://github.com/xaionaro/documentation/blob/master/golang/reduce-binary-size.md

const goVersion = getGoVersion();
log.info(`当前Go版本号%s`, goVersion);
const packageName = `github.com/leizongmin/jssh`;
const binName = `jssh`;
const goBuild = `go build -v -a -gcflags=all="-l -B" -ldflags "-s -w ${getReleaseLdflags()}"`;

const releaseDir = path.join(__dirname, `release`);
const cacheDir = path.join(releaseDir, `cross_compile_cache`);

exec(`mkdir -p ${fixFilePath(releaseDir)}`);
fs.readdir(releaseDir).forEach((s) => {
  const p = path.join(releaseDir, s.name);
  if (p !== cacheDir) {
    exec(`rm -rf ${fixFilePath(p)}`);
  }
});

//**********************************************************************************************************************

buildHostOSVersion();
buildReleaseFiles();

//**********************************************************************************************************************

function getGoVersion() {
  const goVersionOutput = exec2(`go version`).output.match(
    /go version go(.*) /
  );
  if (!goVersionOutput) {
    log.error(`无法通过命令[go version]获得Go版本号`);
    exit(1);
  }
  return goVersionOutput[1];
}

function getReleaseLdflags() {
  const date = exec2(`date +%Y%m%d`).output.trim();
  const time = exec2(`date +%H%M`).output.trim();
  const commitHash = exec2(`git rev-parse --short HEAD`).output.trim();
  const commitDate = exec2(
    `git for-each-ref --sort=-committerdate refs/heads/ --format="%(authordate:short)" | head -n 1`
  )
    .output.trim()
    .replace(/\-/g, ``);
  if (!date || !commitHash) {
    log.error(`无法获取date和commit信息`);
    exit(1);
  }
  const list = [];
  list.push(`-X '${packageName}/internal/pkginfo.CommitHash=${commitHash}'`);
  list.push(`-X '${packageName}/internal/pkginfo.CommitDate=${commitDate}'`);
  list.push(`-X '${packageName}/internal/pkginfo.GoVersion=${goVersion}'`);
  return list.join(" ");
}

function buildHostOSVersion() {
  log.info(`构建宿主系统版本`);
  let binPath = path.join(releaseDir, `${__os}-${__arch}`, binName);
  if (__os === "windows") {
    binPath += ".exe";
  }
  exec(`rm -f ${fixFilePath(binPath)}`);
  const cmd = `${goBuild} -o ${fixFilePath(binPath)} ${packageName}`;
  log.info(cmd);
  exec(cmd);
  log.info(`构建输出到%s`, binPath);
  const version = exec2(`${fixFilePath(binPath)} version`).output.trim();
  log.info(`已构建的jssh版本：${version}`);
}

function buildReleaseFiles() {
  log.info(`输出发布压缩包`);
  const dtsFile = path.join(__dirname, `jssh.d.ts`);
  fs.readdir(releaseDir).forEach((s) => {
    if (s.name.startsWith(`.`)) return;
    const p = path.join(releaseDir, s.name);
    if (p !== cacheDir) {
      cd(__dirname);
      exec(`cp -f ${fixFilePath(dtsFile)} ${fixFilePath(p)}`);
      cd(p);
      if (__os === "darwin") {
        // macOS系统使用zip压缩，解决github action的runner用tar打包出来的文件损坏问题
        const zipFile = `${binName}-${s.name}.zip`;
        const cmd = `zip ../${zipFile} *`;
        log.info(cmd);
        exec(cmd);
        log.info(`输出压缩包%s`, zipFile);
      } else {
        const tarFile = `${binName}-${s.name}.tar.gz`;
        const cmd = `tar -czvf ../${tarFile} *`;
        log.info(cmd);
        exec(cmd);
        log.info(`输出压缩包%s`, tarFile);
      }
      cd(__dirname);
    }
  });
}

function fixFilePath(p) {
  return p.replace(/\\/g, "/");
}
