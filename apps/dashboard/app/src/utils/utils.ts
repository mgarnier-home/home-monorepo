namespace Utils {
  export function getPingColor(nb: number): string {
    if (nb > 400) {
      return "hsl(0, 0%, 0%)";
    } else if (nb > 200) {
      const h = 30 - ((nb - 200) / 200) * 30;

      return `hsl(${h}, 100%, 50%)`;
    } else if (nb > 80) {
      const h = 60 - ((nb - 80) / 120) * 30;

      return `hsl(${h}, 100%, 50%)`;
    } else {
      const h = 120 - (nb / 80) * 60;

      return `hsl(${h}, 100%, 50%)`;
    }
  }

  export function getLineColor(str: string) {
    let hash = 0;
    for (let i = 0; i < str.length; i++) {
      hash = str.charCodeAt(i) + ((hash << 5) - hash);
    }
    let color = "#";
    for (let i = 0; i < 3; i++) {
      const v = (hash >> (i * 8)) & 0xff;
      color += ("00" + v.toString(16)).substr(-2);
    }

    return color;
  }

  export function padLeft(str: string, length: number, char: string) {
    return str.length >= length ? str : new Array(length - str.length + 1).join(char) + str;
  }

  export function dateFormatter(timestamp: number) {
    const date = new Date(timestamp);

    return `${date.getHours().toString().padStart(2, "0")}:${date.getMinutes().toString().padStart(2, "0")}:${date
      .getSeconds()
      .toString()
      .padStart(2, "0")}`;
  }
}

export default Utils;
