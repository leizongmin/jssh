function exec1(cmd, env = {}) {
  return exec(cmd, env, 1);
}

function exec2(cmd, env = {}) {
  return exec(cmd, env, 2);
}

ssh.exec1 = function exec1(cmd, env = {}) {
  return ssh.exec(cmd, env, 1);
};

ssh.exec2 = function exec2(cmd, env = {}) {
  return ssh.exec(cmd, env, 2);
};
