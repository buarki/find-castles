const facilities = {
  name: false,
  vaca: true,
};

const x = Object.entries(facilities).map(([key, value]) => ({name: key, value: value}));
console.log({x});