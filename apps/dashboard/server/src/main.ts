import { t2 } from 'utils';

import { t3 } from '@shared/utils';

import { t1 } from './test';

// const t: string = 1;

t1();
t2();
t3();

console.log('dashboard');
console.log('dashboard2');
console.log('dashboard3');
console.log('dashboard4');
console.log('dashboard5');

setTimeout(() => {
  console.log('dashboard6');
}, 999999999);
