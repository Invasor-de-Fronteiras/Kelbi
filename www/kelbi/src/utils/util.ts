export function randomArr<T>(arr: T[]) {
  return arr[Math.floor(Math.random() * arr.length)];
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

/**
 * Zero or less returns 'false', 100 or greater returns 'true'. Else return probability with required percentage.
 * @param percentage for example 100%, 50%, 0%.
 * @return true or false with required probability.
 */
export function probably(percentage) {
  const zeroToOne = Math.random(); // greater than or equal to 0.0 and less than 1.0
  const multiple = zeroToOne * 100; // greater than or equal to 0.0 and less than 100.0
  return multiple < percentage;
}
