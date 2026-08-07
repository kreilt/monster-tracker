global.window = {};
require('./flavors-data.js');
require('fs').writeFileSync(
  'flavors-data.json',
  JSON.stringify(window.MONSTER_DATA, null, 2)
);
console.log('OK, записано записей:', window.MONSTER_DATA.length);