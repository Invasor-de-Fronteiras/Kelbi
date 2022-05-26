export function randomArr<T>(arr: T[]) {
  return arr[Math.floor(Math.random() * arr.length)];
}

export function randomHexColor() {
  return '#' + (((1 << 24) * Math.random()) | 0).toString(16);
}

export function normalizeDate(date: Date) {
  const d = new Date(date);
  const month = d.getMonth() + 1;
  const day = d.getDate();
  const year = d.getFullYear();

  return `${day}/${month}/${year}`;
}

export function removeItem<T>(array: T[], validate: (item: T) => boolean) {
  for (let i = array.length - 1; i >= 0; i--) {
    if (validate(array[i])) {
      array.splice(i, 1);
    }
  }

  return array;
}

// this function was created for performance reasons for when the item has possibility in the final array
export function findByEnd<T>(array: T[], fn: (item: T) => boolean): T | null {
  for (let i = array.length - 1; i >= 0; i--) {
    if (fn(array[i])) {
      return array[i];
    }
  }

  return null;
}
