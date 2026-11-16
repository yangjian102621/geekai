import { spawn } from 'node:child_process';
import { fileURLToPath } from 'node:url';
import { resolve } from 'node:path';

const root = fileURLToPath(new URL('..', import.meta.url));
const processes = [
  spawn(process.execPath, [resolve(root, '../admin-web/node_modules/vite/bin/vite.js')], {
    cwd: resolve(root, '../admin-web'),
    stdio: 'inherit',
    windowsHide: true,
  }),
  spawn(process.execPath, [resolve(root, 'node_modules/vite/bin/vite.js'), ...process.argv.slice(2)], {
    cwd: root,
    stdio: 'inherit',
    windowsHide: true,
  }),
];

let closing = false;
function stop(code) {
  if (closing) return;
  closing = true;
  process.exitCode = code;
  for (const child of processes) child.kill();
}

for (const child of processes) {
  child.on('error', (error) => {
    console.error(error.message);
    stop(1);
  });
  child.on('exit', (code) => stop(code ?? 1));
}
process.on('SIGINT', () => stop(0));
process.on('SIGTERM', () => stop(0));
