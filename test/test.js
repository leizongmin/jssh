function jssh(args, env) {
  const ret = exec2(`${__bin} ${args}`, env);
  if (ret.code !== 0) {
    log.warn(`jssh exec fail wait status #${ret.code}`);
  }
  return ret.output.trim();
}

function eq(expected, actual) {
  assert(expected === actual, `${expected} !== ${actual}`);
}

function test(title, handler) {
  try {
    handler();
  } catch (err) {
    log.error(`test [${title}] fail: ${err.message}\n${err.stack}`);
    return;
  }
  log.info(`test [${title}] succeed`);
}

test(`jssh eval`, function () {
  eq(`579`, jssh(`eval 123+456`));
  eq(`hello`, jssh(`eval "'hello'"`));
  eq(`123`, jssh(`eval "a=1,b=2,c=3,format('%v%v%v',a,b,c)"`));
});
