import { Color, LogLevel, logLevels } from './interfaces';

const inBrowser = typeof window !== 'undefined';

const envLogLevel = inBrowser ? NaN : parseInt(process?.env?.LOG_LEVEL ?? '');

interface ILogger {
  // log(color: Color, ...args: any[]): void;
  info(...args: any[]): void;
  debug(...args: any[]): void;
  warn(...args: any[]): void;
  error(...args: any[]): void;
  verbose(...args: any[]): void;
}

class Logger implements ILogger {
  public colored: ColoredLogger = new ColoredLogger(this);
  private appName: string = 'default-app-name';
  private logLevel: LogLevel = isNaN(envLogLevel) ? LogLevel.INFO : envLogLevel;

  private static getColor(def: Color, arg: any): Color {
    return typeof arg === 'string' && arg.startsWith('\x1b') ? (arg as Color) : def;
  }

  private static getCSSColor(color: Color): string {
    switch (color) {
      case Color.BLACK:
        return 'color: black;';
      case Color.RED:
        return 'color: red;';
      case Color.GREEN:
        return 'color: green;';
      case Color.YELLOW:
        return 'color: yellow;';
      case Color.BLUE:
        return 'color: blue;';
      case Color.MAGENTA:
        return 'color: magenta;';
      case Color.CYAN:
        return 'color: cyan;';
      case Color.WHITE:
        return 'color: white;';
      case Color.GRAY:
        return 'color: gray;';
      default:
        return '';
    }
  }

  private getLog(color: Color, level: LogLevel, ...args: any[]) {
    const date = new Date().toISOString();

    if (inBrowser) {
      return [`%c[${date}] [${logLevels[level]}] [${this.appName}] `, Logger.getCSSColor(color), ...args];
    } else {
      return [color, `[${date}] [${logLevels[level]}] [${this.appName}] `, ...args, '\x1b[0m'];
    }
  }

  setAppName(appName: string) {
    this.appName = appName;
  }

  setLogLevel(logLevel: LogLevel) {
    this.logLevel = logLevel;
  }

  info(...args: any[]) {
    if (this.logLevel >= LogLevel.INFO) {
      console.info(...this.getLog(Logger.getColor(Color.DEFAULT, args[0]), LogLevel.INFO, ...args));
    }
  }

  debug(...args: any[]) {
    if (this.logLevel >= LogLevel.DEBUG) {
      console.debug(...this.getLog(Logger.getColor(Color.BLUE, args[0]), LogLevel.DEBUG, ...args));
    }
  }

  warn(...args: any[]) {
    if (this.logLevel >= LogLevel.WARN) {
      console.warn(...this.getLog(Logger.getColor(Color.YELLOW, args[0]), LogLevel.WARN, ...args));
    }
  }

  error(...args: any[]) {
    if (this.logLevel >= LogLevel.ERROR) {
      console.error(...this.getLog(Logger.getColor(Color.RED, args[0]), LogLevel.ERROR, ...args));
    }
  }

  verbose(...args: any[]) {
    if (this.logLevel >= LogLevel.VERBOSE) {
      console.log(...this.getLog(Logger.getColor(Color.DEFAULT, args[0]), LogLevel.VERBOSE, ...args));
    }
  }
}

class ColoredLogger implements ILogger {
  constructor(private logger: Logger) {}

  info(color: Color, ...args: any[]) {
    this.logger.info(color, ...args);
  }

  debug(color: Color, ...args: any[]) {
    this.logger.debug(color, ...args);
  }

  warn(color: Color, ...args: any[]) {
    this.logger.warn(color, ...args);
  }

  error(color: Color, ...args: any[]) {
    this.logger.error(color, ...args);
  }

  verbose(color: Color, ...args: any[]) {
    this.logger.verbose(color, ...args);
  }
}

export const logger = new Logger();
