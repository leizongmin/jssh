#!/usr/bin/env go run github.com/leizongmin/jssh

updateBuiltinJS();

const packageName = `github.com/leizongmin/jssh`;
const binName = `jssh`;
const goBuild = `go build -v -ldflags "-s -w"`;
const goProxy = `https://goproxy.cn`;

const releaseDir = path.join(__dirname, `release`);
const cacheDir = path.join(releaseDir, `cross_compile_cache`);

const goVersionOutput = exec2(`go version`).output.match(/go version go(.*) /);
if (!goVersionOutput) {
  log.error(`无法通过命令[go version]获得Go版本号`);
  exit(1);
}
const goVersion = goVersionOutput[1];
log.info(`当前Go版本号%s`, goVersion);

setenv(`GO111MODULE`, `on`);
setenv(`GOPROXY`, goProxy);

exec(`mkdir -p ${fixFilePath(releaseDir)}`);
fs.readdir(releaseDir).forEach((s) => {
  const p = path.join(releaseDir, s.name);
  if (p !== cacheDir) {
    exec(`rm -rf ${fixFilePath(p)}`);
  }
});

//**********************************************************************************************************************

updateReleasePkgInfo();
buildHostOSVersion();
if (__os === `darwin`) {
  buildLinuxVersionOnDocker();
}
buildReleaseFiles();
restoreReleasePkgInfo();

//**********************************************************************************************************************

function updateReleasePkgInfo() {
  log.info(`更新版本信息`);
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
  const file = path.join(__dirname, `internal/pkginfo/build_info.go`);
  const data = `
package pkginfo

const CommitHash = "${commitHash}"
const CommitDate = "${commitDate}"
const BuildGoVersion = "${goVersion}"
`.trimLeft();
  fs.writefile(file, data);
  log.info(data);
}

function restoreReleasePkgInfo() {
  exec(`git checkout internal/pkginfo/build_info.go`);
}

function getNormalOSType() {
  switch (__os) {
    case "darwin":
      return "osx";
    case "linux":
      return "linux";
    case "freebsd":
      return "freebsd";
    case "windows":
      return "windows";
    default:
      return "other";
  }
}

function buildHostOSVersion() {
  log.info(`构建宿主系统版本`);
  const binPath = path.join(releaseDir, getNormalOSType(), binName);
  exec(`${goBuild} -o ${fixFilePath(binPath)} ${packageName}`);
  log.info(`构建输出到%s`, binPath);
}

function buildLinuxVersionOnDocker() {
  if (exec(`which docker`).code !== 0) {
    log.info(`未安装Docker，无法构建Linux版本`);
    return;
  }
  log.info(`在macOS上通过Docker构建Linux版本`);
  const binPath = path.join(releaseDir, `linux`, binName);
  exec(`mkdir -p ${fixFilePath(cacheDir)}`);
  const ret = exec(
    `docker run --rm -it -v "${fixFilePath(cacheDir)}:/go" -v ${fixFilePath(
      __dirname
    )}:${fixFilePath(__dirname)} -w ${fixFilePath(
      __dirname
    )} -e GO111MODULE=on -e GOPROXY=${goProxy} golang:${goVersion} ${goBuild} -o ${fixFilePath(
      binPath
    )} ${packageName}`
  );
  if (ret.code !== 0) {
    log.error(`通过Docker构建失败`);
  }
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
      const tarFile = path.join(releaseDir, `${binName}-${s.name}`);
      exec(`tar -czvf --force-local --absolute-names ${fixFilePath(tarFile)}.tar.gz *`);
      cd(__dirname);
      log.info(`输出压缩包%s`, tarFile);
    }
  });

  const currentDir = path.join(releaseDir, "current");
  const currentTarFile = path.join(
    releaseDir,
    `jssh-${getNormalOSType()}.tar.gz`
  );
  const currentBinFile = path.join(releaseDir, getNormalOSType(), "jssh");
  exec(`mkdir -p ${fixFilePath(currentDir)}`);
  exec(
    `cp -f ${fixFilePath(currentTarFile)} ${fixFilePath(
      path.join(currentDir, "jssh.tar.gz")
    )}`
  );
  exec(
    `cp -f ${fixFilePath(currentBinFile)} ${fixFilePath(
      path.join(currentDir, "jssh")
    )}`
  );
}

function updateBuiltinJS() {
  log.info(`更新内置JS模块`);
  const dir = path.join(__dirname, `internal`, `jsbuiltin`);
  const list = [];
  fs.readdir(dir)
    .sort()
    .forEach((s) => {
      const f = path.join(dir, s.name);
      if (!s.isdir && f.endsWith(`.js`)) {
        log.info(`JS模块%s`, f);
        const code = JSON.stringify(JSON.stringify(fs.readfile(f)));
        list.push(`	// ${s.name}`);
        list.push(
          `	modules = append(modules, JsModule{File: "${s.name}", Code: ${code}})`
        );
        list.push(``);
      }
    });
  const goFile = path.join(__dirname, `internal`, `jsbuiltin`, `all.go`);
  fs.writefile(
    goFile,
    `
package jsbuiltin

var modules []JsModule

func init() {
	${list.join(`\n`).trim()}
}
`.trimLeft()
  );
}

function fixFilePath(p) {
  return p.replace(/\\/g, "/");
}
